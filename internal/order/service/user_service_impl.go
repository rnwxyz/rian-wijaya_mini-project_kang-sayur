package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	it "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type orderServiceImpl struct {
	orderRepo or.OrderRepository
	itemRepo  it.ItemRepository
}

// CencelOder implements OrderService
func (s *orderServiceImpl) CencelOder(orderId string, ctx context.Context) error {
	id, err := uuid.Parse(orderId)
	if err != nil {
		return utils.ErrInvalidId
	}
	err = s.orderRepo.CencelOrder(id, ctx)
	return err
}

// CreateOrder implements OrderService
func (s *orderServiceImpl) CreateOrder(body dto.OrderDetailsRequest, userId string, ctx context.Context) (uuid.UUID, error) {
	newId := uuid.New()
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, utils.ErrInvalidId
	}
	totalPrice := 0

	// validating item and sum price
	for _, ord := range body {
		var item model.Item
		item.ID = ord.ItemID
		err := s.itemRepo.FindItemById(&item, ctx)
		if err != nil {
			return uuid.Nil, err
		}
		ord.Total = (ord.Qty * item.Price)
		totalPrice += ord.Total
	}

	orderDetail := body.ToModel()
	loc, err := time.LoadLocation(config.Cfg.TIME_LOCATION)
	if err != nil {
		return uuid.Nil, err
	}

	newOrder := model.Order{
		ID:            newId,
		UserID:        userIdUUID,
		ShippingCost:  constants.Shipping_cost,
		StatusOrderID: constants.Pending_status_order_id,
		TotalPrice:    totalPrice,
		GrandTotal:    totalPrice + constants.Shipping_cost,
		OrderDetail:   *orderDetail,
		ExpiredOrder:  time.Now().In(loc).Add(constants.ExpOrder),
	}
	err = s.orderRepo.CreateOrder(&newOrder, ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return newId, nil
}

// FindOrder implements OrderService
func (s *orderServiceImpl) FindOrder(userId string, ctx context.Context) (dto.OrdersResponse, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	orders, err := s.orderRepo.FindOrder(userIdUUID, ctx)
	if err != nil {
		return nil, err
	}
	var ordersResponse dto.OrdersResponse
	ordersResponse.FromModel(orders)
	return ordersResponse, nil
}

// FindOrderDetail implements OrderService
func (s *orderServiceImpl) FindOrderDetail(userId string, orderId string, ctx context.Context) (*dto.OrderWithDetailResponse, error) {
	orderIdUUID, err := uuid.Parse(orderId)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	order := model.Order{
		ID:     orderIdUUID,
		UserID: userIdUUID,
	}
	err = s.orderRepo.FindOrderDetail(&order, ctx)
	if err != nil {
		return nil, err
	}
	var orderDetailReponse dto.OrderWithDetailResponse
	orderDetailReponse.FromModel(&order)
	return &orderDetailReponse, nil
}

func NewOrderService(orRepository or.OrderRepository, itRepository it.ItemRepository) OrderService {
	return &orderServiceImpl{
		orderRepo: orRepository,
		itemRepo:  itRepository,
	}
}
