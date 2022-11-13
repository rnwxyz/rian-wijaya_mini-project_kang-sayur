package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type OrderRequest struct {
	CheckpointID string              `json:"checkpoint_id" validate:"required"`
	Order        OrderDetailsRequest `json:"order" validate:"required"`
}
type OrderResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	UserName        string    `json:"user_name"`
	CheckpointID    uuid.UUID `json:"checkpoint_id"`
	CheckpointName  string    `json:"checkpoint_name"`
	CreatedAt       time.Time `json:"created_at"`
	StatusOrderName string    `json:"status_order"`
	ShippingCost    int       `json:"shipping_cost"`
	TotalPrice      int       `json:"total_price"`
	GrandTotal      int       `json:"grand_total"`
	ExpiredOrder    time.Time `json:"expired_order"`
}

func (u *OrderResponse) FromModel(model *model.Order) {
	u.ID = model.ID
	u.UserID = model.UserID
	u.UserName = model.User.Name
	u.CheckpointID = model.CheckpointID
	u.CheckpointName = model.Checkpoint.Name
	u.CreatedAt = model.CreatedAt
	u.StatusOrderName = model.StatusOrder.Name
	u.ShippingCost = model.ShippingCost
	u.TotalPrice = model.TotalPrice
	u.GrandTotal = model.GrandTotal
	u.ExpiredOrder = model.ExpiredOrder
}

type OrdersResponse []OrderResponse

func (u *OrdersResponse) FromModel(model []model.Order) {
	for _, each := range model {
		var order OrderResponse
		order.FromModel(&each)
		*u = append(*u, order)
	}
}

type OrderWithDetailResponse struct {
	ID              uuid.UUID            `json:"id"`
	CheckpointID    uuid.UUID            `json:"checkpoint_id"`
	CheckpointName  string               `json:"checkpoint_name"`
	CreatedAt       time.Time            `json:"created_at"`
	StatusOrderName string               `json:"status_order"`
	ShippingCost    int                  `json:"shipping_cost"`
	TotalPrice      int                  `json:"total_price"`
	GrandTotal      int                  `json:"grand_total"`
	Hash            string               `json:"code"`
	ExpiredOrder    time.Time            `json:"expired_order"`
	OrderDetail     OrderDetailsResponse `json:"order_detail"`
}

func (u *OrderWithDetailResponse) FromModel(model *model.Order) {
	var od OrderDetailsResponse
	od.FromModel(model.OrderDetail)

	u.ID = model.ID
	u.CheckpointID = model.CheckpointID
	u.CheckpointName = model.Checkpoint.Name
	u.CreatedAt = model.CreatedAt
	u.StatusOrderName = model.StatusOrder.Name
	u.ShippingCost = model.ShippingCost
	u.TotalPrice = model.TotalPrice
	u.GrandTotal = model.GrandTotal
	u.Hash = model.Hash
	u.ExpiredOrder = model.ExpiredOrder
	u.OrderDetail = od
}

type TakeOrder struct {
	CheckpointID string `json:"checkpoint_id"`
	Code         string `json:"code" validate:"required"`
}

type NewOrder struct {
	OrderID     uuid.UUID `json:"order_id"`
	RedirectURL string    `json:"redirect_url"`
}
