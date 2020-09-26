package entities

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/santi8ago8/pager-service/app/src/constants"
)

func TestMonitoredService_SetStatus(t *testing.T) {

	t.Run("Healthy", func(t *testing.T) {
		monitoredService := MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		monitoredService.SetHealthy()

		assert.Equal(t, monitoredService.Status, constants.ServiceStatusHealthy)
	})

	t.Run("Unhealthy", func(t *testing.T) {
		monitoredService := MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
		monitoredService.SetUnhealthy()

		assert.Equal(t, monitoredService.Status, constants.ServiceStatusUnhealthy)
	})

}
