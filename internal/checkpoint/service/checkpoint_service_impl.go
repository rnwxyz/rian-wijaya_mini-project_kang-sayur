package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type checkpointServiceImpl struct {
	repo repository.CheckpointRepository
}

// CreateCheckpoint implements CheckpointService
func (s *checkpointServiceImpl) CreateCheckpoint(body dto.CheckpointRequest, ctx context.Context) (uuid.UUID, error) {
	newId := uuid.New()
	checkpoint := body.ToModel()
	checkpoint.ID = newId
	err := s.repo.CreateCheckpoint(checkpoint, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot add or update a child row") {
			return uuid.Nil, utils.ErrBadRequestBody
		}
		return uuid.Nil, err
	}
	return newId, err
}

// DeleteCheckpoint implements CheckpointService
func (s *checkpointServiceImpl) DeleteCheckpoint(id string, ctx context.Context) error {
	panic("unimplemented")
}

// FindCheckpoints implements CheckpointService
func (s *checkpointServiceImpl) FindCheckpoints(ctx context.Context) (dto.CheckpointsResponse, error) {
	checkpoints, err := s.repo.FindCheckpoints(ctx)
	if err != nil {
		return nil, err
	}
	var checkpointsResponse dto.CheckpointsResponse
	checkpointsResponse.FromModel(checkpoints)
	return checkpointsResponse, nil
}

// UpdateCheckpoint implements CheckpointService
func (s *checkpointServiceImpl) UpdateCheckpoint(id string, body dto.CheckpointRequest, ctx context.Context) error {
	panic("unimplemented")
}

func NewCheckpointService(repository repository.CheckpointRepository) CheckpointService {
	return &checkpointServiceImpl{
		repo: repository,
	}
}
