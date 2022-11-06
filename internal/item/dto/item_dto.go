package dto

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type ItemRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Qty         int    `json:"qty"`
	Price       int    `json:"price"`
	CategoryID  uint   `json:"category_id" validate:"required"`
}

func (u *ItemRequest) ToModel() *model.Item {
	return &model.Item{
		Name:        u.Name,
		Description: u.Description,
		Qty:         u.Qty,
		Price:       u.Price,
		CategoryID:  u.CategoryID,
	}
}

type ItemResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Qty          int    `json:"qty"`
	Price        int    `json:"price"`
	CategoryName string `json:"category_name"`
}

func (u *ItemResponse) FromModel(model *model.Item) {
	u.ID = model.ID
	u.Name = model.Name
	u.Description = model.Description
	u.Qty = model.Qty
	u.Price = model.Price
	u.CategoryName = model.Category.Name
}

type ItemsResponse []ItemResponse

func (u *ItemsResponse) FromModel(model []model.Item) {
	for _, each := range model {
		var item ItemResponse
		item.FromModel(&each)
		*u = append(*u, item)
	}
}
