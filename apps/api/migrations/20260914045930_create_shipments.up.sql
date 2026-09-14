CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resi_number VARCHAR(32) NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers (id),
    status VARCHAR(32) NOT NULL DEFAULT 'DRAFT',
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    created_by VARCHAR(255) NOT NULL DEFAULT 'system',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_shipments_resi_number UNIQUE (resi_number),
    CONSTRAINT chk_shipments_status CHECK (
        status IN ('DRAFT', 'BOOKED', 'IN_TRANSIT', 'DELIVERED', 'CANCELLED')
    )
);

CREATE INDEX IF NOT EXISTS idx_shipments_customer_id ON shipments (customer_id)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_shipments_status ON shipments (status)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_shipments_resi_number ON shipments (resi_number)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_shipments_deleted_at ON shipments (deleted_at);
