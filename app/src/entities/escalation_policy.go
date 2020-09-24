package entities

import "github.com/santi8ago8/pager-service/app/src/constants"

type EscalationPolicy struct {
	ID                 string   `json:"id"`
	MonitoredServiceID string   `json:"monitored_service_id"`
	Levels             []*Level `json:"levels"`
}

type Level struct {
	ID      int       `json:"id"`
	Targets []*Target `json:"targets"`
}

type Target struct {
	ID          string               `json:"id"`
	Type        constants.TargetType `json:"type"`
	Email       string               `json:"email,omitempty"`
	PhoneNumber string               `json:"phone_number,omitempty"`
}

func (policy *EscalationPolicy) GetLevelToNotify(levelNumber int) *Level {

	for _, level := range policy.Levels {
		if level.ID == levelNumber {
			return level
		}
	}
	return nil
}
