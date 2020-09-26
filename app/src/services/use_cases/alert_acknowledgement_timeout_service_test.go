package use_cases

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
	"github.com/santi8ago8/pager-service/app/src/mocks"
	services_base "github.com/santi8ago8/pager-service/app/src/services/base"
)

var serviceAckTimeout AlertAcknowledgementTimeoutService
var lockService *mocks.LockServiceMock

func initializeTestAckTimeout() {
	alertRepository = mocks.NewAlertRepositoryMock()
	monitoredServiceRepository = mocks.NewMonitoredServiceRepositoryMock()
	escalationPolicyRepository = mocks.NewEscalationPolicyRepositoryMock()
	smsService := &mocks.SmsServiceMock{}
	emailService := &mocks.EmailServiceMock{}
	notificationService = services_base.NewNotificationService(smsService, emailService)
	timerService = mocks.NewTimerServiceMock()
	lockService = mocks.NewLockServiceMock()

	serviceAckTimeout = NewAlertAcknowledgementTimeoutService(alertRepository, monitoredServiceRepository, escalationPolicyRepository, notificationService, timerService, lockService)
}

func TestAlertAcknowledgementTimeoutService_AcknowledgementTimeout(t *testing.T) {
	initializeTestAckTimeout()
	t.Run("Happy path acknowledgement timeout", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		alert := entities.NewAlert(monService.ID, "Error rate > 10% in the last 10 minutes")
		level := 0
		alert.NotifiedLevelID = &level
		alertRepository.Save(alert)
		monitoredServiceRepository.Save(&monService)
		escPolicy := entities.EscalationPolicy{
			ID:                 uuid.New().String(),
			MonitoredServiceID: monService.ID,
			Levels: []*entities.Level{
				{
					ID: 0,
					Targets: []*entities.Target{
						{
							ID:          uuid.New().String(),
							Type:        constants.TargetTypeSms,
							PhoneNumber: "+54923959292",
						},
					},
				},
				{
					ID: 1,
					Targets: []*entities.Target{
						{
							ID:          uuid.New().String(),
							Type:        constants.TargetTypeSms,
							PhoneNumber: "+3493223463",
						},
						{
							ID:    uuid.New().String(),
							Type:  constants.TargetTypeEmail,
							Email: "test@test.com",
						},
					},
				},
			},
		}
		escalationPolicyRepository.Save(&escPolicy)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, nil)

	})

	t.Run("alert not found", func(t *testing.T) {

		err := serviceAckTimeout.AcknowledgementTimeout("FAKEUUID_ALERT_ID")

		assert.Equal(t, err, constants.ErrorAlertNotFound)

	})

	t.Run("alert previously acknowledge", func(t *testing.T) {

		alert := entities.NewAlert("FAKEUUID_MONITORED_SERVICE", "Error rate > 10% in the last 10 minutes")
		alert.Acknowledge()
		alertRepository.Save(alert)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, nil)

	})

	t.Run("Healthy service", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
		alert := entities.NewAlert(monService.ID, "Error rate > 10% in the last 10 minutes")
		alertRepository.Save(alert)
		monitoredServiceRepository.Save(&monService)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, nil)

	})

	t.Run("Monitored service not found", func(t *testing.T) {
		alert := entities.NewAlert("FAKEUUID_MONITORED_SERVICE", "Error rate > 10% in the last 10 minutes")
		alertRepository.Save(alert)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, constants.ErrorMonitoredServiceNotFound)

	})

	t.Run("Can't lock resources", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		alert := entities.NewAlert(monService.ID, "Error rate > 10% in the last 10 minutes")
		alertRepository.Save(alert)
		monitoredServiceRepository.Save(&monService)
		lockService.SetMockResponse(false)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, constants.ErrorLockedResource)
		lockService.SetMockResponse(true)

	})

	t.Run("Escalation policy not found", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		alert := entities.NewAlert(monService.ID, "Error rate > 10% in the last 10 minutes")
		level := 0
		alert.NotifiedLevelID = &level
		alertRepository.Save(alert)
		monitoredServiceRepository.Save(&monService)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, constants.ErrorEscalationPolicyNotFound)

	})

	t.Run("Escalation policy with one level", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		alert := entities.NewAlert(monService.ID, "Error rate > 10% in the last 10 minutes")
		level := 0
		alert.NotifiedLevelID = &level
		alertRepository.Save(alert)
		monitoredServiceRepository.Save(&monService)
		escPolicy := entities.EscalationPolicy{
			ID:                 uuid.New().String(),
			MonitoredServiceID: monService.ID,
			Levels: []*entities.Level{
				{
					ID: 0,
					Targets: []*entities.Target{
						{
							ID:          uuid.New().String(),
							Type:        constants.TargetTypeSms,
							PhoneNumber: "+54923959292",
						},
					},
				},
			},
		}
		escalationPolicyRepository.Save(&escPolicy)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, nil)

	})

	t.Run("Fail save alert", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		alert := entities.NewAlert(monService.ID, "Error rate > 10% in the last 10 minutes")
		level := 0
		alert.NotifiedLevelID = &level
		alertRepository.Save(alert)
		monitoredServiceRepository.Save(&monService)
		escPolicy := entities.EscalationPolicy{
			ID:                 uuid.New().String(),
			MonitoredServiceID: monService.ID,
			Levels: []*entities.Level{
				{
					ID: 0,
					Targets: []*entities.Target{
						{
							ID:          uuid.New().String(),
							Type:        constants.TargetTypeSms,
							PhoneNumber: "+54923959292",
						},
					},
				},
				{
					ID: 1,
					Targets: []*entities.Target{
						{
							ID:          uuid.New().String(),
							Type:        constants.TargetTypeSms,
							PhoneNumber: "+3493223463",
						},
						{
							ID:    uuid.New().String(),
							Type:  constants.TargetTypeEmail,
							Email: "test@test.com",
						},
					},
				},
			},
		}
		escalationPolicyRepository.Save(&escPolicy)
		errorSave := errors.New("Optimistic locking")
		alertRepository.AddSaveError(alert.ID, errorSave)

		err := serviceAckTimeout.AcknowledgementTimeout(alert.ID)

		assert.Equal(t, err, errorSave)

	})
}
