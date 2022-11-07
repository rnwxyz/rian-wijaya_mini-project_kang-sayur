package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestSuiteUserRepository(t *testing.T) {
	suite.Run(t, new(suiteUserRepository))
}
