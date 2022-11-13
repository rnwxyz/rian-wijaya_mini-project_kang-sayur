package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	it "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository"
	itemRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/item/repository/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/dto"
	or "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository"
	orderRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/order/repository/mock"
	urp "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	userRepositoryMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository/mock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	midtransMock "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/payment/mock"
	"github.com/stretchr/testify/suite"
)

type suiteOrderService struct {
	suite.Suite
	orderRepositoryMock *orderRepositoryMock.OrderRepositoryMock
	itemRepositoryMock  *itemRepositoryMock.ItemRepositoryMock
	userRepositoryMock  *userRepositoryMock.UserRepositoryMock
	payment             *midtransMock.MidtransMock
	orderService        OrderService
}

func newOrderService(orRepository or.OrderRepository, itRepository it.ItemRepository, midtrans midtrans, userRepo urp.UserRepository) OrderService {
	return &orderServiceImpl{
		orderRepo: orRepository,
		itemRepo:  itRepository,
		payment:   midtrans,
		userRepo:  userRepo,
	}
}

func (s *suiteOrderService) SetupSuit() {
	s.orderRepositoryMock = new(orderRepositoryMock.OrderRepositoryMock)
	s.itemRepositoryMock = new(itemRepositoryMock.ItemRepositoryMock)
	s.userRepositoryMock = new(userRepositoryMock.UserRepositoryMock)
	s.payment = new(midtransMock.MidtransMock)
	s.orderService = newOrderService(s.orderRepositoryMock, s.itemRepositoryMock, s.payment, s.userRepositoryMock)
}

func (s *suiteOrderService) TearDown() {
	s.orderRepositoryMock = nil
	s.itemRepositoryMock = nil
	s.userRepositoryMock = nil
	s.payment = nil
	s.orderService = nil
}

// cant implement mock pass by reference

// func (s *suiteOrderService) TestCreateOrder() {
// 	userId := uuid.New()
// 	orderId := uuid.New()
// 	checkpointId := uuid.New().String()

// 	testCase := []struct {
// 		Name              string
// 		ExpectedErr       error
// 		ExpectedRes       dto.NewOrder
// 		Body              dto.OrderRequest
// 		UserId            string
// 		FindItemByIdErr   error
// 		CreateOrderdErr   error
// 		FindUserByIdRes   model.User
// 		NewTransactionErr error
// 		NewTransactionRes string
// 	}{
// 		{
// 			Name:        "success",
// 			ExpectedErr: nil,
// 			ExpectedRes: dto.NewOrder{
// 				OrderID:     orderId,
// 				RedirectURL: "https://testing/midtrans.com",
// 			},
// 			Body: dto.OrderRequest{
// 				CheckpointID: checkpointId,
// 				Order: dto.OrderDetailsRequest{
// 					{
// 						ItemID: 1,
// 						Qty:    1,
// 					},
// 				},
// 			},
// 			UserId:          userId.String(),
// 			FindItemByIdErr: nil,
// 			CreateOrderdErr: nil,
// 			FindUserByIdRes: model.User{
// 				ID:   userId,
// 				Name: "user",
// 			},
// 			NewTransactionErr: nil,
// 			NewTransactionRes: "https://testing/midtrans.com",
// 		},
// 	}
// 	for _, v := range testCase {
// 		s.T().Run(v.Name, func(t *testing.T) {
// 			// s.itemRepositoryMock.On("FindItemById").Return(v.FindItemByIdErr)
// 			s.orderRepositoryMock.On("CreateOrder").Return(v.CreateOrderdErr)
// 			s.userRepositoryMock.On("FindUserID").Return(v.FindUserByIdRes, v.FindItemByIdErr)
// 			s.payment.On("NewTransaction").Return(v.NewTransactionRes, v.NewTransactionErr)

// 			var ctx context.Context
// 			_, err := s.orderService.CreateOrder(v.Body, v.UserId, ctx)

// 			s.Equal(v.ExpectedErr, err)

// 			s.TearDown()
// 		})
// 	}
// }

