package use_cases

import (
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/interfaces"
)

/*
Given a Monitored Service in an Unhealthy State,
when the Pager receives a Healthy event related to this Monitored Service (monitored_service_healthy_event responsibility) (this file)
and later receives the Acknowledgement Timeout, (alert_acknowledgement_timeout_service responsibility)
then the Monitored Service becomes Healthy,
the Pager doesn’t notify any Target
and doesn’t set an acknowledgement delay
*/

type MonitoredServiceHealthyEvent struct {
	monitoredServiceRepository interfaces.MonitoredServiceRepository
	lockService                interfaces.LockService
}

func NewMonitoredServiceHealthyEvent(monitoredServiceRepository interfaces.MonitoredServiceRepository,
	lockService interfaces.LockService) MonitoredServiceHealthyEvent {
	return MonitoredServiceHealthyEvent{
		monitoredServiceRepository: monitoredServiceRepository,
		lockService:                lockService,
	}
}

func (service *MonitoredServiceHealthyEvent) ServiceHealthy(serviceID string) error {
	monService, err := service.monitoredServiceRepository.GetByID(serviceID)

	//error getting monitored service or not monitored service
	if err != nil {
		return err
	}

	//the service is already healthy, finish the process.
	if monService.Status == constants.ServiceStatusHealthy {
		return nil
	}

	lockMonService := service.lockService.Lock(monService.ID)
	defer service.lockService.Unlock(monService.ID)

	if !lockMonService {
		return constants.ErrorLockedResource
	}

	monService.SetHealthy()
	err = service.monitoredServiceRepository.Save(monService)
	if err != nil {
		return err
	}

	return nil
}
