package mocks

type EmailServiceMock struct {
}

func (service *EmailServiceMock) Send(mailAddress string, text string) {
	//do something
}
