package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

// type OrderRequest struct {
// 	OrderDetail OrderDetailsRequest `json:"order_detail" validate:"required"`
// }

// func (u *OrderRequest) ToModel() *model.Order {
// 	return &model.Order{
// 		OrderDetail: *u.OrderDetail.ToModel(),
// 	}
// }

type OrderResponse struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	StatusOrderName string    `json:"status_order"`
	ShippingCost    int       `json:"shipping_cost"`
	TotalPrice      int       `json:"total_price"`
	GrandTotal      int       `json:"grand_total"`
	Code            string    `json:"code"`
	ExpiredOrder    time.Time `json:"expired_order"`
}

func (u *OrderResponse) FromModel(model *model.Order) {
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.StatusOrderName = model.StatusOrder.Name
	u.ShippingCost = model.ShippingCost
	u.TotalPrice = model.TotalPrice
	u.GrandTotal = model.GrandTotal
	u.Code = model.Code
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
	CreatedAt       time.Time            `json:"created_at"`
	StatusOrderName string               `json:"status_order"`
	ShippingCost    int                  `json:"shipping_cost"`
	TotalPrice      int                  `json:"total_price"`
	GrandTotal      int                  `json:"grand_total"`
	Code            string               `json:"code"`
	ExpiredOrder    time.Time            `json:"expired_order"`
	OrderDetail     OrderDetailsResponse `json:"order_detail"`
}

func (u *OrderWithDetailResponse) FromModel(model *model.Order) {
	var od OrderDetailsResponse
	od.FromModel(model.OrderDetail)

	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.StatusOrderName = model.StatusOrder.Name
	u.ShippingCost = model.ShippingCost
	u.TotalPrice = model.TotalPrice
	u.GrandTotal = model.GrandTotal
	u.Code = model.Code
	u.ExpiredOrder = model.ExpiredOrder
	u.OrderDetail = od
}
