package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type suiteUserRepository struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *userRepositoryImpl
}

func (s *suiteUserRepository) SetupSuite() {
	db, mocking, _ := sqlmock.New()

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	s.mock = mocking
	s.repository = &userRepositoryImpl{
		db: dbGorm,
	}
}

func (s *suiteUserRepository) TestCreateUser() {
	testCase := []struct {
		Name        string
		Body        model.User
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name: "success",
			Body: model.User{
				Email:    "test@gmail.com",
				Password: "test",
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "error",
			Body: model.User{
				Email:    "test@gmail.com",
				Password: "test",
			},
			ExpectedErr: errors.New("internal error"),
			MockReturn:  errors.New("internal error"),
		},
		{
			Name: "error",
			Body: model.User{
				Email:    "test@gmail.com",
				Password: "test",
			},
			ExpectedErr: utils.ErrEmailAlredyExist,
			MockReturn:  errors.New("Duplicate entry"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`id`,`created_at`,`updated_at`,`deleted_at`,`name`,`email`,`phone`,`password`,`role_id`,`province_id`,`regency_id`,`district_id`,`village_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}
			var ctx context.Context
			err := s.repository.CreateUser(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)

		})
	}
}
func (s *suiteUserRepository) TestFindUserByEmail() {
	testCase := []struct {
		Name           string
		Email          string
		ExpectedErr    error
		MockReturn     error
		ExpeckedResult *model.User
	}{
		{
			Name:        "success",
			Email:       "test@gmail.com",
			ExpectedErr: nil,
			MockReturn:  nil,
			ExpeckedResult: &model.User{
				Email: "test@gmail.com",
			},
		},
		{
			Name:           "not found",
			Email:          "test@gmail.com",
			ExpectedErr:    utils.ErrNotFound,
			MockReturn:     gorm.ErrRecordNotFound,
			ExpeckedResult: nil,
		},
		{
			Name:           "error database",
			Email:          "test@gmail.com",
			ExpectedErr:    errors.New("err"),
			MockReturn:     errors.New("err"),
			ExpeckedResult: nil,
		},
	}
	row := sqlmock.NewRows([]string{"email"}).AddRow("test@gmail.com")
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `id`,`email`,`password`,`role_id` FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
			} else {
				db.WillReturnRows(row)
			}
			var ctx context.Context
			res, err := s.repository.FindUserByEmail(v.Email, ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpeckedResult, res)
		})
	}
}
func (s *suiteUserRepository) TestFindUserById() {
	testCase := []struct {
		Name        string
		ID          string
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name:        "success",
			ID:          "123",
			ExpectedErr: nil,
		},
		{
			Name:        "not found",
			ID:          "123",
			ExpectedErr: utils.ErrNotFound,
			MockReturn:  gorm.ErrRecordNotFound,
		},
		{
			Name:        "error database",
			ID:          "123",
			ExpectedErr: errors.New("err"),
			MockReturn:  errors.New("err"),
		},
	}
	row := sqlmock.NewRows([]string{"email"}).AddRow("test@gmail.com")
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
			} else {
				db.WillReturnRows(row)
			}
			var ctx context.Context
			_, err := s.repository.FindUserByID(v.ID, ctx)

			s.Equal(v.ExpectedErr, err)
		})
	}
}
func (s *suiteUserRepository) TestFindAllUser() {
	testCase := []struct {
		Name           string
		Email          string
		ExpectedErr    error
		MockReturn     error
		ExpeckedResult []model.User
	}{
		{
			Name:        "success",
			Email:       "test@gmail.com",
			ExpectedErr: nil,
			MockReturn:  nil,
			ExpeckedResult: []model.User{
				{
					Email: "test@gmail.com",
				},
			},
		},
		{
			Name:           "error",
			Email:          "test@gmail.com",
			ExpectedErr:    errors.New("error"),
			MockReturn:     errors.New("error"),
			ExpeckedResult: nil,
		},
	}
	row := sqlmock.NewRows([]string{"email"}).AddRow("test@gmail.com")
	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			db := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
			} else {
				db.WillReturnRows(row)
			}
			var ctx context.Context
			res, err := s.repository.FindAllUsers(ctx)

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpeckedResult, res)
		})
	}
}

func (s *suiteUserRepository) TestUpdateUser() {
	testCase := []struct {
		Name        string
		Body        model.User
		RowAffected int64
		ExpectedErr error
		MockErr     error
	}{
		{
			Name: "success",
			Body: model.User{
				ID:    uuid.New(),
				Email: "test@gmail.com",
			},
			RowAffected: 1,
			ExpectedErr: nil,
			MockErr:     nil,
		},
		{
			Name: "user not found",
			Body: model.User{
				ID:    uuid.New(),
				Email: "test@gmail.com",
			},
			RowAffected: 0,
			ExpectedErr: utils.ErrNotFound,
			MockErr:     nil,
		},
		{
			Name: "error",
			Body: model.User{
				ID:    uuid.New(),
				Email: "test@gmail.com",
			},
			RowAffected: 0,
			ExpectedErr: errors.New("error"),
			MockErr:     errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL"))
			if v.MockErr != nil {
				db.WillReturnError(v.MockErr)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(0, v.RowAffected))
				s.mock.ExpectCommit()
			}
			var ctx context.Context
			err := s.repository.UpdateUser(&v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
		})
	}
}

func (s *suiteUserRepository) TestDeleteUser() {
	testCase := []struct {
		Name        string
		ExpectedErr error
		RowAffected int64
		MockErr     error
		Body        *model.User
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			RowAffected: 1,
			MockErr:     nil,
			Body: &model.User{
				ID: uuid.New(),
			},
		},
		{
			Name:        "user not found",
			ExpectedErr: utils.ErrNotFound,
			RowAffected: 0,
			MockErr:     nil,
			Body: &model.User{
				ID: uuid.New(),
			},
		},
		{
			Name:        "error database went find user",
			ExpectedErr: errors.New("err"),
			RowAffected: 0,
			MockErr:     errors.New("err"),
			Body: &model.User{
				ID: uuid.New(),
			},
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"))
			if v.MockErr != nil {
				db.WillReturnError(v.MockErr)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(0, v.RowAffected))
				s.mock.ExpectCommit()
			}

			var ctx context.Context
			err := s.repository.DeleteUser(v.Body, ctx)

			s.Equal(v.ExpectedErr, err)
		})
	}
}

func TestSuiteUserRepository(t *testing.T) {
	suite.Run(t, new(suiteUserRepository))
}
