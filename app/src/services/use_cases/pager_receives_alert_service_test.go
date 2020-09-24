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

var service PagerReceiveAlertService
var alertRepository *mocks.AlertRepositoryMock
var monitoredServiceRepository *mocks.MonitoredServiceRepositoryMock
var escalationPolicyRepository *mocks.EscalationPolicyRepositoryMock
var timerService *mocks.TimerServiceMock
var notificationService *services_base.NotificationService

func initializeTest() {
	alertRepository = mocks.NewAlertRepositoryMock()
	monitoredServiceRepository = mocks.NewMonitoredServiceRepositoryMock()
	escalationPolicyRepository = mocks.NewEscalationPolicyRepositoryMock()
	smsService := &mocks.SmsServiceMock{}
	emailService := &mocks.EmailServiceMock{}
	notificationService = services_base.NewNotificationService(smsService, emailService)
	timerService = mocks.NewTimerServiceMock()

	service = NewPagerReceiveAlertService(alertRepository, monitoredServiceRepository, escalationPolicyRepository, notificationService, timerService)
}

func TestPagerReceiveAlertService_ReceivesAlert(t *testing.T) {
	initializeTest()
	t.Run("Happy path receives alert", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
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
					}},
			},
		}
		escalationPolicyRepository.Save(&escPolicy)

		alert, err := service.ReceivesAlert(monService.ID, "Error rate > 5% in the last 5 minutes")

		assert.Equal(t, err, nil)
		assert.Equal(t, monService.Status, constants.ServiceStatusUnhealthy)
		assert.NotEqual(t, alert, nil)
		assert.Equal(t, alert.Status, constants.AlertStatusOpen)
	})

	t.Run("Monitored service already unhealthy", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusUnhealthy,
		}
		monitoredServiceRepository.Save(&monService)

		alert, err := service.ReceivesAlert(monService.ID, "Error rate > 5% in the last 5 minutes")

		assert.Equal(t, err, nil)
		assert.Equal(t, alert, nil)
	})

	t.Run("Monitored service not found", func(t *testing.T) {

		alert, err := service.ReceivesAlert("FAKEUUID", "Error rate > 5% in the last 5 minutes")

		assert.Equal(t, err, constants.ErrorMonitoredServiceNotFound)
		assert.Equal(t, alert, nil)
	})

	t.Run("Escalation policy not found", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
		monitoredServiceRepository.Save(&monService)

		alert, err := service.ReceivesAlert(monService.ID, "Error rate > 5% in the last 5 minutes")

		assert.Equal(t, err, constants.ErrorEscalationPolicyNotFound)
		assert.Equal(t, alert, nil)
	})

	t.Run("Monitored service already unhealthy", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
		monitoredServiceRepository.Save(&monService)

		alert, err := service.ReceivesAlert(monService.ID, "Error rate > 5% in the last 5 minutes")

		assert.Equal(t, err, constants.ErrorEscalationPolicyNotFound)
		assert.Equal(t, alert, nil)
	})

	t.Run("Fail save monitored service", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
		errorSave := errors.New("Error in repository")
		monitoredServiceRepository.Save(&monService)
		monitoredServiceRepository.AddSaveError(monService.ID, errorSave)

		alert, err := service.ReceivesAlert(monService.ID, "Error rate > 5% in the last 5 minutes")

		assert.Equal(t, err, errorSave)
		assert.Equal(t, alert, nil)
	})

	t.Run("Fail save alert", func(t *testing.T) {
		monService := entities.MonitoredService{
			ID:     uuid.New().String(),
			Status: constants.ServiceStatusHealthy,
		}
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
					}},
			},
		}
		escalationPolicyRepository.Save(&escPolicy)
		alertRepository.ChangeSaveAllError(true)

		alert, err := service.ReceivesAlert(monService.ID, "Error rate > 5% in the last 5 minutes")

		alertRepository.ChangeSaveAllError(false)
		assert.Equal(t, err != nil, true)
		assert.Equal(t, alert, nil)
	})
}
