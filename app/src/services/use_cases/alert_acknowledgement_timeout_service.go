package use_cases

import (
	"time"

	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/interfaces"
)

/*
Given a Monitored Service in an Unhealthy State,
the corresponding Alert is not Acknowledged
and the last level has not been notified,
when the Pager receives the Acknowledgement Timeout,
then the Pager notifies all targets of the next level of the escalation policy
and sets a 15-minutes acknowledgement delay.
*/
type AlertAcknowledgementTimeoutService struct {
	alertRepository            interfaces.AlertRepository
	monitoredServiceRepository interfaces.MonitoredServiceRepository
	escalationPolicyRepository interfaces.EscalationPolicyRepository
	notificationService        interfaces.NotificationService
	timerService               interfaces.TimerService
	lockService                interfaces.LockService
}

func NewAlertAcknowledgementTimeoutService(alertRepository interfaces.AlertRepository,
	monitoredServiceRepository interfaces.MonitoredServiceRepository,
	escalationPolicyRepository interfaces.EscalationPolicyRepository,
	notificationService interfaces.NotificationService,
	timerService interfaces.TimerService,
	lockService interfaces.LockService) AlertAcknowledgementTimeoutService {
	return AlertAcknowledgementTimeoutService{
		alertRepository:            alertRepository,
		monitoredServiceRepository: monitoredServiceRepository,
		escalationPolicyRepository: escalationPolicyRepository,
		notificationService:        notificationService,
		timerService:               timerService,
		lockService:                lockService,
	}
}

func (service *AlertAcknowledgementTimeoutService) AcknowledgementTimeout(alertID string) error {
	alert, err := service.alertRepository.GetByID(alertID)

	//error getting alert or not found alert
	if err != nil {
		return err
	}

	if alert.Status == constants.AlertStatusAcknowledge {
		return nil
	}

	monitoredService, err := service.monitoredServiceRepository.GetByID(alert.MonitoredServiceID)
	if err != nil {
		return err
	}
	//the service is already healthy, finish the process.
	if monitoredService.Status != constants.ServiceStatusUnhealthy {
		return nil
	}

	lockMonService := service.lockService.Lock(monitoredService.ID)
	defer service.lockService.Unlock(monitoredService.ID)

	lockAlert := service.lockService.Lock(alert.ID)
	defer service.lockService.Unlock(alert.ID)

	if !lockMonService || !lockAlert {
		return constants.ErrorLockedResource
	}

	escalationPolicy, err := service.escalationPolicyRepository.GetByServiceID(monitoredService.ID)
	if err != nil {
		return err
	}

	level := escalationPolicy.GetLevelToNotify(alert.GetNextLevelToNotify())
	if level == nil {
		return nil
	}

	alert.NotifiedLevelID = &level.ID
	err = service.alertRepository.Save(alert)
	if err != nil {
		return err
	}

	service.notificationService.SendNotifications(alert, level.Targets)

	//set scheduler 15 minutes.
	service.timerService.Schedule("alert", alert.ID, time.Now().Add(time.Minute*15))

	return nil
}
