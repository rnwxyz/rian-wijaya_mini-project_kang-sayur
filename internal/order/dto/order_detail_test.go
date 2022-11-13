package dto

import (
	"testing"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestOrderDetailRequest_ToModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Dto      OrderDetailRequest
		Expected *model.OrderDetail
	}{
		{
			Name: "all filled",
			Dto: OrderDetailRequest{
				ItemID: 1,
				Qty:    1,
				Price:  10,
				Total:  10,
			},
			Expected: &model.OrderDetail{
				ItemID: 1,
				Qty:    1,
				Price:  10,
				Total:  10,
			},
		},
		{
			Name: "some filled",
			Dto: OrderDetailRequest{
				ItemID: 1,
			},
			Expected: &model.OrderDetail{
				ItemID: 1,
			},
		},
		{
			Name:     "dto empty",
			Dto:      OrderDetailRequest{},
			Expected: &model.OrderDetail{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			result := v.Dto.ToModel()
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestOrderDetailsRequest_ToModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Dto      OrderDetailsRequest
		Expected []model.OrderDetail
	}{
		{
			Name: "all filled",
			Dto: OrderDetailsRequest{{
				ItemID: 1,
				Qty:    1,
				Price:  10,
				Total:  10,
			}},
			Expected: []model.OrderDetail{{
				ItemID: 1,
				Qty:    1,
				Price:  10,
				Total:  10,
			}},
		},
		{
			Name: "some filled",
			Dto: OrderDetailsRequest{{
				ItemID: 1,
			}},
			Expected: []model.OrderDetail{{
				ItemID: 1,
			}},
		},
		{
			Name:     "dto empty",
			Dto:      OrderDetailsRequest{},
			Expected: []model.OrderDetail(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			result := v.Dto.ToModel()
			assert.Equal(t, &v.Expected, result)
		})
	}
}

func TestOrderDetailResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    *model.OrderDetail
		Expected OrderDetailResponse
	}{
		{
			Name: "all filled",
			Model: &model.OrderDetail{
				ID:    1,
				Qty:   1,
				Price: 10,
				Total: 10,
				Item: model.Item{
					Name: "item",
				},
			},
			Expected: OrderDetailResponse{
				ID:       1,
				ItemName: "item",
				Qty:      1,
				Price:    10,
				Total:    10,
			},
		},
		{
			Name: "some filled",
			Model: &model.OrderDetail{
				ID: 1,
			},
			Expected: OrderDetailResponse{
				ID: 1,
			},
		},
		{
			Name:     "dto empty",
			Model:    &model.OrderDetail{},
			Expected: OrderDetailResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result OrderDetailResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestOrderDetailsResponse_FromModel(t *testing.T) {
	testCase := []struct {
		Name     string
		Model    []model.OrderDetail
		Expected OrderDetailsResponse
	}{
		{
			Name: "all filled",
			Model: []model.OrderDetail{{
				ID:    1,
				Qty:   1,
				Price: 10,
				Total: 10,
				Item: model.Item{
					Name: "item",
				},
			}},
			Expected: OrderDetailsResponse{{
				ID:       1,
				ItemName: "item",
				Qty:      1,
				Price:    10,
				Total:    10,
			}},
		},
		{
			Name: "some filled",
			Model: []model.OrderDetail{{
				ID: 1,
			}},
			Expected: OrderDetailsResponse{{
				ID: 1,
			}},
		},
		{
			Name:     "dto empty",
			Model:    []model.OrderDetail{},
			Expected: OrderDetailsResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result OrderDetailsResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
