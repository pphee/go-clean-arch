package mocks

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/stretchr/testify/mock"
)

type MockBMIRepository struct {
	mock.Mock
}

func (m *MockBMIRepository) Store(ctx context.Context, bmi *domain.BMI) error {
	args := m.Called(ctx, bmi)
	return args.Error(0)
}

func (m *MockBMIRepository) GetByID(ctx context.Context, id int64) (*domain.BMI, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.BMI), args.Error(1)
}

func (m *MockBMIRepository) GetAll(ctx context.Context) ([]*domain.BMI, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.BMI), args.Error(1)
}

func (m *MockBMIRepository) Update(ctx context.Context, bmi *domain.BMI) error {
	args := m.Called(ctx, bmi)
	return args.Error(0)
}

func (m *MockBMIRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
