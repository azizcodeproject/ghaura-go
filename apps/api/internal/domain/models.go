package domain

import "time"

// OwnershipType membedakan armada/driver milik sendiri vs partner.
type OwnershipType string

const (
	OwnershipOwned   OwnershipType = "OWNED"
	OwnershipPartner OwnershipType = "PARTNER"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	NPWP      string    `json:"npwp,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type Vehicle struct {
	ID            string        `json:"id"`
	PlateNumber   string        `json:"plateNumber"`
	OwnershipType OwnershipType `json:"ownershipType"`
	PartnerName   string        `json:"partnerName,omitempty"`
}

type Driver struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Phone         string        `json:"phone,omitempty"`
	OwnershipType OwnershipType `json:"ownershipType"`
}

type ShipmentStatus string

const (
	ShipmentDraft     ShipmentStatus = "DRAFT"
	ShipmentBooked    ShipmentStatus = "BOOKED"
	ShipmentInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentDelivered ShipmentStatus = "DELIVERED"
	ShipmentCancelled ShipmentStatus = "CANCELLED"
)

type Shipment struct {
	ID         string         `json:"id"`
	ResiNumber string         `json:"resiNumber"`
	CustomerID string         `json:"customerId"`
	Status     ShipmentStatus `json:"status"`
	Origin     string         `json:"origin"`
	Destination string        `json:"destination"`
	CreatedAt  time.Time      `json:"createdAt"`
}

type Assignment struct {
	ID         string    `json:"id"`
	ShipmentID string    `json:"shipmentId"`
	DriverID   string    `json:"driverId"`
	VehicleID  string    `json:"vehicleId"`
	AssignedAt time.Time `json:"assignedAt"`
}

type TrackingEvent struct {
	ID         string    `json:"id"`
	ShipmentID string    `json:"shipmentId"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	StatusNote string    `json:"statusNote,omitempty"`
	RecordedAt time.Time `json:"recordedAt"`
}

type POD struct {
	ID           string    `json:"id"`
	ShipmentID   string    `json:"shipmentId"`
	ReceiverName string    `json:"receiverName"`
	PhotoURL     string    `json:"photoUrl,omitempty"`
	Lat          float64   `json:"lat,omitempty"`
	Lng          float64   `json:"lng,omitempty"`
	SignedAt     time.Time `json:"signedAt"`
}
