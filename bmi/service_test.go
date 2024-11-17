package bmi_test

import (
	"context"
	"github.com/bxcodec/go-clean-arch/bmi"
	"github.com/bxcodec/go-clean-arch/bmi/mocks"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCalculateAndStoreBMI(t *testing.T) {
	mockRepo := new(mocks.MockBMIRepository)
	service := bmi.NewServices(mockRepo)

	height := 1.70
	weight := 70.0
	expectedValue := 24.221453287197235

	expectedBMI := &domain.BMI{
		Height:    height,
		Weight:    weight,
		Value:     expectedValue,
		CreatedAt: time.Now(),
	}

	mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.BMI")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*domain.BMI)
		arg.ID = 1
	})

	ctx := context.Background()
	bmiResult, err := service.CalculateAndStoreBMI(ctx, height, weight)

	assert.NoError(t, err)
	assert.NotNil(t, bmiResult)
	assert.Equal(t, expectedBMI.Height, bmiResult.Height)
	assert.Equal(t, expectedBMI.Weight, bmiResult.Weight)
	assert.InDelta(t, expectedBMI.Value, bmiResult.Value, 0.001)
	mockRepo.AssertExpectations(t)
}

func TestCalculateBMICategoryAndRisk(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected struct {
			category string
			risk     string
		}
	}{
		{
			name:  "Underweight",
			value: 17.5,
			expected: struct {
				category string
				risk     string
			}{
				category: "น้ำหนักน้อย / ผอม",
				risk:     "มากกว่าคนปกติ",
			},
		},
		{
			name:  "Normal weight",
			value: 21.0,
			expected: struct {
				category string
				risk     string
			}{
				category: "ปกติ (สุขภาพดี)",
				risk:     "เท่าคนปกติ",
			},
		},
		{
			name:  "Overweight",
			value: 24.5,
			expected: struct {
				category string
				risk     string
			}{
				category: "ท้วม / โรคอ้วนระดับ 1",
				risk:     "อันตรายระดับ 1",
			},
		},
		{
			name:  "Obese level 2",
			value: 27.0,
			expected: struct {
				category string
				risk     string
			}{
				category: "อ้วน / โรคอ้วนระดับ 2",
				risk:     "อันตรายระดับ 2",
			},
		},
		{
			name:  "Obese level 3",
			value: 31.0,
			expected: struct {
				category string
				risk     string
			}{
				category: "อ้วนมาก / โรคอ้วนระดับ 3",
				risk:     "อันตรายระดับ 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, risk := bmi.CalculateBMICategoryAndRisk(tt.value)
			assert.Equal(t, tt.expected.category, category)
			assert.Equal(t, tt.expected.risk, risk)
		})
	}
}

func TestGetBMIByID(t *testing.T) {
	mockRepo := new(mocks.MockBMIRepository)
	service := bmi.NewServices(mockRepo)

	expectedBMI := &domain.BMI{
		ID:        1,
		Height:    1.70,
		Weight:    70.0,
		Value:     24.22,
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedBMI, nil)

	ctx := context.Background()
	result, err := service.GetBMIByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBMI.ID, result.ID)
	assert.Equal(t, "ท้วม / โรคอ้วนระดับ 1", result.Category)
	assert.Equal(t, "อันตรายระดับ 1", result.Risk)

	mockRepo.AssertExpectations(t)
}

func TestGetAllBMI(t *testing.T) {
	mockRepo := new(mocks.MockBMIRepository)
	service := bmi.NewServices(mockRepo)

	createdAt1 := time.Now()
	createdAt2 := time.Now().Add(-time.Hour)

	expectedBMIs := []*domain.BMI{
		{
			ID:        1,
			Height:    1.70,
			Weight:    70.0,
			Value:     24.22,
			CreatedAt: createdAt1,
		},
		{
			ID:        2,
			Height:    1.80,
			Weight:    75.0,
			Value:     23.15,
			CreatedAt: createdAt2,
		},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expectedBMIs, nil)

	ctx := context.Background()
	result, err := service.GetAllBMI(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(expectedBMIs), len(result))

	assert.Equal(t, "ท้วม / โรคอ้วนระดับ 1", result[0].Category)
	assert.Equal(t, "อันตรายระดับ 1", result[0].Risk)
	assert.Equal(t, "ท้วม / โรคอ้วนระดับ 1", result[1].Category)
	assert.Equal(t, "อันตรายระดับ 1", result[1].Risk)

	mockRepo.AssertExpectations(t)
}

func TestUpdateBMI(t *testing.T) {
	mockRepo := new(mocks.MockBMIRepository)
	service := bmi.NewServices(mockRepo)

	bmiToUpdate := &domain.BMI{
		ID:     1,
		Height: 1.75,
		Weight: 72.0,
		Value:  0,
	}

	expectedValue := 23.51
	bmiToUpdate.Value = expectedValue

	mockRepo.On("Update", mock.Anything, bmiToUpdate).Return(nil)

	ctx := context.Background()
	err := service.UpdateBMI(ctx, bmiToUpdate)

	assert.NoError(t, err)
	assert.InDelta(t, expectedValue, bmiToUpdate.Value, 0.001)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBMI(t *testing.T) {
	mockRepo := new(mocks.MockBMIRepository)
	service := bmi.NewServices(mockRepo)

	idToDelete := int64(1)
	mockRepo.On("Delete", mock.Anything, idToDelete).Return(nil)

	ctx := context.Background()
	err := service.DeleteBMI(ctx, idToDelete)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
