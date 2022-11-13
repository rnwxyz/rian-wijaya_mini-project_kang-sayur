package mock

import (
	"context"

	"github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"
	"github.com/stretchr/testify/mock"
)

type CheckpointRepositoryMock struct {
	mock.Mock
}

func (b *CheckpointRepositoryMock) CreateCheckpoint(checkpoint *model.Checkpoint, ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}

func (b *CheckpointRepositoryMock) FindCheckpoints(ctx context.Context) ([]model.Checkpoint, error) {
	args := b.Called()
	return args.Get(0).([]model.Checkpoint), args.Error(1)
}
func (b *CheckpointRepositoryMock) FindCheckpointByProvince(id model.User, ctx context.Context) ([]model.Checkpoint, error) {
	args := b.Called()
	return args.Get(0).([]model.Checkpoint), args.Error(1)
}

func (b *CheckpointRepositoryMock) FindCheckpointByRegency(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	args := b.Called()
	return args.Get(0).([]model.Checkpoint), args.Error(1)
}

func (b *CheckpointRepositoryMock) FindCheckpointByDistrict(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	args := b.Called()
	return args.Get(0).([]model.Checkpoint), args.Error(1)
}

func (b *CheckpointRepositoryMock) FindCheckpointByVilage(user model.User, ctx context.Context) ([]model.Checkpoint, error) {
	args := b.Called()
	return args.Get(0).([]model.Checkpoint), args.Error(1)
}
