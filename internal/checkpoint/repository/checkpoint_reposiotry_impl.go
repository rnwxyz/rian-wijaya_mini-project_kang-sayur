package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"gorm.io/gorm"
)

type checkpointRepositoryImpl struct {
	db *gorm.DB
}

// FindCheckpointByDistrict implements CheckpointRepository
func (r *checkpointRepositoryImpl) FindCheckpointByDistrict(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	var checkpoints []model.Checkpoint
	if user.DistrictID == nil {
		return checkpoints, nil
	}
	err := r.db.WithContext(ctx).Where("district_id = ?", *user.DistrictID).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Find(&checkpoints).Error
	if err != nil {
		return nil, err
	}
	return checkpoints, nil
}

// FindCheckpointByProvince implements CheckpointRepository
func (r *checkpointRepositoryImpl) FindCheckpointByProvince(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	var checkpoints []model.Checkpoint
	if user.ProvinceID == nil {
		return checkpoints, nil
	}
	err := r.db.WithContext(ctx).Where("province_id = ?", *user.ProvinceID).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Find(&checkpoints).Error
	if err != nil {
		return nil, err
	}
	return checkpoints, nil
}

// FindCheckpointByRegency implements CheckpointRepository
func (r *checkpointRepositoryImpl) FindCheckpointByRegency(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	var checkpoints []model.Checkpoint
	if user.RegencyID == nil {
		return checkpoints, nil
	}
	err := r.db.WithContext(ctx).Where("regency_id = ?", *user.RegencyID).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Find(&checkpoints).Error
	if err != nil {
		return nil, err
	}
	return checkpoints, nil
}

// FindCheckpointByVilage implements CheckpointRepository
func (r *checkpointRepositoryImpl) FindCheckpointByVilage(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	var checkpoints []model.Checkpoint
	if user.VillageID == nil {
		return checkpoints, nil
	}
	err := r.db.WithContext(ctx).Where("village_id = ?", *user.VillageID).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Find(&checkpoints).Error
	if err != nil {
		return nil, err
	}
	return checkpoints, nil
}

// CreateCheckpoint implements CheckpointRepository
func (r *checkpointRepositoryImpl) CreateCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(checkpoint).Error
	return err
}

// DeleteCheckpoint implements CheckpointRepository
func (r *checkpointRepositoryImpl) DeleteCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error {
	panic("unimplemented")
}

// FindCheckpoints implements CheckpointRepository
func (r *checkpointRepositoryImpl) FindCheckpoints(ctx context.Context) ([]model.Checkpoint, error) {
	var checkpoints []model.Checkpoint
	err := r.db.WithContext(ctx).Preload("Province").Preload("Regency").Preload("District").Preload("Village").Find(&checkpoints).Error
	if err != nil {
		return nil, err
	}
	return checkpoints, nil
}

// UpdateCheckpoint implements CheckpointRepository
func (r *checkpointRepositoryImpl) UpdateCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error {
	panic("unimplemented")
}

func NewCheckpointRepository(db *gorm.DB) CheckpointRepository {
	return &checkpointRepositoryImpl{
		db: db,
	}
}
