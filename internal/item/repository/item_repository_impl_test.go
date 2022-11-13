package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	customerrors "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils/custom_errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type suiteItemRepository struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *itemRepositoryImpl
}

func (s *suiteItemRepository) SetupSuite() {
	db, mocking, _ := sqlmock.New()

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	s.mock = mocking
	s.repository = &itemRepositoryImpl{
		db: dbGorm,
	}
}

func (s *suiteItemRepository) TearDown() {
	s.mock = nil
	s.repository = nil
}

func (s *suiteItemRepository) TestCreateItem() {
	testCase := []struct {
		Name        string
		Body        model.Item
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name: "success create",
			Body: model.Item{
				ID:    1,
				Name:  "item",
				Qty:   1,
				Price: 1,
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "duplicate item",
			Body: model.Item{
				ID:    1,
				Name:  "item",
				Qty:   1,
				Price: 1,
			},
			ExpectedErr: customerrors.ErrDuplicateData,
			MockReturn:  errors.New("Duplicate entry"),
		},
		{
			Name: "invalid foreign key",
			Body: model.Item{
				ID:    1,
				Name:  "item",
				Qty:   1,
				Price: 1,
			},
			ExpectedErr: customerrors.ErrBadRequestBody,
			MockReturn:  errors.New("Cannot add or update a child row"),
		},
		{
			Name: "other error",
			Body: model.Item{
				ID:    1,
				Name:  "item",
				Qty:   1,
				Price: 1,
			},
			ExpectedErr: errors.New("error"),
			MockReturn:  errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `items` (`created_at`,`updated_at`,`deleted_at`,`name`,`category_id`,`description`,`qty`,`price`,`id`) VALUES (?,?,?,?,?,?,?,?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}

			err := s.repository.CreateItem(&v.Body, context.Background())

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}
func (s *suiteItemRepository) TestCreateCategory() {
	testCase := []struct {
		Name        string
		Body        model.Category
		ExpectedErr error
		MockReturn  error
	}{
		{
			Name: "success create",
			Body: model.Category{
				ID:   1,
				Name: "item",
			},
			ExpectedErr: nil,
			MockReturn:  nil,
		},
		{
			Name: "duplicate item",
			Body: model.Category{
				ID:   1,
				Name: "item",
			},
			ExpectedErr: customerrors.ErrDuplicateData,
			MockReturn:  errors.New("Duplicate entry"),
		},
		{
			Name: "other error",
			Body: model.Category{
				ID:   1,
				Name: "item",
			},
			ExpectedErr: errors.New("error"),
			MockReturn:  errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()
			s.mock.ExpectBegin()
			db := s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `categories` (`created_at`,`updated_at`,`deleted_at`,`name`,`description`,`id`) VALUES (?,?,?,?,?,?)"))
			if v.ExpectedErr != nil {
				db.WillReturnError(v.MockReturn)
				s.mock.ExpectRollback()
			} else {
				db.WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			}

			err := s.repository.CreateCategory(&v.Body, context.Background())

			s.Equal(v.ExpectedErr, err)

			s.TearDown()
		})
	}
}

func (s *suiteItemRepository) TestFindCategory() {
	testCase := []struct {
		Name            string
		ExpectedErr     error
		ExpectedRes     []model.Category
		FindCategoryErr error
		FindCategoryRes *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: []model.Category{
				{ID: 1, Name: "categoryA"},
			},
			FindCategoryErr: nil,
			FindCategoryRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "categoryA"),
		},
		{
			Name:            "error",
			ExpectedErr:     errors.New("error"),
			ExpectedRes:     []model.Category(nil),
			FindCategoryErr: errors.New("error"),
			FindCategoryRes: sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`deleted_at` IS NULL")).WillReturnRows(v.FindCategoryRes).WillReturnError(v.FindCategoryErr)

			res, err := s.repository.FindCategories(context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}
func (s *suiteItemRepository) TestFindItemById() {
	testCase := []struct {
		Name        string
		ExpectedErr error
		Item        *model.Item
		ExpectedRes model.Item
		FindItemErr error
		FindItemRes *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			Item: &model.Item{
				ID: 1,
			},
			ExpectedRes: model.Item{
				ID:    1,
				Name:  "itemA",
				Qty:   1,
				Price: 1,
			},
			FindItemErr: nil,
			FindItemRes: sqlmock.NewRows([]string{"id", "name", "qty", "price"}).
				AddRow(1, "itemA", 1, 1),
		},
		{
			Name:        "error",
			ExpectedErr: customerrors.ErrNotFound,
			Item: &model.Item{
				ID: 1,
			},
			ExpectedRes: model.Item{ID: 1},
			FindItemErr: gorm.ErrRecordNotFound,
			FindItemRes: sqlmock.NewRows([]string{"id", "name", "qty", "price"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `items` WHERE `items`.`deleted_at` IS NULL AND `items`.`id` = ? ORDER BY `items`.`id` LIMIT 1")).WillReturnRows(v.FindItemRes).WillReturnError(v.FindItemErr)

			err := s.repository.FindItemById(v.Item, context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(&v.ExpectedRes, v.Item)

			s.TearDown()
		})
	}
}
func (s *suiteItemRepository) TestFindItems() {
	testCase := []struct {
		Name               string
		ExpectedErr        error
		ExpectedRes        []model.Item
		FindItemErr        error
		FindItemRes        *sqlmock.Rows
		PreloadCategoryErr error
		PreloadCategoryRes *sqlmock.Rows
	}{
		{
			Name:        "success",
			ExpectedErr: nil,
			ExpectedRes: []model.Item{
				{
					ID:         1,
					Name:       "item",
					CategoryID: 1,
					Category: model.Category{
						ID:   1,
						Name: "category",
					},
				},
			},
			FindItemErr: nil,
			FindItemRes: sqlmock.NewRows([]string{"id", "name", "category_id"}).
				AddRow(1, "item", 1),
			PreloadCategoryErr: nil,
			PreloadCategoryRes: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "category"),
		},
		{
			Name:               "error",
			ExpectedErr:        errors.New("error"),
			ExpectedRes:        []model.Item(nil),
			FindItemErr:        errors.New("error"),
			FindItemRes:        sqlmock.NewRows([]string{"id", "name", "category_id"}),
			PreloadCategoryErr: nil,
			PreloadCategoryRes: sqlmock.NewRows([]string{"id", "name"}),
		},
	}

	for _, v := range testCase {
		s.T().Run(v.Name, func(t *testing.T) {
			s.SetupSuite()

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `items` WHERE `items`.`deleted_at` IS NULL")).WillReturnRows(v.FindItemRes).WillReturnError(v.FindItemErr)

			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL")).WillReturnRows(v.PreloadCategoryRes).WillReturnError(v.PreloadCategoryErr)

			res, err := s.repository.FindItems(context.Background())

			s.Equal(v.ExpectedErr, err)
			s.Equal(v.ExpectedRes, res)

			s.TearDown()
		})
	}
}

func TestSuiteItemRepository(t *testing.T) {
	suite.Run(t, new(suiteItemRepository))
}
