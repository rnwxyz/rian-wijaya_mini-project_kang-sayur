package dto

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type OrderDetailRequest struct {
	ItemID uint `json:"item_id" validate:"gte=1,required"`
	Qty    int  `json:"qty" validate:"gte=1,required"`
	Total  int  `json:"total_price"`
}

type OrderDetailsRequest []OrderDetailRequest

func (u *OrderDetailRequest) ToModel() *model.OrderDetail {
	return &model.OrderDetail{
		ItemID: u.ItemID,
		Qty:    u.Qty,
		Total:  u.Total,
	}
}

func (u *OrderDetailsRequest) ToModel() *[]model.OrderDetail {
	var model []model.OrderDetail
	for _, each := range *u {
		dto := each.ToModel()
		model = append(model, *dto)
	}
	return &model
}

type OrderDetailResponse struct {
	ID       uint   `json:"order_detail_id"`
	ItemName string `json:"item_name"`
	Qty      int    `json:"qty"`
	Price    int    `json:"price"`
	Total    int    `json:"total"`
}

func (u *OrderDetailResponse) FromModel(model *model.OrderDetail) {
	u.ID = model.ID
	u.ItemName = model.Item.Name
	u.Qty = model.Qty
	u.Price = model.Price
	u.Total = model.Total
}

type OrderDetailsResponse []OrderDetailResponse

func (u *OrderDetailsResponse) FromModel(model []model.OrderDetail) {
	for _, each := range model {
		var orderDetail OrderDetailResponse
		orderDetail.FromModel(&each)
		*u = append(*u, orderDetail)
	}
}
