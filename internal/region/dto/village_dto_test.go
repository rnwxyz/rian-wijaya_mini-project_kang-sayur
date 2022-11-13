package dto

import (
	"testing"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestVillageResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    *model.Village
		Expected villageResponse
	}{
		{
			Name: "all filled",
			Model: &model.Village{
				ID:         1,
				Name:       "villageA",
				DistrictID: 1,
				District: model.District{
					ID:        1,
					Name:      "districtA",
					RegencyID: 1,
				},
			},
			Expected: villageResponse{
				ID:           1,
				Name:         "villageA",
				DistrictID:   1,
				DistrictName: "districtA",
			},
		},
		{
			Name: "some filled",
			Model: &model.Village{
				ID:   1,
				Name: "villageA",
			},
			Expected: villageResponse{
				ID:   1,
				Name: "villageA",
			},
		},
		{
			Name:     "dto empty",
			Model:    &model.Village{},
			Expected: villageResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result villageResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestVillagesResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    []model.Village
		Expected VillagesResponse
	}{
		{
			Name: "all filled",
			Model: []model.Village{{
				ID:         1,
				Name:       "villageA",
				DistrictID: 1,
				District: model.District{
					ID:        1,
					Name:      "districtA",
					RegencyID: 1,
				},
			}},
			Expected: VillagesResponse{{
				ID:           1,
				Name:         "villageA",
				DistrictID:   1,
				DistrictName: "districtA",
			}},
		},
		{
			Name: "some filled",
			Model: []model.Village{{
				ID:   1,
				Name: "villageA",
			}},
			Expected: VillagesResponse{{
				ID:   1,
				Name: "villageA",
			}},
		},
		{
			Name:     "dto empty",
			Model:    []model.Village{},
			Expected: VillagesResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result VillagesResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
