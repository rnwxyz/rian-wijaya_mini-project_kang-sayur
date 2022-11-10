package mock

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type MidtransMock struct {
	mock.Mock
}

func (m *MidtransMock) NewTransaction(order model.Order, user model.User) (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}
