package entities

import "github.com/santi8ago8/pager-service/app/src/constants"

type MonitoredService struct {
	ID     string                  `json:"id"`
	Name   string                  `json:"name"`
	Status constants.ServiceStatus `json:"status"`
}



func (monitoredService *MonitoredService) SetUnhealthy() {
	monitoredService.Status = constants.ServiceStatusUnhealthy
}

func (monitoredService *MonitoredService) SetHealthy() {
	monitoredService.Status = constants.ServiceStatusHealthy
}
