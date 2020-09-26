package use_cases

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
	"github.com/santi8ago8/pager-service/app/src/mocks"
)

var alertAcknowledgementService AlertAcknowledgementService

func initializeTestAcknowledgementService() {
	alertRepository = mocks.NewAlertRepositoryMock()
	lockService = mocks.NewLockServiceMock()

	alertAcknowledgementService = NewAlertAcknowledgementService(alertRepository, lockService)
}

func TestAlertAcknowledgementService_AcknowledgeAlarm(t *testing.T) {
	initializeTestAcknowledgementService()
	t.Run("Happy path acknowledge alarm", func(t *testing.T) {
		alert := entities.NewAlert("FAKE_MON_SERVICE_UUID", "Apdex < 0.7 in the last 5 minutes")
		alertRepository.Save(alert)

		err := alertAcknowledgementService.AcknowledgeAlarm(alert.ID)

		assert.Equal(t, err, nil)
		assert.Equal(t, alert.Status, constants.AlertStatusAcknowledge)

	})

	t.Run("Alert not found", func(t *testing.T) {
		err := alertAcknowledgementService.AcknowledgeAlarm("FAKE_ALERT_UUID")

		assert.Equal(t, err, constants.ErrorAlertNotFound)

	})

	t.Run("Alert is already acknowledged", func(t *testing.T) {
		alert := entities.NewAlert("FAKE_MON_SERVICE_UUID", "Apdex < 0.7 in the last 5 minutes")
		alert.Acknowledge()
		alertRepository.Save(alert)

		err := alertAcknowledgementService.AcknowledgeAlarm(alert.ID)

		assert.Equal(t, err, nil)

	})

	t.Run("Alert locked", func(t *testing.T) {
		alert := entities.NewAlert("FAKE_MON_SERVICE_UUID", "Apdex < 0.7 in the last 5 minutes")
		alertRepository.Save(alert)
		lockService.SetMockResponse(false)

		err := alertAcknowledgementService.AcknowledgeAlarm(alert.ID)

		assert.Equal(t, err, constants.ErrorLockedResource)
		lockService.SetMockResponse(true)

	})

	t.Run("Alarm save error", func(t *testing.T) {
		alert := entities.NewAlert("FAKE_MON_SERVICE_UUID", "Apdex < 0.7 in the last 5 minutes")
		alertRepository.Save(alert)
		errorSave := errors.New("Optimistic locking")
		alertRepository.AddSaveError(alert.ID, errorSave)

		err := alertAcknowledgementService.AcknowledgeAlarm(alert.ID)

		assert.Equal(t, err, errorSave)

	})

}
