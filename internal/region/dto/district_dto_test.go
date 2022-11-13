package dto

import (
	"testing"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestDistrictResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    *model.District
		Expected districtResponse
	}{
		{
			Name: "all filled",
			Model: &model.District{
				ID:        1,
				Name:      "districtA",
				RegencyID: 1,
				Regency: model.Regency{
					ID:         1,
					Name:       "regencyA",
					ProvinceID: 1,
				},
			},
			Expected: districtResponse{
				ID:          1,
				Name:        "districtA",
				RegencyID:   1,
				RegencyName: "regencyA",
			},
		},
		{
			Name: "some filled",
			Model: &model.District{
				ID:   1,
				Name: "districtA",
			},
			Expected: districtResponse{
				ID:   1,
				Name: "districtA",
			},
		},
		{
			Name:     "dto empty",
			Model:    &model.District{},
			Expected: districtResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result districtResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestDistrictsResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    []model.District
		Expected DistrictsResponse
	}{
		{
			Name: "all filled",
			Model: []model.District{{
				ID:        1,
				Name:      "districtA",
				RegencyID: 1,
				Regency: model.Regency{
					ID:         1,
					Name:       "regencyA",
					ProvinceID: 1,
				},
			}},
			Expected: DistrictsResponse{{
				ID:          1,
				Name:        "districtA",
				RegencyID:   1,
				RegencyName: "regencyA",
			}},
		},
		{
			Name: "some filled",
			Model: []model.District{{
				ID:   1,
				Name: "districtA",
			}},
			Expected: DistrictsResponse{{
				ID:   1,
				Name: "districtA",
			}},
		},
		{
			Name:     "dto empty",
			Model:    []model.District{},
			Expected: DistrictsResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result DistrictsResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
