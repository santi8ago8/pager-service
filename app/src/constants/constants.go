package constants

//ServiceStatus Service status types type
type ServiceStatus string

const (
	//ServiceStatusHealthy service is healthy
	ServiceStatusHealthy ServiceStatus = "healthy"
	//ServiceStatusUnhealthy service is unhealthy
	ServiceStatusUnhealthy ServiceStatus = "unhealthy"
)

//TargetType Target types type
type TargetType string

const (
	//TargetTypeSms Target using SMS
	TargetTypeSms TargetType = "SMS"
	//TargetTypeEmail Target using Email
	TargetTypeEmail TargetType = "email"
)

//AlertStatus Alert status types type
type AlertStatus string

const (
	//AlertStatusOpen open alert
	AlertStatusOpen AlertStatus = "open"
	//AlertStatusAcknowledge acknowledge alert
	AlertStatusAcknowledge AlertStatus = "acknowledge"
	//AlertStatusClosed closed alert
	AlertStatusClosed AlertStatus = "closed"
)
