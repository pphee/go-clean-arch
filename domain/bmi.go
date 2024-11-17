package domain

import "time"

type BMI struct {
	ID        int64     `json:"id"`
	Height    float64   `json:"height"`
	Weight    float64   `json:"weight"`
	Value     float64   `json:"value"`
	Category  string    `json:"category,omitempty"`
	Risk      string    `json:"risk,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type BMICalculationRequest struct {
	Height float64 `json:"height" validate:"required,gt=0"`
	Weight float64 `json:"weight" validate:"required,gt=0"`
}
