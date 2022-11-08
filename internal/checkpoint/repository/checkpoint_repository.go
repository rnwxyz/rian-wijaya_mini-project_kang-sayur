package repository

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
)

type CheckpointRepository interface {
	CreateCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error
	UpdateCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error
	DeleteCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error
	FindCheckpoints(ctx context.Context) ([]model.Checkpoint, error)
	FindCheckpointByProvince(id model.User, ctx context.Context) ([]model.Checkpoint, error)
	FindCheckpointByRegency(user model.User, ctx context.Context) ([]model.Checkpoint, error)
	FindCheckpointByDistrict(user model.User, ctx context.Context) ([]model.Checkpoint, error)
	FindCheckpointByVilage(user model.User, ctx context.Context) ([]model.Checkpoint, error)
}