func (s *suiteOrderService) TestFindOrder() {
	orderId := uuid.New()
	userId := uuid.New()

	testCase := []struct {
		Name         string
		ExpectedErr  error
		ExpectedRes  dto.OrdersResponse
		UserId       string
		FindOrderErr error
		FindOrderRes []model.Order
	}{
		{
			Name:        "find order success",
			ExpectedErr: nil,
			ExpectedRes: dto.OrdersResponse{
				{
					ID:             orderId,
					UserID:         userId,
					UserName:       "user",
					CheckpointName: "checkpoint",
				},
			},
			UserId:       userId.String(),
			FindOrderErr: nil,
			FindOrderRes: []model.Order{
				{
					ID:     orderId,
					UserID: userId,
					User: model.User{
						ID:   userId,
						Name: "user",
					},
					Checkpoint: model.Checkpoint{
						Name: "checkpoint",
					},
				},
			},
		},
		{
			Name:         "find order invalid id",
			ExpectedErr:  customerrors.ErrInvalidId,
			ExpectedRes:  dto.OrdersResponse(nil),
			UserId:       "123",
			FindOrderErr: nil,
			FindOrderRes: []model.Order{},
		},
		{
			Name:         "find order internal error",
			ExpectedErr:  errors.New("internal error"),
			ExpectedRes:  dto.OrdersResponse(nil),
			UserId:       userId.String(),
			FindOrderErr: errors.New("internal error"),
			FindOrderRes: []model.Order{},
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.orderRepositoryMock.On("FindOrder").Return(v.FindOrderRes, v.FindOrderErr)

			var ctx context.Context
			res, err := s.orderService.FindOrder(v.UserId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteOrderService) TestFindOrderDetail() {
	orderId := uuid.New()
	userId := uuid.New()

	testCase := []struct {
		Name         string
		ExpectedErr  error
		ExpectedRes  *dto.OrderWithDetailResponse
		UserId       string
		OrderId      string
		FindOrderErr error
	}{
		{
			Name:        "find order detail success",
			ExpectedErr: nil,
			ExpectedRes: &dto.OrderWithDetailResponse{
				ID: orderId,
			},
			UserId:       userId.String(),
			OrderId:      orderId.String(),
			FindOrderErr: nil,
		},
		{
			Name:         "find order invalid user id",
			ExpectedErr:  customerrors.ErrInvalidId,
			ExpectedRes:  nil,
			UserId:       "123",
			OrderId:      orderId.String(),
			FindOrderErr: nil,
		},
		{
			Name:         "find order invalid order id",
			ExpectedErr:  customerrors.ErrInvalidId,
			ExpectedRes:  nil,
			UserId:       userId.String(),
			OrderId:      "123",
			FindOrderErr: nil,
		},
		{
			Name:         "find order detail internal error",
			ExpectedErr:  errors.New("internal error"),
			ExpectedRes:  nil,
			UserId:       userId.String(),
			OrderId:      orderId.String(),
			FindOrderErr: errors.New("internal error"),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.orderRepositoryMock.On("FindOrderDetail").Return(v.FindOrderErr)

			var ctx context.Context
			res, err := s.orderService.FindOrderDetail(v.UserId, v.OrderId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteOrderService) TestFindAllOrders() {
	orderId := uuid.New()
	userId := uuid.New()

	testCase := []struct {
		Name             string
		ExpectedErr      error
		ExpectedRes      dto.OrdersResponse
		FindAllOrdersErr error
		FindAllOrdersRes []model.Order
	}{
		{
			Name:        "find order detail success",
			ExpectedErr: nil,
			ExpectedRes: dto.OrdersResponse{
				{
					ID:             orderId,
					UserID:         userId,
					UserName:       "user",
					CheckpointName: "checkpoint",
				},
			},
			FindAllOrdersErr: nil,
			FindAllOrdersRes: []model.Order{
				{
					ID:     orderId,
					UserID: userId,
					User: model.User{
						ID:   userId,
						Name: "user",
					},
					Checkpoint: model.Checkpoint{
						Name: "checkpoint",
					},
				},
			},
		},
		{
			Name:             "error find all orders",
			ExpectedErr:      errors.New("internal error"),
			ExpectedRes:      dto.OrdersResponse(nil),
			FindAllOrdersErr: errors.New("internal error"),
			FindAllOrdersRes: []model.Order(nil),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuit()

			s.orderRepositoryMock.On("FindAllOrders").Return(v.FindAllOrdersRes, v.FindAllOrdersErr)

			var ctx context.Context
			res, err := s.orderService.FindAllOrders(ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestSuiteOrderService(t *testing.T) {
	suite.Run(t, new(suiteOrderService))
}
