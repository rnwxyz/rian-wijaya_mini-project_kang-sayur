package mock

import "github.com/stretchr/testify/mock"

type QRCodeMock struct {
	mock.Mock
}

func (b *QRCodeMock) GenerateQRCode(hashCode string) ([]byte, error) {
	args := b.Called()
	return args.Get(0).([]byte), args.Error(1)
}
