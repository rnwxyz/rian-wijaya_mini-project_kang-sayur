package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

// CreateOrder implements OrderRepository
func (r *orderRepositoryImpl) CreateOrder(order *model.Order, ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrDuplicateData
		}
		if strings.Contains(err.Error(), "Cannot add or update a child row") {
			return utils.ErrBadRequestBody
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
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", order.UserID, order.ID).Preload("OrderDetail").Preload("StatusOrder").Preload("Checkpoint").Find(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
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
		return utils.ErrNotFound
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
		return utils.ErrNotFound
	}
	return nil
}

// OrderDone implements OrderRepository
func (r *orderRepositoryImpl) OrderDone(orderId uuid.UUID, ctx context.Context) error {
	order := model.Order{
		ID: orderId,
	}
	loc, err := time.LoadLocation(config.Cfg.TIME_LOCATION)
	if err != nil {
		return utils.ErrTimeLocation
	}
	// set zero
	zeroTime := time.Date(int(1), time.January, int(1), int(0), int(0), int(0), int(0), loc)
	res := r.db.WithContext(ctx).Model(&order).Updates(&model.Order{
		StatusOrderID: constants.Success_status_order_id,
		ExpiredOrder:  zeroTime,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return utils.ErrNotFound
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
		return utils.ErrNotFound
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
