package qrcode

import (
	"github.com/skip2/go-qrcode"
)

type QRCode struct {
}

func (*QRCode) GenerateQRCode(hashCode string) ([]byte, error) {
	var qrCode []byte
	qrCode, err := qrcode.Encode(hashCode, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}
