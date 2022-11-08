package service

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
	it "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/payment"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type orderServiceImpl struct {
	orderRepo or.OrderRepository
	itemRepo  it.ItemRepository
	payment   payment.Midtrans
}

func NewOrderService(orRepository or.OrderRepository, itRepository it.ItemRepository, midtrans payment.Midtrans) OrderService {
	return &orderServiceImpl{
		orderRepo: orRepository,
		itemRepo:  itRepository,
		payment:   midtrans,
	}
}

// CreateOrder implements OrderService
func (s *orderServiceImpl) CreateOrder(body dto.OrderRequest, userId string, ctx context.Context) (*dto.NewOrder, error) {
	newId := uuid.New()
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	checkpointIdUUID, err := uuid.Parse(body.CheckpointID)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	totalPrice := 0

	// validating item and sum price
	for i, ord := range body.Order {
		var item model.Item
		item.ID = ord.ItemID
		err := s.itemRepo.FindItemById(&item, ctx)
		if err != nil {
			return nil, utils.ErrBadRequestBody
		}
		if item.Qty < ord.Qty || ord.Qty < 1 {
			return nil, utils.ErrQtyOrder
		}
		body.Order[i].Price = item.Price
		body.Order[i].Total = (ord.Qty * item.Price)
		totalPrice += body.Order[i].Total
	}

	orderDetail := body.Order.ToModel()

	newOrder := model.Order{
		ID:            newId,
		UserID:        userIdUUID,
		CheckpointID:  checkpointIdUUID,
		ShippingCost:  constants.Shipping_cost,
		StatusOrderID: constants.Pending_status_order_id,
		TotalPrice:    totalPrice,
		GrandTotal:    totalPrice + constants.Shipping_cost,
		OrderDetail:   *orderDetail,
		ExpiredOrder:  time.Now().Add(constants.ExpOrder),
	}
	err = s.orderRepo.CreateOrder(&newOrder, ctx)
	if err != nil {
		return nil, err
	}
	transaction := s.payment.NewTransaction(newOrder)
	newOrderResponse := dto.NewOrder{
		OrderID:     newId,
		RedirectURL: transaction,
	}
	return &newOrderResponse, nil
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

// TakeOrder implements OrderService
func (s *orderServiceImpl) TakeOrder(code string, ctx context.Context) error {
	byt, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return utils.ErrOrderCode
	}
	// data have 3 index orderId, order code, and checkpoint id
	data := strings.Split(string(byt), " ")
	id, err := uuid.Parse(data[0])
	if err != nil {
		return utils.ErrOrderCode
	}
	orderCode := data[1]
	checkpointId, err := uuid.Parse(data[2])
	if err != nil {
		return utils.ErrOrderCode
	}
	order := model.Order{
		ID:           id,
		Code:         orderCode,
		CheckpointID: checkpointId,
	}
	err = s.orderRepo.FindOrderDetail(&order, ctx)
	if err != nil {
		return utils.ErrOrderCode
	}
	err = s.orderRepo.OrderDone(id, ctx)
	if err != nil {
		return err
	}
	return nil
}

// FindAllOrders implements OrderService
func (s *orderServiceImpl) FindAllOrders(ctx context.Context) (dto.OrdersResponse, error) {
	orders, err := s.orderRepo.FindAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	var ordersResponse dto.OrdersResponse
	ordersResponse.FromModel(orders)
	return ordersResponse, nil
}

// OderReady implements OrderService
func (s *orderServiceImpl) OrderReady(orderId string, ctx context.Context) error {
	id, err := uuid.Parse(orderId)
	if err != nil {
		return utils.ErrInvalidId
	}
	err = s.orderRepo.OrderReady(id, ctx)
	return err
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
