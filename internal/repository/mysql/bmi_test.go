package mysql_test

import (
	"context"
	repository "github.com/bxcodec/go-clean-arch/internal/repository/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bxcodec/go-clean-arch/domain"
)

func TestBMIRepository_Store(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := "INSERT INTO bmi_records(height, weight, value, created_at) VALUES(?,?,?,NOW())"
	bmi := &domain.BMI{
		Height: 1.70,
		Weight: 70.0,
		Value:  24.221453287197235,
	}

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(bmi.Height, bmi.Weight, bmi.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewBMIRepository(db)
	err = repo.Store(context.Background(), bmi)
	require.NoError(t, err)

	assert.Equal(t, int64(1), bmi.ID)
}

func TestBMIRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := "SELECT id, height, weight, value, created_at FROM bmi_records WHERE id = ?"
	id := int64(1)
	bmiValue := 24.221453287197235
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "height", "weight", "value", "created_at"}).
		AddRow(id, 1.75, 70.0, bmiValue, createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(id).
		WillReturnRows(rows)

	repo := repository.NewBMIRepository(db)
	bmi, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)

	assert.NotNil(t, bmi)
	assert.Equal(t, id, bmi.ID)
	assert.Equal(t, 1.75, bmi.Height)
	assert.Equal(t, 70.0, bmi.Weight)
	assert.Equal(t, bmiValue, bmi.Value)
	assert.WithinDuration(t, createdAt, time.Now(), time.Second)
}

func TestBMIRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := "SELECT id, height, weight, value, created_at FROM bmi_records"

	createdAt1 := time.Now().Add(-1 * time.Hour)
	createdAt2 := time.Now().Add(-2 * time.Hour)

	rows := sqlmock.NewRows([]string{"id", "height", "weight", "value", "created_at"}).
		AddRow(1, 1.75, 70.0, 22.857142857142858, createdAt1).
		AddRow(2, 1.80, 75.0, 23.148148148148145, createdAt2)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(rows)

	repo := repository.NewBMIRepository(db)
	bmis, err := repo.GetAll(context.Background())
	require.NoError(t, err)

	assert.Len(t, bmis, 2)
	assert.Equal(t, int64(1), bmis[0].ID)
	assert.Equal(t, 1.75, bmis[0].Height)
	assert.Equal(t, 70.0, bmis[0].Weight)
	assert.Equal(t, 22.857142857142858, bmis[0].Value)
	assert.WithinDuration(t, createdAt1, bmis[0].CreatedAt, time.Second)

	assert.Equal(t, int64(2), bmis[1].ID)
	assert.Equal(t, 1.80, bmis[1].Height)
	assert.Equal(t, 75.0, bmis[1].Weight)
	assert.Equal(t, 23.148148148148145, bmis[1].Value)
	assert.WithinDuration(t, createdAt2, bmis[1].CreatedAt, time.Second)
}

func TestBMIRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `UPDATE bmi_records SET height = ?, weight = ?, value = ? WHERE id = ?`
	bmi := &domain.BMI{
		ID:     1,
		Height: 1.75,
		Weight: 75.0,
		Value:  24.49,
	}

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(bmi.Height, bmi.Weight, bmi.Value, bmi.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := repository.NewBMIRepository(db)
	err = repo.Update(context.Background(), bmi)
	require.NoError(t, err)
}

func TestBMIRepository_Update_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `UPDATE bmi_records SET height = ?, weight = ?, value = ? WHERE id = ?`
	bmi := &domain.BMI{
		ID:     999, // Non-existent ID
		Height: 1.75,
		Weight: 75.0,
		Value:  24.49,
	}

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(bmi.Height, bmi.Weight, bmi.Value, bmi.ID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

	repo := repository.NewBMIRepository(db)
	err = repo.Update(context.Background(), bmi)
	require.Error(t, err)
	assert.Equal(t, "no record found to update", err.Error())
}

func TestBMIRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `DELETE FROM bmi_records WHERE id = ?`
	id := int64(1)

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := repository.NewBMIRepository(db)
	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
}

func TestBMIRepository_Delete_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `DELETE FROM bmi_records WHERE id = ?`
	id := int64(999)

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := repository.NewBMIRepository(db)
	err = repo.Delete(context.Background(), id)
	require.Error(t, err)
	assert.Equal(t, "no record found to delete", err.Error())
}
