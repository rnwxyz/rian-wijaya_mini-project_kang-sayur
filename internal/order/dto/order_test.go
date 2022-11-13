package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestOrderResponse_FromModel(t *testing.T) {
	orderId := uuid.New()
	checkpointId := uuid.New()

	testCase := []struct {
		Name     string
		Model    *model.Order
		Expected OrderResponse
	}{
		{
			Name: "all filled",
			Model: &model.Order{
				ID:           orderId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint",
				},
				StatusOrderID: 1,
				StatusOrder: model.StatusOrder{
					ID:   1,
					Name: "pending",
				},
				ShippingCost: 5000,
				TotalPrice:   20000,
				GrandTotal:   25000,
				ExpiredOrder: time.Time{},
			},
			Expected: OrderResponse{
				ID:              orderId,
				CheckpointID:    checkpointId,
				CheckpointName:  "checkpoint",
				StatusOrderName: "pending",
				ShippingCost:    5000,
				TotalPrice:      20000,
				GrandTotal:      25000,
				ExpiredOrder:    time.Time{},
			},
		},
		{
			Name: "some filled",
			Model: &model.Order{
				ID:           orderId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint",
				},
				StatusOrderID: 1,
				StatusOrder: model.StatusOrder{
					ID:   1,
					Name: "pending",
				},
			},
			Expected: OrderResponse{
				ID:              orderId,
				CheckpointID:    checkpointId,
				CheckpointName:  "checkpoint",
				StatusOrderName: "pending",
			},
		},
		{
			Name:     "dto empty",
			Model:    &model.Order{},
			Expected: OrderResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result OrderResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestOrderWithDetailResponse_FromModel(t *testing.T) {
	orderId := uuid.New()
	checkpointId := uuid.New()

	testCase := []struct {
		Name     string
		Model    *model.Order
		Expected OrderWithDetailResponse
	}{
		{
			Name: "all filled",
			Model: &model.Order{
				ID:           orderId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint",
				},
				StatusOrderID: 1,
				StatusOrder: model.StatusOrder{
					ID:   1,
					Name: "pending",
				},
				ShippingCost: 5000,
				TotalPrice:   20000,
				GrandTotal:   25000,
				ExpiredOrder: time.Time{},
				Hash:         "qwerty",
				OrderDetail: []model.OrderDetail{
					{
						ID: uint(1),
						Item: model.Item{
							ID:   uint(1),
							Name: "item",
						},
						Qty:   1,
						Price: 10,
					},
				},
			},
			Expected: OrderWithDetailResponse{
				ID:              orderId,
				CheckpointID:    checkpointId,
				CheckpointName:  "checkpoint",
				StatusOrderName: "pending",
				ShippingCost:    5000,
				TotalPrice:      20000,
				GrandTotal:      25000,
				ExpiredOrder:    time.Time{},
				Hash:            "qwerty",
				OrderDetail: OrderDetailsResponse{
					{
						ID:       1,
						ItemName: "item",
						Qty:      1,
						Price:    10,
					},
				},
			},
		},
		{
			Name: "some filled",
			Model: &model.Order{
				ID:           orderId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint",
				},
				StatusOrderID: 1,
				StatusOrder: model.StatusOrder{
					ID:   1,
					Name: "pending",
				},
			},
			Expected: OrderWithDetailResponse{
				ID:              orderId,
				CheckpointID:    checkpointId,
				CheckpointName:  "checkpoint",
				StatusOrderName: "pending",
			},
		},
		{
			Name:     "dto empty",
			Model:    &model.Order{},
			Expected: OrderWithDetailResponse{},
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result OrderWithDetailResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
func TestOrdersResponse_FromModel(t *testing.T) {
	orderId := uuid.New()
	checkpointId := uuid.New()

	testCase := []struct {
		Name     string
		Model    []model.Order
		Expected OrdersResponse
	}{
		{
			Name: "all filled",
			Model: []model.Order{{
				ID:           orderId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint",
				},
				StatusOrderID: 1,
				StatusOrder: model.StatusOrder{
					ID:   1,
					Name: "pending",
				},
				ShippingCost: 5000,
				TotalPrice:   20000,
				GrandTotal:   25000,
				ExpiredOrder: time.Time{},
			}},
			Expected: OrdersResponse{{
				ID:              orderId,
				CheckpointID:    checkpointId,
				CheckpointName:  "checkpoint",
				StatusOrderName: "pending",
				ShippingCost:    5000,
				TotalPrice:      20000,
				GrandTotal:      25000,
				ExpiredOrder:    time.Time{},
			}},
		},
		{
			Name: "some filled",
			Model: []model.Order{{
				ID:           orderId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint",
				},
				StatusOrderID: 1,
				StatusOrder: model.StatusOrder{
					ID:   1,
					Name: "pending",
				},
			}},
			Expected: OrdersResponse{{
				ID:              orderId,
				CheckpointID:    checkpointId,
				CheckpointName:  "checkpoint",
				StatusOrderName: "pending",
			}},
		},
		{
			Name:     "dto empty",
			Model:    []model.Order{},
			Expected: OrdersResponse(nil),
		},
	}
	for _, v := range testCase {
		t.Run(v.Name, func(t *testing.T) {
			var result OrdersResponse
			result.FromModel(v.Model)
			assert.Equal(t, v.Expected, result)
		})
	}
}
