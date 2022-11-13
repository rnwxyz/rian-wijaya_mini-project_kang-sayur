package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/internal/checkpoint/dto"
	"github.com/stretchr/testify/mock"
)

type CheckpointServiceMock struct {
	mock.Mock
}

func (b *CheckpointServiceMock) CreateCheckpoint(body dto.CheckpointRequest, ctx context.Context) (uuid.UUID, error) {
	args := b.Called()
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (b *CheckpointServiceMock) FindCheckpoints(ctx context.Context) (dto.CheckpointsResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.CheckpointsResponse), args.Error(1)

}

func (b *CheckpointServiceMock) FindCheckpointsByUser(id string, ctx context.Context) (dto.CheckpointsResponse, error) {
	args := b.Called()
	return args.Get(0).(dto.CheckpointsResponse), args.Error(1)
}
