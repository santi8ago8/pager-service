package mocks

import (
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
)

type EscalationPolicyRepositoryMock struct {
	escalationPolicies map[string]*entities.EscalationPolicy
	saveErrors         map[string]error
}

func NewEscalationPolicyRepositoryMock() *EscalationPolicyRepositoryMock {
	return &EscalationPolicyRepositoryMock{
		escalationPolicies: map[string]*entities.EscalationPolicy{},
		saveErrors:         map[string]error{},
	}
}

func (service *EscalationPolicyRepositoryMock) GetByServiceID(serviceID string) (*entities.EscalationPolicy, error) {

	for _, escalation := range service.escalationPolicies {
		if escalation.MonitoredServiceID == serviceID {
			return escalation, nil
		}
	}
	return nil, constants.ErrorEscalationPolicyNotFound

}

func (service *EscalationPolicyRepositoryMock) Save(escalation *entities.EscalationPolicy) error {
	if service.saveErrors[escalation.ID] != nil {
		return service.saveErrors[escalation.ID]
	}
	service.escalationPolicies[escalation.ID] = escalation
	return nil
}

func (service *EscalationPolicyRepositoryMock) AddSaveError(escalationID string, e error) {
	service.saveErrors[escalationID] = e
}
