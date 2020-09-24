package use_cases

import (
	"time"

	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
	"github.com/santi8ago8/pager-service/app/src/interfaces"
)

/*
Given a Monitored Service in a Healthy State,
when the Pager receives an Alert related to this Monitored Service,
then the Monitored Service becomes Unhealthy,
the Pager notifies all targets of the first level of the escalation policy,
and sets a 15-minutes acknowledgement delay

Given a Monitored Service in an Unhealthy State
when the Pager receives the Acknowledgement
and later receives the Acknowledgement Timeout,
then the Pager doesn't notify any Target
and doesn't set an acknowledgement delay.
*/

type PagerReceiveAlertService struct {
	alertRepository            interfaces.AlertRepository
	monitoredServiceRepository interfaces.MonitoredServiceRepository
	escalationPolicyRepository interfaces.EscalationPolicyRepository
	notificationService        interfaces.NotificationService
	timerService               interfaces.TimerService
}

func NewPagerReceiveAlertService(alertRepository interfaces.AlertRepository,
	monitoredServiceRepository interfaces.MonitoredServiceRepository,
	escalationPolicyRepository interfaces.EscalationPolicyRepository,
	notificationService interfaces.NotificationService,
	timerService interfaces.TimerService) PagerReceiveAlertService {
	return PagerReceiveAlertService{
		alertRepository:            alertRepository,
		monitoredServiceRepository: monitoredServiceRepository,
		escalationPolicyRepository: escalationPolicyRepository,
		notificationService:        notificationService,
		timerService:               timerService,
	}
}

func (service *PagerReceiveAlertService) ReceivesAlert(serviceID string, alertMessage string) (*entities.Alert, error) {

	monitoredService, err := service.monitoredServiceRepository.GetByID(serviceID)
	if err != nil {
		return nil, err
	}
	//the service is already unhealthy, finish the process.
	if monitoredService.Status == constants.ServiceStatusUnhealthy {
		return nil, nil
	}
	monitoredService.SetUnhealthy()
	err = service.monitoredServiceRepository.Save(monitoredService)
	if err != nil {
		return nil, err
	}

	alert := entities.NewAlert(serviceID, alertMessage)

	escalationPolicy, err := service.escalationPolicyRepository.GetByServiceID(monitoredService.ID)
	if err != nil {
		return nil, err
	}

	level := escalationPolicy.GetLevelToNotify(alert.GetNextLevelToNotify())

	alert.NotifiedLevelID = &level.ID
	err = service.alertRepository.Save(alert)
	if err != nil {
		return nil, err
	}

	service.notificationService.SendNotifications(alert, level.Targets)

	//set scheduler 15 minutes.
	service.timerService.Schedule("alert", alert.ID, time.Now().Add(time.Minute*15))
	return alert, nil
}
