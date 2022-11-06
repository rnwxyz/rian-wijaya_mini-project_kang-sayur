package constants

import (
	"time"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

// code konfirmation order to pick up order
const Item_code_length = 6

// set shipping cost
const Shipping_cost = 5000

// order expired duration
const ExpOrder = 24 * time.Hour

// status id
const Pending_status_order_id = 1
const Waiting_status_order_id = 2
const Ready_status_order_id = 3
const Success_status_order_id = 4
const Refund_status_order_id = 5
const Refund_success_status_order_id = 6
const Cencel_status_order_id = 7

// default value status order model
var (
	StatusOrder = []model.StatusOrder{
		{
			ID:          Pending_status_order_id,
			Name:        "pending",
			Description: "order waiting for payment",
		},
		{
			ID:          Waiting_status_order_id,
			Name:        "waiting",
			Description: "order send to checkpoint",
		},
		{
			ID:          Ready_status_order_id,
			Name:        "ready",
			Description: "order is ready in checkpoint. Take before expired",
		},
		{
			ID:          Success_status_order_id,
			Name:        "success",
			Description: "order success",
		},
		{
			ID:          Refund_status_order_id,
			Name:        "refund",
			Description: "order expired, refund 1/2 total cost",
		},
		{
			ID:          Refund_success_status_order_id,
			Name:        "refund_success",
			Description: "refund is success",
		},
		{
			ID:          Cencel_status_order_id,
			Name:        "cencel",
			Description: "order is cenceled or not paid",
		},
	}
)
