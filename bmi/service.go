package bmi

import (
	"context"
	"fmt"
	"github.com/bxcodec/go-clean-arch/domain"
	client "github.com/qdrant/go-client/qdrant"
	"time"
)

type bmiRepository interface {
	Store(ctx context.Context, bmi *domain.BMI) error
	GetByID(ctx context.Context, id int64) (*domain.BMI, error)
	GetAll(ctx context.Context) ([]*domain.BMI, error)
	Update(ctx context.Context, bmi *domain.BMI) error
	Delete(ctx context.Context, id int64) error
}

type bmiQdrantRepository interface {
	CreateCollection(ctx context.Context) error
	Store(ctx context.Context, bmi *domain.BMI) error
	Query(ctx context.Context, queryVector []float32) ([]*client.ScoredPoint, error)
}

type Service struct {
	bmiRepo       bmiRepository
	bmiQdrantRepo bmiQdrantRepository
}

func NewServices(b bmiRepository, bq bmiQdrantRepository) *Service {
	return &Service{
		bmiRepo:       b,
		bmiQdrantRepo: bq,
	}
}

func (u *Service) CalculateAndStoreBMI(ctx context.Context, height, weight float64) (*domain.BMI, error) {
	if height <= 0 {
		return nil, fmt.Errorf("height must be greater than 0")
	}
	value := weight / (height * height)
	bmi := &domain.BMI{
		Height:    height,
		Weight:    weight,
		Value:     value,
		CreatedAt: time.Now(),
	}
	err := u.bmiRepo.Store(ctx, bmi)
	if err != nil {
		return nil, err
	}
	return bmi, nil
}

func CalculateBMICategoryAndRisk(value float64) (string, string) {
	switch {
	case value < 18.5:
		return "น้ำหนักน้อย / ผอม", "มากกว่าคนปกติ"
	case value >= 18.5 && value <= 22.9:
		return "ปกติ (สุขภาพดี)", "เท่าคนปกติ"
	case value >= 23 && value <= 24.9:
		return "ท้วม / โรคอ้วนระดับ 1", "อันตรายระดับ 1"
	case value >= 25 && value <= 29.9:
		return "อ้วน / โรคอ้วนระดับ 2", "อันตรายระดับ 2"
	case value > 30:
		return "อ้วนมาก / โรคอ้วนระดับ 3", "อันตรายระดับ 3"
	default:
		return "", ""
	}
}

func CalculateBMI(height, weight float64) (string, string) {
	bmi := weight / (height * height)

	switch {
	case bmi < 18.5:
		return "น้ำหนักน้อย / ผอม", "มากกว่าคนปกติ"
	case bmi >= 18.5 && bmi <= 22.9:
		return "ปกติ (สุขภาพดี)", "เท่าคนปกติ"
	case bmi >= 23 && bmi <= 24.9:
		return "ท้วม / โรคอ้วนระดับ 1", "อันตรายระดับ 1"
	case bmi >= 25 && bmi <= 29.9:
		return "อ้วน / โรคอ้วนระดับ 2", "อันตรายระดับ 2"
	case bmi >= 30:
		return "อ้วนมาก / โรคอ้วนระดับ 3", "อันตรายระดับ 3"
	default:
		return "", ""
	}
}

func (u *Service) GetBMIByID(ctx context.Context, id int64) (*domain.BMI, error) {
	bmi, err := u.bmiRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bmi.Category, bmi.Risk = CalculateBMICategoryAndRisk(bmi.Value)

	return bmi, nil
}

func (u *Service) GetAllBMI(ctx context.Context) ([]*domain.BMI, error) {
	bmiRecords, err := u.bmiRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, bmi := range bmiRecords {
		bmi.Category, bmi.Risk = CalculateBMICategoryAndRisk(bmi.Value)
	}

	return bmiRecords, nil
}

func (u *Service) UpdateBMI(ctx context.Context, bmi *domain.BMI) error {
	if bmi.Height <= 0 || bmi.Weight <= 0 {
		return fmt.Errorf("height and weight must be greater than 0")
	}
	bmi.Value = bmi.Weight / (bmi.Height * bmi.Height)
	return u.bmiRepo.Update(ctx, bmi)
}

func (u *Service) DeleteBMI(ctx context.Context, id int64) error {
	return u.bmiRepo.Delete(ctx, id)
}

func (u *Service) StoreBMI(ctx context.Context, height, weight float64) (*domain.BMI, error) {
	if height <= 0 {
		return nil, fmt.Errorf("height must be greater than 0")
	}

	value := weight / (height * height)
	bmi := &domain.BMI{
		Height:    height,
		Weight:    weight,
		Value:     value,
		CreatedAt: time.Now(),
	}

	bmi.Category, bmi.Risk = CalculateBMICategoryAndRisk(value)

	if err := u.bmiRepo.Store(ctx, bmi); err != nil {
		return nil, fmt.Errorf("failed to store BMI in MySQL: %w", err)
	}

	if err := u.bmiQdrantRepo.Store(ctx, bmi); err != nil {
		return nil, fmt.Errorf("failed to store BMI in Qdrant: %w", err)
	}

	return bmi, nil
}

func (u *Service) QueryBMI(ctx context.Context, queryVector []float32) ([]*client.ScoredPoint, error) {
	results, err := u.bmiQdrantRepo.Query(ctx, queryVector)
	if err != nil {
		return nil, fmt.Errorf("failed to query BMI from Qdrant: %w", err)
	}
	return results, nil
}
