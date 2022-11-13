package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
)

type CheckpointService interface {
	CreateCheckpoint(body dto.CheckpointRequest, ctx context.Context) (uuid.UUID, error)
	FindCheckpoints(ctx context.Context) (dto.CheckpointsResponse, error)
	FindCheckpointsByUser(id string, ctx context.Context) (dto.CheckpointsResponse, error)
}
