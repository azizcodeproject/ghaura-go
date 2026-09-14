"use client";

import { useParams } from "next/navigation";
import { ShipmentDetail } from "@/components/ShipmentDetail";

export default function ShipmentByResiPage() {
  const params = useParams<{ resiNumber: string }>();
  const resiNumber = decodeURIComponent(params.resiNumber ?? "");

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: 32 }}>
      <ShipmentDetail resiNumber={resiNumber} />
    </main>
  );
}
