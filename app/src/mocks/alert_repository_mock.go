package mocks

import (
	"errors"

	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
)

type AlertRepositoryMock struct {
	alerts     map[string]*entities.Alert
	saveErrors map[string]error
	failAll    bool
}

func NewAlertRepositoryMock() *AlertRepositoryMock {
	return &AlertRepositoryMock{
		alerts:     map[string]*entities.Alert{},
		saveErrors: map[string]error{},
	}
}

func (service *AlertRepositoryMock) GetByID(alertID string) (*entities.Alert, error) {
	alert := service.alerts[alertID]
	if alert == nil {
		return nil, constants.ErrorAlertNotFound
	}
	return alert, nil
}
func (service *AlertRepositoryMock) Save(alert *entities.Alert) error {
	if service.saveErrors[alert.ID] != nil {
		return service.saveErrors[alert.ID]
	}
	if service.failAll {
		return errors.New("Fail all error")
	}
	service.alerts[alert.ID] = alert
	return nil
}
func (service *AlertRepositoryMock) AddSaveError(alertID string, e error) {
	service.saveErrors[alertID] = e
}

func (service *AlertRepositoryMock) ChangeSaveAllError(fail bool) {
	service.failAll = fail
}
