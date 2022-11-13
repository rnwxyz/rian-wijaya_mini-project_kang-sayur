package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type suiteOrderRepository struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *orderRepositoryImpl
}

func (s *suiteOrderRepository) SetupSuite() {
	db, mocking, _ := sqlmock.New()

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	s.mock = mocking
	s.repository = &orderRepositoryImpl{
		db: dbGorm,
	}
}

func (s *suiteOrderRepository) TestCreateOrder() {
	orderId := uuid.New()

	testCase := []struct {
		Name           string
		Order          model.Order
		ExpectedErr    error
		CreateOrderErr error
	}{
		{
			Name: "create order success",
			Order: model.Order{
				ID: orderId,
			},
			ExpectedErr:    nil,
			CreateOrderErr: nil,
		},
		{
			Name: "create error duplicate entry",
			Order: model.Order{
				ID: orderId,
			},
			ExpectedErr:    customerrors.ErrDuplicateData,
			CreateOrderErr: errors.New("Duplicate entry"),
		},
		{
			Name: "create error foreign key",
			Order: model.Order{
				ID: orderId,
			},
			ExpectedErr:    customerrors.ErrBadRequestBody,
			CreateOrderErr: errors.New("Cannot add or update a child row"),
		},
		{
			Name: "other error when create order",
			Order: model.Order{
				ID: orderId,
			},
			ExpectedErr:    errors.New("other error"),
			CreateOrderErr: errors.New("other error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `orders` (`id`,`created_at`,`updated_at`,`deleted_at`,`user_id`,`checkpoint_id`,`status_order_id`,`shipping_cost`,`total_price`,`grand_total`,`code`,`hash`,`expired_order`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.CreateOrderErr)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}

			var ctx context.Context
			err := s.repository.CreateOrder(&v.Order, ctx)

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}

func (s *suiteOrderRepository) TestFindOrderById() {
	orderId := uuid.New()

	testCase := []struct {
		Name         string
		Order        model.Order
		ExpectedErr  error
		FistOrderErr error
		FistOrderRes *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			Order: model.Order{
				ID: orderId,
			},
			FistOrderErr: nil,
			FistOrderRes: sqlmock.NewRows([]string{"id", "user_id", "checkpoint_id"}).
				AddRow(orderId, uuid.New(), uuid.New()),
		},
		{
			Name:        "error",
			ExpectedErr: customerrors.ErrNotFound,
			Order: model.Order{
				ID: orderId,
			},
			FistOrderErr: gorm.ErrRecordNotFound,
			FistOrderRes: sqlmock.NewRows([]string{"id", "user_id", "checkpoint_id"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE `orders`.`deleted_at` IS NULL AND `orders`.`id` = ? ORDER BY `orders`.`id` LIMIT 1"))
			if v.FistOrderErr != nil {
				db.WillReturnError(v.FistOrderErr)
			} else {
				db.WillReturnRows(v.FistOrderRes)
			}
			var ctx context.Context
			err := s.repository.FindOrderById(&v.Order, ctx)

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}

func (s *suiteOrderRepository) TestFindOrder() {
	userId := uuid.New()
	orderId1 := uuid.New()
	// orderId2 := uuid.New()
	checkpointId1 := uuid.New()
	// checkpointId2 := uuid.New()
	query := regexp.QuoteMeta("SELECT * FROM `orders` WHERE user_id = ? AND `orders`.`deleted_at` IS NULL")
	preloadCheckpoint := regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE `checkpoints`.`id` = ? AND `checkpoints`.`deleted_at` IS NULL")
	preloadStatusOrder := regexp.QuoteMeta("SELECT * FROM `status_orders` WHERE `status_orders`.`id` = ? AND `status_orders`.`deleted_at` IS NULL")
	preloadUser := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")

	testCase := []struct {
		Name                  string
		UserId                uuid.UUID
		ExpectedErr           error
		ExpectedRes           []model.Order
		FindOrderErr          error
		FindOrderRes          *sqlmock.Rows
		PreloadUserRes        *sqlmock.Rows
		PreloadStatusOrderRes *sqlmock.Rows
		PreloadCheckpointRes  *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: []model.Order{
				{
					ID:     orderId1,
					UserID: userId,
					User: model.User{
						ID:   userId,
						Name: "user",
					},
					CheckpointID: checkpointId1,
					Checkpoint: model.Checkpoint{
						ID:   checkpointId1,
						Name: "checkpoint 1",
					},
					StatusOrderID: constants.Pending_status_order_id,
					StatusOrder: model.StatusOrder{
						ID:   constants.Pending_status_order_id,
						Name: "pending",
					},
				},
			},
			UserId:       userId,
			FindOrderErr: nil,
			FindOrderRes: sqlmock.NewRows([]string{"id", "user_id", "checkpoint_id", "status_order_id"}).
				AddRow(orderId1, userId, checkpointId1, 1),
			PreloadUserRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(userId, "user"),
			PreloadCheckpointRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(checkpointId1, "checkpoint 1"),
			PreloadStatusOrderRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(constants.Pending_status_order_id, "pending"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(query).WillReturnError(v.FindOrderErr).WillReturnRows(v.FindOrderRes)

			s.mock.ExpectQuery(preloadCheckpoint).WillReturnRows(v.PreloadCheckpointRes)

			s.mock.ExpectQuery(preloadStatusOrder).WillReturnRows(v.PreloadStatusOrderRes)

			s.mock.ExpectQuery(preloadUser).WillReturnRows(v.PreloadUserRes)

			res, err := s.repository.FindOrder(v.UserId, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func (s *suiteOrderRepository) TestFindOrderDetail() {
	userId := uuid.New()
	orderId := uuid.New()
	checkpointId := uuid.New()

	testCase := []struct {
		Name                  string
		Order                 model.Order
		ExpectedErr           error
		ExpectedRes           model.Order
		FindOrderDetailErr    error
		FindOrderDetailRes    *sqlmock.Rows
		PreloadOrderDetailRes *sqlmock.Rows
		PreloadStatusOrderRes *sqlmock.Rows
		PreloadCheckpointRes  *sqlmock.Rows
		PreloadItemRes        *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: model.Order{
				ID:           orderId,
				UserID:       userId,
				CheckpointID: checkpointId,
				Checkpoint: model.Checkpoint{
					ID:   checkpointId,
					Name: "checkpoint 1",
				},
				StatusOrderID: constants.Pending_status_order_id,
				StatusOrder: model.StatusOrder{
					ID:   constants.Pending_status_order_id,
					Name: "pending",
				},
				OrderDetail: []model.OrderDetail{
					{
						OrderID: orderId,
						ItemID:  1,
						Item: model.Item{
							ID:   1,
							Name: "item",
						},
						Qty:   1,
						Price: 1,
					},
				},
			},
			Order: model.Order{
				ID:     orderId,
				UserID: userId,
			},
			FindOrderDetailErr: nil,
			FindOrderDetailRes: sqlmock.NewRows([]string{"id", "user_id", "checkpoint_id", "status_order_id"}).
				AddRow(orderId, userId, checkpointId, constants.Pending_status_order_id),
			PreloadOrderDetailRes: sqlmock.NewRows([]string{"order_id", "item_id", "qty", "price"}).
				AddRow(orderId, 1, 1, 1),
			PreloadCheckpointRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(checkpointId, "checkpoint 1"),
			PreloadStatusOrderRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(constants.Pending_status_order_id, "pending"),
			PreloadItemRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "item"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			findOrderMock := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE (user_id = ? AND id = ?) AND `orders`.`deleted_at` IS NULL AND `orders`.`id` = ?"))
			findOrderMock.WillReturnError(v.FindOrderDetailErr)
			findOrderMock.WillReturnRows(v.FindOrderDetailRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `checkpoints` WHERE `checkpoints`.`id` = ? AND `checkpoints`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadCheckpointRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `order_details` WHERE `order_type` = ? AND `order_details`.`order_id` = ? AND `order_details`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadOrderDetailRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `items` WHERE `items`.`id` = ? AND `items`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadItemRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status_orders` WHERE `status_orders`.`id` = ? AND `status_orders`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadStatusOrderRes)

			// var ctx context.Context
			err := s.repository.FindOrderDetail(&v.Order, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, v.Order)

			s.TearDown()
		})
	}
}
func (s *suiteOrderRepository) TestFindAllOrders() {
	userId := uuid.New()
	orderId := uuid.New()
	checkpointId := uuid.New()

	testCase := []struct {
		Name                  string
		Order                 model.Order
		ExpectedErr           error
		ExpectedRes           []model.Order
		FindOrderAllErr       error
		FindOrderAllRes       *sqlmock.Rows
		PreloadOrderDetailRes *sqlmock.Rows
		PreloadStatusOrderRes *sqlmock.Rows
		PreloadUserRes        *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: []model.Order{{
				ID:     orderId,
				UserID: userId,
				User: model.User{
					ID:   userId,
					Name: "user",
				},
				CheckpointID:  checkpointId,
				StatusOrderID: constants.Pending_status_order_id,
				StatusOrder: model.StatusOrder{
					ID:   constants.Pending_status_order_id,
					Name: "pending",
				},
				OrderDetail: []model.OrderDetail{
					{
						OrderID: orderId,
						ItemID:  1,
						Qty:     1,
						Price:   1,
					},
				},
			}},
			Order: model.Order{
				ID:     orderId,
				UserID: userId,
			},
			FindOrderAllErr: nil,
			FindOrderAllRes: sqlmock.NewRows([]string{"id", "user_id", "checkpoint_id", "status_order_id"}).
				AddRow(orderId, userId, checkpointId, constants.Pending_status_order_id),
			PreloadOrderDetailRes: sqlmock.NewRows([]string{"order_id", "item_id", "qty", "price"}).
				AddRow(orderId, 1, 1, 1),
			PreloadStatusOrderRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(constants.Pending_status_order_id, "pending"),
			PreloadUserRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(userId, "user"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			findOrderMock := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE `orders`.`deleted_at` IS NULL"))
			findOrderMock.WillReturnError(v.FindOrderAllErr)
			findOrderMock.WillReturnRows(v.FindOrderAllRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `order_details` WHERE `order_type` = ? AND `order_details`.`order_id` = ? AND `order_details`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadOrderDetailRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status_orders` WHERE `status_orders`.`id` = ? AND `status_orders`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadStatusOrderRes)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadUserRes)

			var ctx context.Context
			res, err := s.repository.FindAllOrders(ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func (s *suiteOrderRepository) TearDown() {
	s.mock = nil
	s.repository = nil
}

func TestSuiteOrderRepository(t *testing.T) {
	suite.Run(t, new(suiteOrderRepository))
}
