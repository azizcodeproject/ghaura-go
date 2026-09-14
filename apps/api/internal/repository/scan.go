package repository

import (
	"database/sql"

	"github.com/azizcodeproject/ghaura-go/apps/api/internal/domain"
)

const shipmentSelectColumns = `
	s.id,
	s.resi_number,
	s.customer_id,
	c.name,
	s.status,
	s.origin,
	s.destination,
	s.created_by,
	s.created_at,
	s.updated_at
`

func scanCustomer(scanner interface {
	Scan(dest ...any) error
}) (domain.Customer, error) {
	var customer domain.Customer
	var npwp sql.NullString
	err := scanner.Scan(
		&customer.ID,
		&customer.Name,
		&npwp,
		&customer.CreatedBy,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return domain.Customer{}, err
	}
	if npwp.Valid {
		customer.NPWP = npwp.String
	}
	return customer, nil
}

func scanShipment(scanner interface {
	Scan(dest ...any) error
}) (domain.Shipment, error) {
	var shipment domain.Shipment
	err := scanner.Scan(
		&shipment.ID,
		&shipment.ResiNumber,
		&shipment.CustomerID,
		&shipment.CustomerName,
		&shipment.Status,
		&shipment.Origin,
		&shipment.Destination,
		&shipment.CreatedBy,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)
	if err != nil {
		return domain.Shipment{}, err
	}
	return shipment, nil
}
