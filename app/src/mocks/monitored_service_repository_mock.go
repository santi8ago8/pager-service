package mocks

import (
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
)

type MonitoredServiceRepositoryMock struct {
	monitoredServices map[string]*entities.MonitoredService
	saveErrors        map[string]error
}

func NewMonitoredServiceRepositoryMock() *MonitoredServiceRepositoryMock {
	return &MonitoredServiceRepositoryMock{
		monitoredServices: map[string]*entities.MonitoredService{},
		saveErrors:        map[string]error{},
	}
}

func (service *MonitoredServiceRepositoryMock) GetByID(serviceID string) (*entities.MonitoredService, error) {
	mService := service.monitoredServices[serviceID]
	if mService == nil {
		return nil, constants.ErrorMonitoredServiceNotFound
	}
	return mService, nil
}
func (service *MonitoredServiceRepositoryMock) Save(monitoredService *entities.MonitoredService) error {
	if service.saveErrors[monitoredService.ID] != nil {
		return service.saveErrors[monitoredService.ID]
	}
	service.monitoredServices[monitoredService.ID] = monitoredService
	return nil
}
func (service *MonitoredServiceRepositoryMock) AddSaveError(serviceID string, e error) {
	service.saveErrors[serviceID] = e
}
