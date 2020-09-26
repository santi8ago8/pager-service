package use_cases

import (
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/interfaces"
)

/*
Given a Monitored Service in an Unhealthy State
when the Pager receives the Acknowledgement (alert_acknowledgement_service responsibility) (this file)
and later receives the Acknowledgement Timeout, (alert_acknowledgement_timeout_service responsibility)
then the Pager doesn't notify any Target
and doesn't set an acknowledgement delay.
*/

type AlertAcknowledgementService struct {
	alertRepository interfaces.AlertRepository
	lockService     interfaces.LockService
}

func NewAlertAcknowledgementService(alertRepository interfaces.AlertRepository, lockService interfaces.LockService) AlertAcknowledgementService {
	return AlertAcknowledgementService{
		alertRepository: alertRepository,
		lockService:     lockService,
	}
}

func (service *AlertAcknowledgementService) AcknowledgeAlarm(alertID string) error {
	alert, err := service.alertRepository.GetByID(alertID)

	//error getting alert or not found alert
	if err != nil {
		return err
	}

	//the alert is already acknowledged, finish the process.
	if alert.Status == constants.AlertStatusAcknowledge {
		return nil
	}

	lockAlarm := service.lockService.Lock(alert.ID)
	defer service.lockService.Unlock(alert.ID)

	if !lockAlarm {
		return constants.ErrorLockedResource
	}

	alert.Acknowledge()
	err = service.alertRepository.Save(alert)
	if err != nil {
		return err
	}

	return nil
}
