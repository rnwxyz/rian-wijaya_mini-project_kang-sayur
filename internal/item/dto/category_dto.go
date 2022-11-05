package dto

import (
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type CategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (u *CategoryRequest) ToModel() *model.Category {
	return &model.Category{
		Name:        u.Name,
		Description: u.Description,
	}
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (u *CategoryResponse) FromModel(model *model.Category) {
	u.ID = model.ID
	u.Name = model.Name
	u.Description = model.Description
}

type CategoriesResponse []CategoryResponse

func (u *CategoriesResponse) FromModel(model []model.Category) {
	for _, each := range model {
		var categorie CategoryResponse
		categorie.FromModel(&each)
		*u = append(*u, categorie)
	}
}
