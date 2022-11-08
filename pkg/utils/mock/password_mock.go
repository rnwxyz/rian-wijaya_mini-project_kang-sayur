package mock

import "github.com/stretchr/testify/mock"

type PasswordMock struct {
	mock.Mock
}

func (m *PasswordMock) HashPassword(password string) (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}

func (m *PasswordMock) CheckPasswordHash(password, hash string) bool {
	args := m.Called()
	return args.Get(0).(bool)
}
