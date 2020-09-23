package entities

import "github.com/santi8ago8/pager-service/app/src/constants"

type EscalationPolicy struct {
	ID                 string `json:"id"`
	MonitoredServiceID string `json:"monitored_service_id"`
	Levels             []struct {
		ID      int `json:"id"`
		Targets []struct {
			ID          string               `json:"id"`
			Type        constants.TargetType `json:"type"`
			Email       string               `json:"email,omitempty"`
			PhoneNumber string               `json:"phone_number,omitempty"`
		} `json:"targets"`
	} `json:"levels"`
}
