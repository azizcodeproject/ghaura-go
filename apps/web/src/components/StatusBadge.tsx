import { STATUS_LABELS, type ShipmentStatus } from "@/lib/types";

const STATUS_COLORS: Record<ShipmentStatus, { background: string; color: string }> = {
  DRAFT: { background: "#334155", color: "#e2e8f0" },
  BOOKED: { background: "#1d4ed8", color: "#dbeafe" },
  IN_TRANSIT: { background: "#b45309", color: "#ffedd5" },
  DELIVERED: { background: "#047857", color: "#d1fae5" },
  CANCELLED: { background: "#be123c", color: "#ffe4e6" },
};

export function StatusBadge({ status }: { status: ShipmentStatus }) {
  const tone = STATUS_COLORS[status];
  return (
    <span
      style={{
        display: "inline-flex",
        alignItems: "center",
        padding: "4px 10px",
        borderRadius: 999,
        fontSize: 12,
        fontWeight: 600,
        background: tone.background,
        color: tone.color,
      }}
    >
      {STATUS_LABELS[status]}
    </span>
  );
}
