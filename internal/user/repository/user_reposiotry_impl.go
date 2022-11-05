package repository

import (
	"context"
	"strings"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/config"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/constants"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// CreateUser implements UserRepository
func (u *userRepositoryImpl) CreateUser(user *model.User, ctx context.Context) error {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrEmailAlredyExist
		}
		return err
	}
	return nil
}

// FindUserByEmail implements UserRepository
func (u *userRepositoryImpl) FindUserByEmail(email string, ctx context.Context) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Select([]string{"id", "email", "password", "role_id"}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByID implements UserRepository
func (u *userRepositoryImpl) FindUserByID(id string, ctx context.Context) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindAllUsers implements UserRepository
func (u *userRepositoryImpl) FindAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := u.db.WithContext(ctx).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Find(&users).Error
	return users, err
}

// UpdateUser implements UserRepository
func (u *userRepositoryImpl) UpdateUser(user *model.User, ctx context.Context) error {
	res := u.db.WithContext(ctx).Model(model.User{}).Where("id = ?", user.ID).Updates(&model.User{
		Name:       user.Name,
		Phone:      user.Phone,
		ProvinceID: user.ProvinceID,
		RegencyID:  user.RegencyID,
		DistrictID: user.DistrictID,
		VillageID:  user.VillageID,
	})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return utils.ErrNotFound
	}
	return nil
}

// DeleteUser implements UserRepository
func (u *userRepositoryImpl) DeleteUser(user *model.User, ctx context.Context) error {
	err := u.db.WithContext(ctx).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return err
	}
	if user.Email == config.Cfg.DEFAULT_ADMIN_EMAIL {
		return utils.ErrNotAllowedDeleteDefaultAdmin
	}
	res := u.db.WithContext(ctx).Delete(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return utils.ErrFailedDeleteUser
	}
	return nil
}

// InitRole implements UserRepository
func (u *userRepositoryImpl) InitRole() error {
	var count int64
	err := u.db.Model(&model.Role{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count >= 1 {
		return nil
	}
	role := constants.Role
	err = u.db.Create(&role).Error
	return err
}

func NewUserRepository(db *gorm.DB) UserRepository {
	userRepository := &userRepositoryImpl{
		db: db,
	}
	err := userRepository.InitRole()
	if err != nil {
		panic(err)
	}

	return userRepository
}
