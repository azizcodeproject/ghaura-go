package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/domain"
	"github.com/google/uuid"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, name, npwp, createdBy string) (domain.Customer, error) {
	customerID := uuid.NewString()
	const query = `
		INSERT INTO customers (id, name, npwp, created_by)
		VALUES ($1, $2, NULLIF($3, ''), $4)
		RETURNING id, name, npwp, created_by, created_at, updated_at
	`

	row := r.db.QueryRowContext(ctx, query, customerID, name, npwp, createdBy)
	customer, err := scanCustomer(row)
	if err != nil {
		return domain.Customer{}, fmt.Errorf("insert customer: %w", err)
	}
	return customer, nil
}

func (r *CustomerRepository) List(ctx context.Context) ([]domain.Customer, error) {
	const query = `
		SELECT id, name, npwp, created_by, created_at, updated_at
		FROM customers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	customers := make([]domain.Customer, 0)
	for rows.Next() {
		customer, err := scanCustomer(rows)
		if err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate customers: %w", err)
	}
	return customers, nil
}

func (r *CustomerRepository) Exists(ctx context.Context, customerID string) (bool, error) {
	const query = `
		SELECT 1
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`

	var exists int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check customer: %w", err)
	}
	return true, nil
}
