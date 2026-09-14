"use client";

import { useQuery } from "@tanstack/react-query";
import { apiHealth } from "@/lib/api";

export default function HomePage() {
  const health = useQuery({ queryKey: ["health"], queryFn: apiHealth, refetchInterval: 15000 });

  return (
    <main style={{ maxWidth: 720, margin: "0 auto", padding: 32 }}>
      <h1 style={{ fontSize: 28, marginBottom: 8 }}>Ghaura</h1>
      <p style={{ color: "#94a3b8", marginBottom: 24 }}>
        Ekspedisi B2B kargo — MVP scaffold (tracking → KPI driver → aset)
      </p>
      <div
        style={{
          padding: 16,
          borderRadius: 12,
          border: "1px solid #334155",
          background: "#1e293b",
        }}
      >
        <div style={{ fontSize: 13, color: "#94a3b8", marginBottom: 8 }}>API /health</div>
        {health.isLoading && <div>Checking…</div>}
        {health.isError && <div style={{ color: "#f87171" }}>API belum reachable</div>}
        {health.data && (
          <pre style={{ margin: 0, fontSize: 13, whiteSpace: "pre-wrap" }}>
            {JSON.stringify(health.data, null, 2)}
          </pre>
        )}
      </div>
    </main>
  );
}
