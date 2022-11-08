package mock

import "github.com/stretchr/testify/mock"

type CustomValidatorMock struct {
	mock.Mock
}

func (m *CustomValidatorMock) Validate(i interface{}) error {
	args := m.Called()
	return args.Error(0)
}
