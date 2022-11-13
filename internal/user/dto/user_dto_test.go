package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestUserSignup_ToModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Dto      UserSignup
		Expected *model.User
	}{
		{
			Name: "all filled",
			Dto: UserSignup{
				Name:     "test",
				Email:    "test@gmail.com",
				Password: "test",
			},
			Expected: &model.User{
				Name:  "test",
				Email: "test@gmail.com",
			},
		},
		{
			Name: "some filled",
			Dto: UserSignup{
				Email: "test@gmail.com",
			},
			Expected: &model.User{
				Email: "test@gmail.com",
			},
		},
		{
			Name:     "empty",
			Dto:      UserSignup{},
			Expected: &model.User{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			result := v.Dto.ToModel()
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestUserUpdate_ToModel(t *testing.T) {
	regionId := uint(1)

	testCase := []struct {
		Name     string
		Dto      UserUpdate
		Expected *model.User
	}{
		{
			Name: "all filled",
			Dto: UserUpdate{
				Name:       "test",
				Phone:      "081234567890",
				ProvinceID: &regionId,
				RegencyID:  &regionId,
				DistrictID: &regionId,
				VillageID:  &regionId,
			},
			Expected: &model.User{
				Name:       "test",
				Phone:      "081234567890",
				ProvinceID: &regionId,
				RegencyID:  &regionId,
				DistrictID: &regionId,
				VillageID:  &regionId,
			},
		},
		{
			Name: "some filled",
			Dto: UserUpdate{
				Phone:      "081234567890",
				ProvinceID: &regionId,
			},
			Expected: &model.User{
				Phone:      "081234567890",
				ProvinceID: &regionId,
			},
		},
		{
			Name:     "empty",
			Dto:      UserUpdate{},
			Expected: &model.User{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			result := v.Dto.ToModel()
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestUserResponse_FromModel(t *testing.T) {
	regionId := uint(1)
	userId := uuid.New()

	testCase := []struct {
		Name     string
		Model    *model.User
		Expected UserResponse
	}{
		{
			Name: "all filled",
			Model: &model.User{
				ID:         userId,
				Name:       "test",
				Email:      "test@gmail.com",
				Phone:      "081234567890",
				ProvinceID: &regionId,
				Province: model.Province{
					ID:   regionId,
					Name: "provinceTest",
				},
				RegencyID: &regionId,
				Regency: model.Regency{
					ID:   regionId,
					Name: "regencyTest",
				},
				DistrictID: &regionId,
				District: model.District{
					ID:   regionId,
					Name: "districtTest",
				},
				VillageID: &regionId,
				Village: model.Village{
					ID:   regionId,
					Name: "villageTest",
				},
			},
			Expected: UserResponse{
				ID:           userId,
				Name:         "test",
				Email:        "test@gmail.com",
				Phone:        "081234567890",
				ProvinceName: "provinceTest",
				RegencyName:  "regencyTest",
				DistrictName: "districtTest",
				VillageName:  "villageTest",
			},
		},
		{
			Name: "some filled",
			Model: &model.User{
				ID:         userId,
				Name:       "test",
				Email:      "test@gmail.com",
				ProvinceID: &regionId,
				Province: model.Province{
					ID:   regionId,
					Name: "provinceTest",
				},
			},
			Expected: UserResponse{
				ID:           userId,
				Name:         "test",
				Email:        "test@gmail.com",
				ProvinceName: "provinceTest",
			},
		},
		{
			Name:     "empty",
			Model:    &model.User{},
			Expected: UserResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result UserResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestUsersResponse_FromModel(t *testing.T) {
	regionId := uint(1)
	userId := uuid.New()

	testCase := []struct {
		Name     string
		Model    []model.User
		Expected UsersResponse
	}{
		{
			Name: "all filled",
			Model: []model.User{
				{
					ID:         userId,
					Name:       "test",
					Email:      "test@gmail.com",
					Phone:      "081234567890",
					ProvinceID: &regionId,
					Province: model.Province{
						ID:   regionId,
						Name: "provinceTest",
					},
					RegencyID: &regionId,
					Regency: model.Regency{
						ID:   regionId,
						Name: "regencyTest",
					},
					DistrictID: &regionId,
					District: model.District{
						ID:   regionId,
						Name: "districtTest",
					},
					VillageID: &regionId,
					Village: model.Village{
						ID:   regionId,
						Name: "villageTest",
					},
				},
			},
			Expected: UsersResponse{
				{
					ID:           userId,
					Name:         "test",
					Email:        "test@gmail.com",
					Phone:        "081234567890",
					ProvinceName: "provinceTest",
					RegencyName:  "regencyTest",
					DistrictName: "districtTest",
					VillageName:  "villageTest",
				},
			},
		},
		{
			Name: "some filled",
			Model: []model.User{
				{
					ID:         userId,
					Name:       "test",
					Email:      "test@gmail.com",
					ProvinceID: &regionId,
					Province: model.Province{
						ID:   regionId,
						Name: "provinceTest",
					},
				},
			},
			Expected: UsersResponse{
				{
					ID:           userId,
					Name:         "test",
					Email:        "test@gmail.com",
					ProvinceName: "provinceTest",
				},
			},
		},
		{
			Name:     "empty",
			Model:    []model.User{},
			Expected: UsersResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result UsersResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
