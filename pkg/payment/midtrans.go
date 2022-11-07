package payment

import (
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type Midtrans struct {
}

func (*Midtrans) NewTransaction(order model.Order) string {
	var s snap.Client
	s.New("SB-Mid-server-E3oJluTQ3CHxYL19HG8fyJAK", midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.ID.String(),
			GrossAmt: int64(order.GrandTotal),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.User.Name,
			Email: order.User.Email,
		},
		EnabledPayments: snap.AllSnapPaymentType,
	}
	resp, err := s.CreateTransactionUrl(req)
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	return resp
}
