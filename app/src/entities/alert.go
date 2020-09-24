package entities

import (
	"github.com/google/uuid"
	"github.com/santi8ago8/pager-service/app/src/constants"
)

type Alert struct {
	ID                 string                `json:"id"`
	MonitoredServiceID string                `json:"monitored_service_id"`
	Message            string                `json:"message"`
	NotifiedLevelID    *int                  `json:"notified_level_id"`
	Status             constants.AlertStatus `json:"status"`
}

func NewAlert(serviceID string, alertMessage string) *Alert {
	IDv4, _ := uuid.NewRandom()
	return &Alert{
		ID:                 IDv4.String(),
		MonitoredServiceID: serviceID,
		Message:            alertMessage,
		Status:             constants.AlertStatusOpen,
	}
}

func (alert *Alert) Acknowledge() {
	alert.Status = constants.AlertStatusAcknowledge
}

func (alert *Alert) GetNextLevelToNotify() int {
	if alert.NotifiedLevelID == nil {
		return 0
	}
	return *alert.NotifiedLevelID + 1
}
