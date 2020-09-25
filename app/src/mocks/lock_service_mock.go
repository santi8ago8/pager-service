package mocks

type LockServiceMock struct {
	defaultValue bool
}

func NewLockServiceMock() *LockServiceMock {
	return &LockServiceMock{
		defaultValue: true,
	}
}

//SetMockResponse change lock mock response.
func (service *LockServiceMock) SetMockResponse(value bool) {
	service.defaultValue = value
}

//Unlock best effort unlock
func (service *LockServiceMock) Unlock(resourceID string) {
	//do something
}

//Lock returns true if the resource is locked correctly, and false if the resource was previously locked
func (service *LockServiceMock) Lock(resourceID string) bool {
	return service.defaultValue
}
