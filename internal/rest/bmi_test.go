package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/rest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type MockBMIService struct {
	mock.Mock
}

func (m *MockBMIService) CalculateAndStoreBMI(ctx context.Context, height, weight float64) (*domain.BMI, error) {
	args := m.Called(ctx, height, weight)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BMI), args.Error(1)
}

func (m *MockBMIService) GetBMIByID(ctx context.Context, id int64) (*domain.BMI, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BMI), args.Error(1)
}

func (m *MockBMIService) GetAllBMI(ctx context.Context) ([]*domain.BMI, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.BMI), args.Error(1)
}

func (m *MockBMIService) UpdateBMI(ctx context.Context, bmi *domain.BMI) error {
	args := m.Called(ctx, bmi)
	return args.Error(0)
}

func (m *MockBMIService) DeleteBMI(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCalculateAndStoreBMIHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockBMIService)
	handler := &rest.BmiHandler{
		BmiSrv: mockService,
	}

	timestamp := time.Date(2024, time.November, 17, 15, 16, 15, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"height": 1.75,
			"weight": 70.0,
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedBMI := &domain.BMI{
			ID:        1,
			Height:    1.75,
			Weight:    70.0,
			Value:     22.857142857142858,
			CreatedAt: timestamp,
		}
		mockService.On("CalculateAndStoreBMI", mock.Anything, 1.75, 70.0).Return(expectedBMI, nil)

		req := httptest.NewRequest(http.MethodPost, "/bmi", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := handler.CalculateAndStoreBMI(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var respBody domain.BMI
		err = json.Unmarshal(rec.Body.Bytes(), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBMI.Height, respBody.Height)
		assert.Equal(t, expectedBMI.Weight, respBody.Weight)
		assert.InDelta(t, expectedBMI.Value, respBody.Value, 0.001)

		mockService.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"height": -1.75,
			"weight": 70.0,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/bmi", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := handler.CalculateAndStoreBMI(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

}

func TestCalculateAndStoreBMIHandler_ServiceError(t *testing.T) {
	e := echo.New()
	mockService := new(MockBMIService)
	handler := &rest.BmiHandler{
		BmiSrv: mockService,
	}

	t.Run("service error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"height": 1.75,
			"weight": 70.0,
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("CalculateAndStoreBMI", mock.Anything, 1.75, 70.0).Return(nil, errors.New("internal error"))

		req := httptest.NewRequest(http.MethodPost, "/bmi", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := handler.CalculateAndStoreBMI(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal error")

		mockService.AssertExpectations(t)
	})

}

func TestGetBMIByIDHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockBMIService)
	handler := &rest.BmiHandler{BmiSrv: mockService}

	timestamp := time.Date(2024, time.November, 17, 15, 16, 15, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		expectedBMI := &domain.BMI{
			ID:        1,
			Height:    1.75,
			Weight:    70.0,
			Value:     22.857142857142858,
			CreatedAt: timestamp,
		}
		mockService.On("GetBMIByID", mock.Anything, int64(1)).Return(expectedBMI, nil)

		req := httptest.NewRequest(http.MethodGet, "/bmi/1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(1, 10))

		err := handler.GetBMIByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var respBody domain.BMI
		err = json.Unmarshal(rec.Body.Bytes(), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBMI, &respBody)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bmi/invalid", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := handler.GetBMIByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetAllBMIHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockBMIService)
	handler := &rest.BmiHandler{BmiSrv: mockService}

	timestamp := time.Date(2024, time.November, 17, 15, 16, 15, 0, time.UTC)
	t.Run("success", func(t *testing.T) {
		expectedBMIs := []*domain.BMI{
			{ID: 1, Height: 1.75, Weight: 70.0, Value: 22.857142857142858, CreatedAt: timestamp},
			{ID: 2, Height: 1.80, Weight: 75.0, Value: 23.148148148148145, CreatedAt: timestamp},
		}
		mockService.On("GetAllBMI", mock.Anything).Return(expectedBMIs, nil)

		req := httptest.NewRequest(http.MethodGet, "/bmi", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := handler.GetAllBMI(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var respBody []*domain.BMI
		err = json.Unmarshal(rec.Body.Bytes(), &respBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedBMIs, respBody)

		mockService.AssertExpectations(t)
	})
}

func TestDeleteBMIHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockBMIService)
	handler := &rest.BmiHandler{BmiSrv: mockService}

	t.Run("success", func(t *testing.T) {
		mockService.On("DeleteBMI", mock.Anything, int64(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/bmi/1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.DeleteBMI(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "BMI record deleted successfully")

		mockService.AssertExpectations(t)
	})

	t.Run("record not found", func(t *testing.T) {
		mockService.On("DeleteBMI", mock.Anything, int64(999)).Return(domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/bmi/999", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.DeleteBMI(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "BMI record not found")

		mockService.AssertExpectations(t)
	})
}
