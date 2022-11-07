package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"gorm.io/gorm"
)

type checkpointRepositoryImpl struct {
	db *gorm.DB
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
