package domain

import "testing"

func TestParseShipmentStatus(t *testing.T) {
	validStatuses := []string{"DRAFT", "BOOKED", "IN_TRANSIT", "DELIVERED", "CANCELLED"}
	for _, rawStatus := range validStatuses {
		if _, ok := ParseShipmentStatus(rawStatus); !ok {
			t.Fatalf("expected %s to be valid", rawStatus)
		}
	}

	if _, ok := ParseShipmentStatus("SHIPPED"); ok {
		t.Fatal("expected SHIPPED to be invalid")
	}
}
