package qdrantrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	client "github.com/qdrant/go-client/qdrant"
)

type BMIRepository struct {
	client         *client.Client
	collectionName string
}

func NewBMIRepository(endpoint, apiKey, collectionName string) (*BMIRepository, error) {
	qdrantClient, err := client.NewClient(&client.Config{
		Host:   endpoint,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Qdrant client: %w", err)
	}

	return &BMIRepository{
		client:         qdrantClient,
		collectionName: collectionName,
	}, nil
}

func (r *BMIRepository) CreateCollection(ctx context.Context) error {
	err := r.client.CreateCollection(ctx, &client.CreateCollection{
		CollectionName: r.collectionName,
		VectorsConfig: client.NewVectorsConfig(&client.VectorParams{
			Size:     3,
			Distance: client.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil
}

func (r *BMIRepository) Store(ctx context.Context, bmi *domain.BMI) error {
	vector := client.NewVectors(float32(bmi.Height), float32(bmi.Weight), float32(bmi.Value))

	payload := client.NewValueMap(map[string]any{
		"category":   bmi.Category,
		"risk":       bmi.Risk,
		"created_at": bmi.CreatedAt.Format(time.RFC3339),
	})

	point := &client.PointStruct{
		Id:      client.NewIDNum(uint64(bmi.ID)),
		Vectors: vector,
		Payload: payload,
	}

	_, err := r.client.Upsert(ctx, &client.UpsertPoints{
		CollectionName: r.collectionName,
		Points:         []*client.PointStruct{point},
	})
	if err != nil {
		return fmt.Errorf("failed to upsert point: %w", err)
	}
	return nil
}

func (r *BMIRepository) Query(ctx context.Context, queryVector []float32) ([]*client.ScoredPoint, error) {
	limit := uint64(10)

	response, err := r.client.Query(ctx, &client.QueryPoints{
		CollectionName: r.collectionName,
		Query:          client.NewQuery(queryVector...),
		Limit:          &limit,
		WithPayload:    client.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query points: %w", err)
	}

	return response, nil
}
