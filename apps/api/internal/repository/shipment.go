package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/domain"
	"github.com/azizcodeproject/ghaura-go/apps/api/internal/resi"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const maxResiAttempts = 5

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Create(ctx context.Context, customerID, origin, destination, createdBy string) (domain.Shipment, error) {
	const query = `
		WITH created AS (
			INSERT INTO shipments (id, resi_number, customer_id, status, origin, destination, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, resi_number, customer_id, status, origin, destination, created_by, created_at, updated_at
		)
		SELECT
			created.id,
			created.resi_number,
			created.customer_id,
			c.name,
			created.status,
			created.origin,
			created.destination,
			created.created_by,
			created.created_at,
			created.updated_at
		FROM created
		JOIN customers c ON c.id = created.customer_id
	`

	var lastErr error
	for attempt := 0; attempt < maxResiAttempts; attempt++ {
		resiNumber, err := resi.Generate()
		if err != nil {
			return domain.Shipment{}, err
		}

		shipmentID := uuid.NewString()
		row := r.db.QueryRowContext(
			ctx,
			query,
			shipmentID,
			resiNumber,
			customerID,
			domain.ShipmentDraft,
			origin,
			destination,
			createdBy,
		)
		shipment, err := scanShipment(row)
		if err == nil {
			return shipment, nil
		}
		if isUniqueViolation(err) {
			lastErr = err
			continue
		}
		return domain.Shipment{}, fmt.Errorf("insert shipment: %w", err)
	}

	return domain.Shipment{}, fmt.Errorf("generate unique resi: %w", lastErr)
}

func (r *ShipmentRepository) List(ctx context.Context, statusFilter string) ([]domain.Shipment, error) {
	query := `
		SELECT ` + shipmentSelectColumns + `
		FROM shipments s
		JOIN customers c ON c.id = s.customer_id
		WHERE s.deleted_at IS NULL
	`
	args := make([]any, 0, 1)
	if statusFilter != "" {
		query += ` AND s.status = $1`
		args = append(args, statusFilter)
	}
	query += ` ORDER BY s.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list shipments: %w", err)
	}
	defer rows.Close()

	shipments := make([]domain.Shipment, 0)
	for rows.Next() {
		shipment, err := scanShipment(rows)
		if err != nil {
			return nil, fmt.Errorf("scan shipment: %w", err)
		}
		shipments = append(shipments, shipment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shipments: %w", err)
	}
	return shipments, nil
}

func (r *ShipmentRepository) GetByID(ctx context.Context, shipmentID string) (domain.Shipment, error) {
	query := `
		SELECT ` + shipmentSelectColumns + `
		FROM shipments s
		JOIN customers c ON c.id = s.customer_id
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`
	return r.getOne(ctx, query, shipmentID)
}

func (r *ShipmentRepository) GetByResi(ctx context.Context, resiNumber string) (domain.Shipment, error) {
	query := `
		SELECT ` + shipmentSelectColumns + `
		FROM shipments s
		JOIN customers c ON c.id = s.customer_id
		WHERE s.resi_number = $1 AND s.deleted_at IS NULL
	`
	return r.getOne(ctx, query, resiNumber)
}

func (r *ShipmentRepository) UpdateStatus(ctx context.Context, shipmentID string, status domain.ShipmentStatus) (domain.Shipment, error) {
	query := `
		UPDATE shipments AS s
		SET status = $2, updated_at = NOW()
		FROM customers AS c
		WHERE s.id = $1
		  AND s.deleted_at IS NULL
		  AND c.id = s.customer_id
		RETURNING ` + shipmentSelectColumns

	return r.getOne(ctx, query, shipmentID, status)
}

func (r *ShipmentRepository) getOne(ctx context.Context, query string, args ...any) (domain.Shipment, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	shipment, err := scanShipment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Shipment{}, sql.ErrNoRows
	}
	if err != nil {
		return domain.Shipment{}, fmt.Errorf("get shipment: %w", err)
	}
	return shipment, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key") || strings.Contains(message, "unique constraint")
}
