"use client";

import { useParams } from "next/navigation";
import { ShipmentDetail } from "@/components/ShipmentDetail";

export default function ShipmentByIdPage() {
  const params = useParams<{ id: string }>();

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: 32 }}>
      <ShipmentDetail shipmentId={params.id} />
    </main>
  );
}
