"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { apiHealth } from "@/lib/api";

export default function HomePage() {
  const health = useQuery({ queryKey: ["health"], queryFn: apiHealth, refetchInterval: 15000 });

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: 32 }}>
      <h1 style={{ fontSize: 28, marginBottom: 8 }}>Ghaura</h1>
      <p style={{ color: "#94a3b8", marginBottom: 24 }}>
        Platform ekspedisi B2B kargo — kelola pelanggan, terbitkan resi, dan pantau status pengiriman.
      </p>

      <Link
        href="/shipments"
        style={{
          display: "inline-block",
          marginBottom: 24,
          padding: "10px 16px",
          borderRadius: 10,
          background: "#0ea5e9",
          color: "#0f172a",
          fontWeight: 700,
          textDecoration: "none",
        }}
      >
        Buka daftar resi
      </Link>

      <div
        style={{
          padding: 16,
          borderRadius: 12,
          border: "1px solid #334155",
          background: "#1e293b",
        }}
      >
        <div style={{ fontSize: 13, color: "#94a3b8", marginBottom: 8 }}>API /health</div>
        {health.isLoading && <div>Memeriksa koneksi API…</div>}
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
