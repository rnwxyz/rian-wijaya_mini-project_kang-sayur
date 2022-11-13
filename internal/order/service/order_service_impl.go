package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	it "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	urp "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
)

type midtrans interface {
	NewTransaction(order model.Order, user model.User) (string, error)
}

type orderServiceImpl struct {
	orderRepo or.OrderRepository
	itemRepo  it.ItemRepository
	payment   midtrans
	userRepo  urp.UserRepository
}

func NewOrderService(orRepository or.OrderRepository, itRepository it.ItemRepository, midtrans midtrans, userRepo urp.UserRepository) OrderService {
	return &orderServiceImpl{
		orderRepo: orRepository,
		itemRepo:  itRepository,
		payment:   midtrans,
		userRepo:  userRepo,
	}
}

// CreateOrder implements OrderService
func (s *orderServiceImpl) CreateOrder(body dto.OrderRequest, userId string, ctx context.Context) (*dto.NewOrder, error) {
	newId := uuid.New()
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, customerrors.ErrInvalidId
	}
	checkpointIdUUID, err := uuid.Parse(body.CheckpointID)
	if err != nil {
		return nil, customerrors.ErrInvalidId
	}
	totalPrice := 0

	// validating item and sum price
	for i, ord := range body.Order {
		var item model.Item
		item.ID = ord.ItemID
		err := s.itemRepo.FindItemById(&item, ctx)
		if err != nil {
			return nil, customerrors.ErrBadRequestBody
		}
		if item.Qty < ord.Qty || ord.Qty < 1 {
			return nil, customerrors.ErrQtyOrder
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

	user, _ := s.userRepo.FindUserByID(userId, ctx)

	transaction, err := s.payment.NewTransaction(newOrder, *user)
	if err != nil {
		return nil, err
	}
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
		return nil, customerrors.ErrInvalidId
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
		return nil, customerrors.ErrInvalidId
	}
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, customerrors.ErrInvalidId
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
func (s *orderServiceImpl) TakeOrder(body dto.TakeOrder, ctx context.Context) error {
	byt, err := base64.StdEncoding.DecodeString(body.Code)
	if err != nil {
		return customerrors.ErrOrderCode
	}
	// data have 3 index orderId, order code, and checkpoint id
	data := strings.Split(string(byt), " ")

	id, err := uuid.Parse(data[0])
	if err != nil {
		return customerrors.ErrOrderCode
	}

	orderCode := data[1]
	checkpointId, err := uuid.Parse(data[2])
	if err != nil {
		return customerrors.ErrOrderCode
	}

	order := model.Order{
		ID: id,
	}
	err = s.orderRepo.FindOrderById(&order, ctx)
	if err != nil {
		return err
	}

	if body.CheckpointID != checkpointId.String() {
		return customerrors.ErrWrongCheckpoint
	}

	fmt.Println(checkpointId, "  ", order.CheckpointID)

	if orderCode != order.Code {
		return customerrors.ErrCodeUsed
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
		return customerrors.ErrInvalidId
	}
	order := model.Order{
		ID: id,
	}
	err = s.orderRepo.FindOrderById(&order, ctx)
	if err != nil {
		return err
	}
	if order.StatusOrderID != constants.Waiting_status_order_id {
		return customerrors.ErrUpdateStatusOrder
	}
	err = s.orderRepo.OrderReady(id, ctx)
	return err
}

// CencelOder implements OrderService
func (s *orderServiceImpl) CencelOder(orderId string, ctx context.Context) error {
	id, err := uuid.Parse(orderId)
	if err != nil {
		return customerrors.ErrInvalidId
	}
	order := model.Order{
		ID: id,
	}
	err = s.orderRepo.FindOrderById(&order, ctx)
	if err != nil {
		return err
	}
	if order.StatusOrderID == constants.Pending_status_order_id {
		return customerrors.ErrUpdateStatusOrder
	}
	err = s.orderRepo.CencelOrder(id, ctx)
	return err
}

// SetOrderStatus implements OrderService
func (s *orderServiceImpl) SetOrderStatus(orderId uuid.UUID, status string, ctx context.Context) error {
	order := model.Order{
		ID: orderId,
	}
	err := s.orderRepo.FindOrderDetail(&order, ctx)
	if err != nil {
		return err
	}
	switch status {
	case "capture":
		err := s.orderRepo.OrderWaiting(orderId, ctx)
		if err != nil {
			return err
		}
	case "settlement":
		err := s.orderRepo.OrderWaiting(orderId, ctx)
		if err != nil {
			return err
		}
	case "deny":
		err := s.orderRepo.CencelOrder(orderId, ctx)
		if err != nil {
			return err
		}
	case "cencel":
		err := s.orderRepo.CencelOrder(orderId, ctx)
		if err != nil {
			return err
		}
	case "expired":
		err := s.orderRepo.CencelOrder(orderId, ctx)
		if err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}
