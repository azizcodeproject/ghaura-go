import type { Metadata } from "next";
import type { ReactNode } from "react";
import { AppNav } from "@/components/AppNav";
import { Providers } from "@/components/Providers";

export const metadata: Metadata = {
  title: "Ghaura — Ekspedisi",
  description: "Platform ekspedisi B2B kargo",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="id">
      <body style={{ margin: 0, fontFamily: "system-ui, sans-serif", background: "#0f172a", color: "#e2e8f0" }}>
        <Providers>
          <AppNav />
          {children}
        </Providers>
      </body>
    </html>
  );
}
