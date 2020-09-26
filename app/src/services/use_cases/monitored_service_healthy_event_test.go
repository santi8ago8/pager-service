package use_cases

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
	"github.com/santi8ago8/pager-service/app/src/mocks"
)

var serviceMonServiceHealthyEvent MonitoredServiceHealthyEvent

func initializeTestMonServiceHealthyEvent() {
	monitoredServiceRepository = mocks.NewMonitoredServiceRepositoryMock()
	lockService = mocks.NewLockServiceMock()

	serviceMonServiceHealthyEvent = NewMonitoredServiceHealthyEvent(monitoredServiceRepository, lockService)
}

func TestMonitoredServiceHealthyEvent_ServiceHealthy(t *testing.T) {
	initializeTestMonServiceHealthyEvent()
	t.Run("Happy path healthy event", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		monitoredServiceRepository.Save(&monService)

		err := serviceMonServiceHealthyEvent.ServiceHealthy(monService.ID)

		assert.Equal(t, err, nil)
		assert.Equal(t, monService.Status, constants.ServiceStatusHealthy)

	})

	t.Run("Monitored service not found", func(t *testing.T) {
		err := serviceMonServiceHealthyEvent.ServiceHealthy("FAKE_MONITORED_SERVICE_UUID_1234")

		assert.Equal(t, err, constants.ErrorMonitoredServiceNotFound)

	})

	t.Run("Monitored service is already healthy", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
		monitoredServiceRepository.Save(&monService)

		err := serviceMonServiceHealthyEvent.ServiceHealthy(monService.ID)

		assert.Equal(t, err, nil)

	})

	t.Run("Monitored service locked", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		lockService.SetMockResponse(false)
		monitoredServiceRepository.Save(&monService)

		err := serviceMonServiceHealthyEvent.ServiceHealthy(monService.ID)

		assert.Equal(t, err, constants.ErrorLockedResource)
		lockService.SetMockResponse(true)

	})

	t.Run("Monitored service save error", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		monitoredServiceRepository.Save(&monService)
		errorSave := errors.New("Optimistic locking")
		monitoredServiceRepository.AddSaveError(monService.ID, errorSave)

		err := serviceMonServiceHealthyEvent.ServiceHealthy(monService.ID)

		assert.Equal(t, err, errorSave)

	})
}
