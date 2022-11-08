package payment

import (
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type Midtrans struct {
}

func (*Midtrans) NewTransaction(order model.Order, user model.User) string {
	var s snap.Client
	s.New(config.Cfg.MIDTRANS_SERVER_KEY, midtrans.Sandbox)
	shipping := strconv.Itoa(constants.Shipping_cost)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.ID.String(),
			GrossAmt: int64(order.GrandTotal),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
		CustomField1:    "Shipping Cost : " + shipping,
		EnabledPayments: snap.AllSnapPaymentType,
	}
	resp, err := s.CreateTransactionUrl(req)
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	return resp
}
