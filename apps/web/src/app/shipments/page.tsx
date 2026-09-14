"use client";

import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useState, type CSSProperties, type FormEvent } from "react";
import { StatusBadge } from "@/components/StatusBadge";
import { createCustomer, createShipment, listCustomers, listShipments } from "@/lib/api";
import { formatDateTime } from "@/lib/format";
import { SHIPMENT_STATUSES, STATUS_LABELS, type ShipmentStatus } from "@/lib/types";

export default function ShipmentsPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<ShipmentStatus | "">("");
  const [resiSearch, setResiSearch] = useState("");
  const [customerName, setCustomerName] = useState("");
  const [customerNpwp, setCustomerNpwp] = useState("");
  const [selectedCustomerId, setSelectedCustomerId] = useState("");
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");

  const customersQuery = useQuery({
    queryKey: ["customers"],
    queryFn: listCustomers,
  });

  const shipmentsQuery = useQuery({
    queryKey: ["shipments", statusFilter],
    queryFn: () => listShipments(statusFilter),
  });

  const createCustomerMutation = useMutation({
    mutationFn: createCustomer,
    onSuccess: (customer) => {
      queryClient.invalidateQueries({ queryKey: ["customers"] });
      setSelectedCustomerId(customer.id);
      setCustomerName("");
      setCustomerNpwp("");
    },
  });

  const createShipmentMutation = useMutation({
    mutationFn: createShipment,
    onSuccess: (shipment) => {
      queryClient.invalidateQueries({ queryKey: ["shipments"] });
      setOrigin("");
      setDestination("");
      router.push(`/shipments/${shipment.id}`);
    },
  });

  const customers = customersQuery.data?.items ?? [];
  const shipments = shipmentsQuery.data?.items ?? [];

  function handleCreateCustomer(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createCustomerMutation.mutate({
      name: customerName.trim(),
      npwp: customerNpwp.trim() || undefined,
    });
  }

  function handleCreateShipment(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createShipmentMutation.mutate({
      customerId: selectedCustomerId,
      origin: origin.trim(),
      destination: destination.trim(),
    });
  }

  function handleSearchResi(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const resiNumber = resiSearch.trim();
    if (resiNumber) {
      router.push(`/shipments/resi/${encodeURIComponent(resiNumber)}`);
    }
  }

  return (
    <main style={{ maxWidth: 1080, margin: "0 auto", padding: 32, display: "grid", gap: 24 }}>
      <div>
        <h1 style={{ margin: 0, fontSize: 28 }}>Resi pengiriman</h1>
        <p style={{ color: "#94a3b8", marginTop: 8 }}>
          Buat resi baru, cek status, dan lihat daftar pengiriman pelanggan B2B.
        </p>
      </div>

      <section style={cardStyle}>
        <h2 style={sectionTitleStyle}>Cari nomor resi</h2>
        <form onSubmit={handleSearchResi} style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
          <input
            value={resiSearch}
            onChange={(event) => setResiSearch(event.target.value)}
            placeholder="Contoh: GHR-260914-A3K7MP"
            style={{ ...inputStyle, flex: "1 1 240px" }}
          />
          <button type="submit" style={buttonStyle}>
            Lihat resi
          </button>
        </form>
      </section>

      <section style={cardStyle}>
        <h2 style={sectionTitleStyle}>Pelanggan baru</h2>
        <form onSubmit={handleCreateCustomer} style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
          <input
            required
            value={customerName}
            onChange={(event) => setCustomerName(event.target.value)}
            placeholder="Nama perusahaan"
            style={inputStyle}
          />
          <input
            value={customerNpwp}
            onChange={(event) => setCustomerNpwp(event.target.value)}
            placeholder="NPWP (opsional)"
            style={inputStyle}
          />
          <button type="submit" disabled={createCustomerMutation.isPending} style={buttonStyle}>
            {createCustomerMutation.isPending ? "Menyimpan…" : "Simpan pelanggan"}
          </button>
        </form>
        {createCustomerMutation.isError && (
          <p style={{ color: "#f87171", marginBottom: 0 }}>
            {createCustomerMutation.error instanceof Error ? createCustomerMutation.error.message : "Gagal menyimpan pelanggan"}
          </p>
        )}
      </section>

      <section style={cardStyle}>
        <h2 style={sectionTitleStyle}>Buat resi</h2>
        <form onSubmit={handleCreateShipment} style={{ display: "grid", gap: 12 }}>
          <select
            required
            value={selectedCustomerId}
            onChange={(event) => setSelectedCustomerId(event.target.value)}
            style={inputStyle}
          >
            <option value="">Pilih pelanggan</option>
            {customers.map((customer) => (
              <option key={customer.id} value={customer.id}>
                {customer.name}
              </option>
            ))}
          </select>
          <div style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
            <input
              required
              value={origin}
              onChange={(event) => setOrigin(event.target.value)}
              placeholder="Asal, misalnya Jakarta"
              style={inputStyle}
            />
            <input
              required
              value={destination}
              onChange={(event) => setDestination(event.target.value)}
              placeholder="Tujuan, misalnya Surabaya"
              style={inputStyle}
            />
          </div>
          <button type="submit" disabled={createShipmentMutation.isPending || !selectedCustomerId} style={buttonStyle}>
            {createShipmentMutation.isPending ? "Membuat resi…" : "Buat resi"}
          </button>
        </form>
        {createShipmentMutation.isError && (
          <p style={{ color: "#f87171", marginBottom: 0 }}>
            {createShipmentMutation.error instanceof Error ? createShipmentMutation.error.message : "Gagal membuat resi"}
          </p>
        )}
      </section>

      <section style={cardStyle}>
        <div style={{ display: "flex", justifyContent: "space-between", gap: 12, flexWrap: "wrap", alignItems: "center" }}>
          <h2 style={{ ...sectionTitleStyle, margin: 0 }}>Daftar resi</h2>
          <select
            value={statusFilter}
            onChange={(event) => setStatusFilter(event.target.value as ShipmentStatus | "")}
            style={{ ...inputStyle, minWidth: 200 }}
          >
            <option value="">Semua status</option>
            {SHIPMENT_STATUSES.map((status) => (
              <option key={status} value={status}>
                {STATUS_LABELS[status]}
              </option>
            ))}
          </select>
        </div>

        {shipmentsQuery.isLoading && <p style={{ color: "#94a3b8" }}>Memuat daftar resi…</p>}
        {shipmentsQuery.isError && (
          <p style={{ color: "#f87171" }}>
            {shipmentsQuery.error instanceof Error ? shipmentsQuery.error.message : "Gagal memuat resi"}
          </p>
        )}
        {!shipmentsQuery.isLoading && shipments.length === 0 && (
          <p style={{ color: "#94a3b8" }}>Belum ada resi. Buat pelanggan lalu terbitkan resi pertama.</p>
        )}

        {shipments.length > 0 && (
          <div style={{ overflowX: "auto" }}>
            <table style={{ width: "100%", borderCollapse: "collapse", fontSize: 14 }}>
              <thead>
                <tr style={{ color: "#94a3b8", textAlign: "left" }}>
                  <th style={thStyle}>Resi</th>
                  <th style={thStyle}>Pelanggan</th>
                  <th style={thStyle}>Rute</th>
                  <th style={thStyle}>Status</th>
                  <th style={thStyle}>Dibuat</th>
                </tr>
              </thead>
              <tbody>
                {shipments.map((shipment) => (
                  <tr key={shipment.id}>
                    <td style={tdStyle}>
                      <Link href={`/shipments/${shipment.id}`} style={{ color: "#38bdf8", textDecoration: "none", fontWeight: 600 }}>
                        {shipment.resiNumber}
                      </Link>
                    </td>
                    <td style={tdStyle}>{shipment.customerName}</td>
                    <td style={tdStyle}>
                      {shipment.origin} → {shipment.destination}
                    </td>
                    <td style={tdStyle}>
                      <StatusBadge status={shipment.status} />
                    </td>
                    <td style={tdStyle}>{formatDateTime(shipment.createdAt)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </main>
  );
}

const cardStyle: CSSProperties = {
  padding: 20,
  borderRadius: 16,
  border: "1px solid #334155",
  background: "#1e293b",
  display: "grid",
  gap: 16,
};

const sectionTitleStyle: CSSProperties = {
  margin: 0,
  fontSize: 18,
};

const inputStyle: CSSProperties = {
  width: "100%",
  boxSizing: "border-box",
  padding: "10px 12px",
  borderRadius: 10,
  border: "1px solid #475569",
  background: "#0f172a",
  color: "#e2e8f0",
};

const buttonStyle: CSSProperties = {
  padding: "10px 16px",
  borderRadius: 10,
  border: "none",
  background: "#0ea5e9",
  color: "#0f172a",
  fontWeight: 700,
  cursor: "pointer",
};

const thStyle: CSSProperties = {
  padding: "10px 8px",
  borderBottom: "1px solid #334155",
  fontWeight: 600,
};

const tdStyle: CSSProperties = {
  padding: "12px 8px",
  borderBottom: "1px solid #1f2937",
};
