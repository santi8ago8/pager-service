package mocks

import "time"

type TimerServiceMock struct {
}

func NewTimerServiceMock() *TimerServiceMock {
	return &TimerServiceMock{}
}

func (service *TimerServiceMock) Schedule(resource string, resourceID string, time time.Time) {
	//do something
}
