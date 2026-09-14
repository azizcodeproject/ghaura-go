const base = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function apiHealth() {
  const res = await fetch(`${base}/health`, { cache: "no-store" });
  if (!res.ok) throw new Error(`health ${res.status}`);
  return res.json();
}
