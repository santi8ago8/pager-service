package base

import (
	"testing"

	"github.com/google/uuid"
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
	"github.com/santi8ago8/pager-service/app/src/mocks"
)

var service *NotificationService

func initializeTest() {
	smsService := &mocks.SmsServiceMock{}
	emailService := &mocks.EmailServiceMock{}

	service = NewNotificationService(smsService, emailService)
}

func TestNotificationService_SendNotifications(t *testing.T) {
	initializeTest()
	t.Run("Call Send Notifications", func(t *testing.T) {
		targets := []*entities.Target{
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
		}
		alert := entities.NewAlert("FAKE_UUID_MONITORED_SERVICE", "Error rate > 10% in the las 15 minutes.")

		service.SendNotifications(alert, targets)
	})
}
