package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

// FindOrderById implements OrderRepository
func (r *orderRepositoryImpl) FindOrderById(order *model.Order, ctx context.Context) error {
	err := r.db.WithContext(ctx).First(&order).Error
	if err != nil {
		return customerrors.ErrNotFound
	}
	return nil
}

// CreateOrder implements OrderRepository
func (r *orderRepositoryImpl) CreateOrder(order *model.Order, ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return customerrors.ErrDuplicateData
		}
		if strings.Contains(err.Error(), "Cannot add or update a child row") {
			return customerrors.ErrBadRequestBody
		}
		return err
	}
	return nil
}

// FindOrder implements OrderRepository
func (r *orderRepositoryImpl) FindOrder(userId uuid.UUID, ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Preload("StatusOrder").Preload("User").Preload("Checkpoint").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

// FindOrderDetail implements OrderRepository
func (r *orderRepositoryImpl) FindOrderDetail(order *model.Order, ctx context.Context) error {
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", order.UserID, order.ID).Preload("OrderDetail.Item").Preload("OrderDetail").Preload("StatusOrder").Preload("Checkpoint").Find(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return customerrors.ErrNotFound
		}
		return err
	}
	return nil
}

// FindAllOrders implements OrderRepository
func (r *orderRepositoryImpl) FindAllOrders(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).Preload("OrderDetail").Preload("StatusOrder").Preload("User").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

// CencelOrder implements OrderRepository
func (r *orderRepositoryImpl) CencelOrder(orderId uuid.UUID, ctx context.Context) error {
	order := model.Order{
		ID: orderId,
	}
	res := r.db.WithContext(ctx).Model(&order).Update("status_order_id", constants.Cencel_status_order_id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// OrderReady implements OrderRepository
func (r *orderRepositoryImpl) OrderReady(orderId uuid.UUID, ctx context.Context) error {
	order := model.Order{
		ID: orderId,
	}
	res := r.db.WithContext(ctx).Model(&order).Update("status_order_id", constants.Ready_status_order_id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// OrderDone implements OrderRepository
func (r *orderRepositoryImpl) OrderDone(orderId uuid.UUID, ctx context.Context) error {
	order := model.Order{
		ID: orderId,
	}

	res := r.db.WithContext(ctx).Model(&order).Updates(&model.Order{
		StatusOrderID: constants.Success_status_order_id,
		Code:          "0",
		ExpiredOrder:  time.Now(),
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// OrderWaiting implements OrderRepository
func (r *orderRepositoryImpl) OrderWaiting(orderId uuid.UUID, ctx context.Context) error {
	order := model.Order{
		ID: orderId,
	}
	res := r.db.WithContext(ctx).Model(&order).Update("status_order_id", constants.Waiting_status_order_id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// InitStatusOrder implements OrderRepository
func (r *orderRepositoryImpl) InitStatusOrder() error {
	var status []model.StatusOrder
	r.db.Find(&status)
	if len(status) > 0 {
		return nil
	}
	err := r.db.Create(constants.StatusOrder).Error
	return err
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	orderRepo := &orderRepositoryImpl{
		db: db,
	}
	if err := orderRepo.InitStatusOrder(); err != nil {
		panic(err)
	}
	return orderRepo
}
