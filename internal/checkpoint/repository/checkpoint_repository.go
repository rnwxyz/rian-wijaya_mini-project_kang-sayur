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
}
