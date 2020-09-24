package base

import (
	"github.com/santi8ago8/pager-service/app/src/constants"
	"github.com/santi8ago8/pager-service/app/src/entities"
	"github.com/santi8ago8/pager-service/app/src/interfaces"
)

type NotificationService struct {
	smsService   interfaces.SmsService
	emailService interfaces.EmailService
}

func NewNotificationService(smsService interfaces.SmsService, emailService interfaces.EmailService) *NotificationService {
	return &NotificationService{
		smsService:   smsService,
		emailService: emailService,
	}
}

func (service *NotificationService) SendNotifications(alert *entities.Alert, targets []*entities.Target) {
	for _, target := range targets {
		switch target.Type {
		case constants.TargetTypeEmail:
			service.emailService.Send(target.Email, alert.Message)
		case constants.TargetTypeSms:
			service.smsService.Send(target.PhoneNumber, alert.Message)
		}
	}
}
