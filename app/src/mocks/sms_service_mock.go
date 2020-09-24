package mocks

type SmsServiceMock struct {
}

func (service *SmsServiceMock) Send(phoneNumber string, text string) {
	//do something
}
