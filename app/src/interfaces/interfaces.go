package interfaces

import (
	"time"

	"github.com/santi8ago8/pager-service/app/src/entities"
)

//Repositories

type EscalationPolicyRepository interface {
	GetByServiceID(serviceID string) (*entities.EscalationPolicy, error)
	Save(*entities.EscalationPolicy) error
}

type AlertRepository interface {
	GetByID(alertID string) (*entities.Alert, error)
	Save(*entities.Alert) error
}

type MonitoredServiceRepository interface {
	GetByID(ID string) (*entities.MonitoredService, error)
	Save(*entities.MonitoredService) error
}

//Services
type NotificationService interface {
	SendNotifications(alert *entities.Alert, targets []*entities.Target)
}

type EmailService interface {
	Send(mailAddress string, text string)
}

type SmsService interface {
	Send(phoneNumber string, text string)
}

type TimerService interface {
	Schedule(resource string, resourceID string, time time.Time)
}
