package mock

import (
	"github.com/stretchr/testify/mock"
)

type ImportCsvMock struct {
	mock.Mock
}

func (b *ImportCsvMock) UnmarshalCsv(filepath string, model interface{}) error {
	args := b.Called()
	return args.Error(0)
}
