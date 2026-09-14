"use client";

import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState, type CSSProperties, type FormEvent } from "react";
import { StatusBadge } from "@/components/StatusBadge";
import { getShipmentById, getShipmentByResi, updateShipmentStatus } from "@/lib/api";
import { formatDateTime } from "@/lib/format";
import { SHIPMENT_STATUSES, STATUS_LABELS, type ShipmentStatus } from "@/lib/types";

type ShipmentDetailProps = {
  shipmentId?: string;
  resiNumber?: string;
};

export function ShipmentDetail({ shipmentId, resiNumber }: ShipmentDetailProps) {
  const queryClient = useQueryClient();
  const [selectedStatus, setSelectedStatus] = useState<ShipmentStatus | "">("");

  const shipmentQuery = useQuery({
    queryKey: shipmentId ? ["shipment", shipmentId] : ["shipment-resi", resiNumber],
    queryFn: () => (shipmentId ? getShipmentById(shipmentId) : getShipmentByResi(resiNumber ?? "")),
    enabled: Boolean(shipmentId || resiNumber),
  });

  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: ShipmentStatus }) => updateShipmentStatus(id, status),
    onSuccess: (updated) => {
      queryClient.setQueryData(["shipment", updated.id], updated);
      queryClient.invalidateQueries({ queryKey: ["shipments"] });
      queryClient.invalidateQueries({ queryKey: ["shipment-resi", updated.resiNumber] });
      setSelectedStatus("");
    },
  });

  if (shipmentQuery.isLoading) {
    return <p style={{ color: "#94a3b8" }}>Memuat detail resi…</p>;
  }

  if (shipmentQuery.isError || !shipmentQuery.data) {
    return (
      <div>
        <p style={{ color: "#f87171", marginBottom: 16 }}>
          {shipmentQuery.error instanceof Error ? shipmentQuery.error.message : "Resi tidak ditemukan"}
        </p>
        <Link href="/shipments" style={{ color: "#38bdf8" }}>
          Kembali ke daftar resi
        </Link>
      </div>
    );
  }

  const shipment = shipmentQuery.data;
  const nextStatus = selectedStatus || shipment.status;

  return (
    <article style={cardStyle}>
      <div style={{ display: "flex", justifyContent: "space-between", gap: 16, flexWrap: "wrap" }}>
        <div>
          <p style={{ color: "#94a3b8", margin: 0, fontSize: 13 }}>Nomor resi</p>
          <h1 style={{ margin: "6px 0 0", fontSize: 28 }}>{shipment.resiNumber}</h1>
        </div>
        <StatusBadge status={shipment.status} />
      </div>

      {shipment.cache && (
        <p style={{ color: shipment.cache === "HIT" ? "#34d399" : "#94a3b8", fontSize: 13 }}>
          {shipment.cache === "HIT" ? "Data diambil dari cache Redis" : "Data diambil dari database"}
        </p>
      )}

      <dl style={detailGridStyle}>
        <DetailItem label="Pelanggan" value={shipment.customerName || shipment.customerId} />
        <DetailItem label="Asal" value={shipment.origin} />
        <DetailItem label="Tujuan" value={shipment.destination} />
        <DetailItem label="Dibuat" value={formatDateTime(shipment.createdAt)} />
        <DetailItem label="Terakhir diubah" value={formatDateTime(shipment.updatedAt)} />
      </dl>

      <form
        onSubmit={(event: FormEvent<HTMLFormElement>) => {
          event.preventDefault();
          statusMutation.mutate({ id: shipment.id, status: nextStatus });
        }}
        style={{ display: "flex", gap: 12, flexWrap: "wrap", alignItems: "end" }}
      >
        <label style={{ display: "grid", gap: 6, fontSize: 13, color: "#94a3b8" }}>
          Ubah status
          <select
            value={nextStatus}
            onChange={(event) => setSelectedStatus(event.target.value as ShipmentStatus)}
            style={inputStyle}
          >
            {SHIPMENT_STATUSES.map((status) => (
              <option key={status} value={status}>
                {STATUS_LABELS[status]}
              </option>
            ))}
          </select>
        </label>
        <button type="submit" disabled={statusMutation.isPending || nextStatus === shipment.status} style={buttonStyle}>
          {statusMutation.isPending ? "Menyimpan…" : "Simpan status"}
        </button>
      </form>

      {statusMutation.isError && (
        <p style={{ color: "#f87171" }}>
          {statusMutation.error instanceof Error ? statusMutation.error.message : "Gagal mengubah status"}
        </p>
      )}
      {statusMutation.isSuccess && <p style={{ color: "#34d399" }}>Status pengiriman sudah diperbarui.</p>}

      <Link href="/shipments" style={{ color: "#38bdf8", display: "inline-block", marginTop: 8 }}>
        Kembali ke daftar resi
      </Link>
    </article>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt style={{ color: "#94a3b8", fontSize: 13, marginBottom: 4 }}>{label}</dt>
      <dd style={{ margin: 0 }}>{value}</dd>
    </div>
  );
}

const cardStyle: CSSProperties = {
  padding: 24,
  borderRadius: 16,
  border: "1px solid #334155",
  background: "#1e293b",
  display: "grid",
  gap: 20,
};

const detailGridStyle: CSSProperties = {
  display: "grid",
  gridTemplateColumns: "repeat(auto-fit, minmax(180px, 1fr))",
  gap: 16,
  margin: 0,
};

const inputStyle: CSSProperties = {
  minWidth: 220,
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
