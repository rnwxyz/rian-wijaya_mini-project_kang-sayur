package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/repository"
	urp "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/user/repository"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/utils"
)

type checkpointServiceImpl struct {
	repo     repository.CheckpointRepository
	userRepo urp.UserRepository
}

// FindCheckpointsByUser implements CheckpointService
func (s *checkpointServiceImpl) FindCheckpointsByUser(id string, ctx context.Context) (dto.CheckpointsResponse, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, utils.ErrInvalidId
	}
	user, _ := s.userRepo.FindUserByID(id, ctx)
	var checkpointsResponse dto.CheckpointsResponse

	checkpoints1, err := s.repo.FindCheckpointByVilage(*user, ctx)
	if err != nil {
		return nil, err
	}
	if len(checkpoints1) != 0 {
		checkpointsResponse.FromModel(checkpoints1)
		return checkpointsResponse, nil
	}

	checkpoints2, err := s.repo.FindCheckpointByDistrict(*user, ctx)
	if err != nil {
		return nil, err
	}
	if len(checkpoints2) != 0 {
		checkpointsResponse.FromModel(checkpoints2)
		return checkpointsResponse, nil
	}

	checkpoints3, err := s.repo.FindCheckpointByRegency(*user, ctx)
	if err != nil {
		return nil, err
	}
	if len(checkpoints3) != 0 {
		checkpointsResponse.FromModel(checkpoints3)
		return checkpointsResponse, nil
	}

	checkpoints4, err := s.repo.FindCheckpointByProvince(*user, ctx)
	if err != nil {
		return nil, err
	}

	if len(checkpoints4) != 0 {
		checkpointsResponse.FromModel(checkpoints4)
		return checkpointsResponse, nil
	}

	return nil, utils.ErrCheckpointNotCovered
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

func NewCheckpointService(repository repository.CheckpointRepository, userRepo urp.UserRepository) CheckpointService {
	return &checkpointServiceImpl{
		repo:     repository,
		userRepo: userRepo,
	}
}
