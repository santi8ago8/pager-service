package entities

import "github.com/santi8ago8/pager-service/app/src/constants"

type Alert struct {
	ID                 string                `json:"id"`
	MonitoredServiceID string                `json:"monitored_service_id"`
	Message            string                `json:"message"`
	NotifiedLevelID    *int                  `json:"notified_level_id"`
	Status             constants.AlertStatus `json:"status"`
}
