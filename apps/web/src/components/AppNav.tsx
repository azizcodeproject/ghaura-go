import Link from "next/link";

export function AppNav() {
  return (
    <header
      style={{
        borderBottom: "1px solid #334155",
        background: "#0b1220",
      }}
    >
      <div
        style={{
          maxWidth: 1080,
          margin: "0 auto",
          padding: "16px 24px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: 16,
        }}
      >
        <Link href="/" style={{ color: "#e2e8f0", textDecoration: "none", fontWeight: 700 }}>
          Ghaura
        </Link>
        <nav style={{ display: "flex", gap: 16, fontSize: 14 }}>
          <Link href="/shipments" style={{ color: "#94a3b8", textDecoration: "none" }}>
            Resi pengiriman
          </Link>
        </nav>
      </div>
    </header>
  );
}
