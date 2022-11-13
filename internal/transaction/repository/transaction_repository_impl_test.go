package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type suiteTransactionRepository struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *transactionRepositoryImpl
}

func (s *suiteTransactionRepository) SetupSuite() {
	db, mocking, _ := sqlmock.New()

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	s.mock = mocking
	s.repository = &transactionRepositoryImpl{
		db: dbGorm,
	}
}

func (s *suiteTransactionRepository) TestCreateTransaction() {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name        string
		Body        model.Transaction
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name: "success",
			Body: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "invalid foreign key (order id)",
			Body: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr: customerrors.ErrBadRequestBody,
			MockReturn:  errors.New("Cannot add or update a child row"),
		},
		{
			Name: "error",
			Body: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr: errors.New("err"),
			MockReturn:  errors.New("err"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`id`,`created_at`,`updated_at`,`deleted_at`,`order_id`,`transaction_status`,`transaction_time`,`signature_key`,`payment_type`,`gross_amount`,`settlement_time`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}
			var ctx context.Context
			err := s.repository.CreateTransaction(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
		})
	}
}

func (s *suiteTransactionRepository) TestUpdateTransaction() {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name        string
		Body        model.Transaction
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name: "success",
			Body: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "invalid foreign key (order id)",
			Body: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr: customerrors.ErrBadRequestBody,
			MockReturn:  errors.New("Cannot add or update a child row"),
		},
		{
			Name: "error",
			Body: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			ExpectedErr: errors.New("err"),
			MockReturn:  errors.New("err"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `transactions` SET `updated_at`=?,`transaction_status`=? WHERE id = ? AND `transactions`.`deleted_at` IS NULL"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}
			var ctx context.Context
			err := s.repository.UpdateTransaction(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
		})
	}
}
func (s *suiteTransactionRepository) TestFindAllTransaction() {
	userId := uuid.New().String()
	transactionId := uuid.New()

	testCase := []struct {
		Name        string
		UserId      string
		ExpectedErr error
		ExpectedRes []model.Transaction
		MockErr     error
		MockRes     *sqlmock.Rows
	}{
		{
			Name:        "success",
			UserId:      userId,
			ExpectedErr: nil,
			ExpectedRes: []model.Transaction{
				{
					ID: transactionId,
				},
			},
			MockErr: nil,
			MockRes: sqlmock.NewRows([]string{"id"}).AddRow(transactionId),
		},
		{
			Name:        "error",
			UserId:      userId,
			ExpectedErr: errors.New("err"),
			ExpectedRes: nil,
			MockErr:     errors.New("err"),
			MockRes:     nil,
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `transactions`.`id`,`transactions`.`created_at`,`transactions`.`updated_at`,`transactions`.`deleted_at`,`transactions`.`order_id`,`transactions`.`transaction_status`,`transactions`.`transaction_time`,`transactions`.`signature_key`,`transactions`.`payment_type`,`transactions`.`gross_amount`,`transactions`.`settlement_time` FROM `transactions` left join orders on orders.id = transactions.order_id WHERE orders.user_id = ? AND `transactions`.`deleted_at` IS NULL"))
			if v.MockRes != nil {
				db.WillReturnRows(v.MockRes)
			}
			db.WillReturnError(v.MockErr)
			var ctx context.Context
			res, err := s.repository.FindAllTransaction(v.UserId, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)
		})
	}
}
func (s *suiteTransactionRepository) TestFindTransaction() {
	transactionId := uuid.New()
	orderId := uuid.New()

	testCase := []struct {
		Name        string
		Transaction model.Transaction
		ExpectedErr error
		ExpectedRes model.Transaction
		MockErr     error
		MockRes     *sqlmock.Rows
	}{
		{
			Name: "success",
			Transaction: model.Transaction{
				ID: transactionId,
			},
			ExpectedErr: nil,
			ExpectedRes: model.Transaction{
				ID:                transactionId,
				OrderID:           orderId,
				TransactionStatus: "test",
				TransactionTime:   "test",
				SignatureKey:      "test",
				PaymentType:       "test",
				GrossAmount:       "test",
				SettlementTime:    "test",
			},
			MockErr: nil,
			MockRes: sqlmock.NewRows([]string{"transaction_id", "order_id", "transaction_status", "transaction_time", "signature_key", "payment_type", "gross_amount", "settlement_time"}).
				AddRow(transactionId, orderId, "test", "test", "test", "test", "test", "test"),
		},
		{
			Name: "error not found",
			Transaction: model.Transaction{
				ID: transactionId,
			},
			ExpectedErr: customerrors.ErrNotFound,
			ExpectedRes: model.Transaction{
				ID: transactionId,
			},
			MockErr: gorm.ErrRecordNotFound,
			MockRes: nil,
		},
		{
			Name: "other error",
			Transaction: model.Transaction{
				ID: transactionId,
			},
			ExpectedErr: errors.New("err"),
			ExpectedRes: model.Transaction{
				ID: transactionId,
			},
			MockErr: errors.New("err"),
			MockRes: nil,
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE `transactions`.`deleted_at` IS NULL AND `transactions`.`id` = ? ORDER BY `transactions`.`id` LIMIT 1"))
			if v.MockRes != nil {
				db.WillReturnRows(v.MockRes)
			}
			db.WillReturnError(v.MockErr)
			var ctx context.Context
			err := s.repository.FindTransaction(&v.Transaction, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, v.Transaction)
		})
	}
}

func TestSuiteTransactionRepository(t *testing.T) {
	suite.Run(t, new(suiteTransactionRepository))
}
