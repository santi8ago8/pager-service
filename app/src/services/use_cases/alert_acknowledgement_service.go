package use_cases

import (
	"github.com/santi8ago8/pager-service/app/src/interfaces"
)

type AlertAcknowledgementService struct {
	alertRepository interfaces.AlertRepository
}

func NewAlertAcknowledgementService(alertRepository interfaces.AlertRepository) *AlertAcknowledgementService {
	return &AlertAcknowledgementService{
		alertRepository: alertRepository,
	}
}

func (service *AlertAcknowledgementService) AcknowledgeAlarm(alertID string) error {
	alert, err := service.alertRepository.GetByID(alertID)

	//error getting alert or not found alert
	if err != nil {
		return err
	}

	alert.Acknowledge()
	err = service.alertRepository.Save(alert)
	if err != nil {
		return err
	}

	return nil
}
