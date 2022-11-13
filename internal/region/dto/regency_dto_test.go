package dto

import (
	"testing"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestRegencyResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    *model.Regency
		Expected regencyResponse
	}{
		{
			Name: "all filled",
			Model: &model.Regency{
				ID:         1,
				Name:       "regencyA",
				ProvinceID: 1,
				Province: model.Province{
					ID:   1,
					Name: "provinceA",
				},
			},
			Expected: regencyResponse{
				ID:           1,
				Name:         "regencyA",
				ProvinceID:   1,
				ProvinceName: "provinceA",
			},
		},
		{
			Name: "some filled",
			Model: &model.Regency{
				ID:   1,
				Name: "regencyA",
			},
			Expected: regencyResponse{
				ID:   1,
				Name: "regencyA",
			},
		},
		{
			Name:     "dto empty",
			Model:    &model.Regency{},
			Expected: regencyResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result regencyResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestRegenciesResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    []model.Regency
		Expected RegenciesResponse
	}{
		{
			Name: "all filled",
			Model: []model.Regency{{
				ID:         1,
				Name:       "regencyA",
				ProvinceID: 1,
				Province: model.Province{
					ID:   1,
					Name: "provinceA",
				},
			}},
			Expected: RegenciesResponse{{
				ID:           1,
				Name:         "regencyA",
				ProvinceID:   1,
				ProvinceName: "provinceA",
			}},
		},
		{
			Name: "some filled",
			Model: []model.Regency{{
				ID:   1,
				Name: "regencyA",
			}},
			Expected: RegenciesResponse{{
				ID:   1,
				Name: "regencyA",
			}},
		},
		{
			Name:     "dto empty",
			Model:    []model.Regency{},
			Expected: RegenciesResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result RegenciesResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
