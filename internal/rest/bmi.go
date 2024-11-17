package rest

import (
	"context"
	"errors"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BmiService interface {
	CalculateAndStoreBMI(ctx context.Context, height, weight float64) (*domain.BMI, error)
	GetBMIByID(ctx context.Context, id int64) (*domain.BMI, error)
	GetAllBMI(ctx context.Context) ([]*domain.BMI, error)
	UpdateBMI(ctx context.Context, bmi *domain.BMI) error
	DeleteBMI(ctx context.Context, id int64) error
}

type BmiHandler struct {
	BmiSrv BmiService
}

func NewBmiHandler(e *echo.Echo, b BmiService) {
	handler := BmiHandler{
		BmiSrv: b,
	}

	e.POST("/bmi", handler.CalculateAndStoreBMI)
	e.GET("/bmi/:id", handler.GetBMIByID)
	e.GET("/bmi", handler.GetAllBMI)
	e.PUT("/bmi/:id", handler.UpdateBMI)
	e.DELETE("/bmi/:id", handler.DeleteBMI)
}

func (h *BmiHandler) CalculateAndStoreBMI(c echo.Context) error {
	var req domain.BMICalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input. Ensure height and weight are positive numbers."})
	}

	if req.Height <= 0 || req.Weight <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Height and weight must be positive numbers."})
	}

	ctx := c.Request().Context()
	bmi, err := h.BmiSrv.CalculateAndStoreBMI(ctx, req.Height, req.Weight)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bmi)
}

func (h *BmiHandler) GetBMIByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	ctx := c.Request().Context()
	bmi, err := h.BmiSrv.GetBMIByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "BMI record not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bmi)
}

func (h *BmiHandler) GetAllBMI(c echo.Context) error {
	ctx := c.Request().Context()
	bmis, err := h.BmiSrv.GetAllBMI(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bmis)
}

func (h *BmiHandler) UpdateBMI(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var req domain.BMI
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	req.ID = id

	ctx := c.Request().Context()
	if err := h.BmiSrv.UpdateBMI(ctx, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "BMI record updated successfully"})
}

func (h *BmiHandler) DeleteBMI(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	ctx := c.Request().Context()
	if err := h.BmiSrv.DeleteBMI(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "BMI record not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "BMI record deleted successfully"})
}
