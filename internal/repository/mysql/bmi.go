package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bxcodec/go-clean-arch/domain"
)

type BMIRepository struct {
	Conn *sql.DB
}

func NewBMIRepository(conn *sql.DB) *BMIRepository {
	return &BMIRepository{
		Conn: conn,
	}
}

func (m *BMIRepository) Store(ctx context.Context, bmi *domain.BMI) error {
	query := `INSERT INTO bmi_records(height, weight, value, created_at) VALUES(?,?,?,NOW())`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, bmi.Height, bmi.Weight, bmi.Value)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	bmi.ID = lastID
	return nil
}

func (m *BMIRepository) GetByID(ctx context.Context, id int64) (*domain.BMI, error) {
	query := `SELECT id, height, weight, value, created_at FROM bmi_records WHERE id = ?`
	row := m.Conn.QueryRowContext(ctx, query, id)

	bmi := &domain.BMI{}
	err := row.Scan(&bmi.ID, &bmi.Height, &bmi.Weight, &bmi.Value, &bmi.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return bmi, nil
}

func (m *BMIRepository) GetAll(ctx context.Context) ([]*domain.BMI, error) {
	query := `SELECT id, height, weight, value, created_at FROM bmi_records`
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bmis []*domain.BMI
	for rows.Next() {
		bmi := &domain.BMI{}
		if err := rows.Scan(&bmi.ID, &bmi.Height, &bmi.Weight, &bmi.Value, &bmi.CreatedAt); err != nil {
			return nil, err
		}
		bmis = append(bmis, bmi)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bmis, nil
}

func (m *BMIRepository) Update(ctx context.Context, bmi *domain.BMI) error {
	query := `UPDATE bmi_records SET height = ?, weight = ?, value = ? WHERE id = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, bmi.Height, bmi.Weight, bmi.Value, bmi.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no record found to update")
	}
	return nil
}

func (m *BMIRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bmi_records WHERE id = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no record found to delete")
	}
	return nil
}
