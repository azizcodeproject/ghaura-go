import type { Customer, ListResponse, Shipment, ShipmentStatus } from "./types";

const apiBaseUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    cache: "no-store",
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });

  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as { error?: string } | null;
    throw new Error(payload?.error ?? `Permintaan gagal (${response.status})`);
  }

  return (await response.json()) as T;
}

export async function apiHealth() {
  return request<{ status: string; service: string; time: string }>("/health");
}

export function listCustomers() {
  return request<ListResponse<Customer>>("/api/v1/customers");
}

export function createCustomer(input: { name: string; npwp?: string }) {
  return request<Customer>("/api/v1/customers", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function listShipments(status?: ShipmentStatus | "") {
  const query = status ? `?status=${encodeURIComponent(status)}` : "";
  return request<ListResponse<Shipment>>(`/api/v1/shipments${query}`);
}

export function createShipment(input: { customerId: string; origin: string; destination: string }) {
  return request<Shipment>("/api/v1/shipments", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function getShipmentById(shipmentId: string) {
  return request<Shipment>(`/api/v1/shipments/${encodeURIComponent(shipmentId)}`);
}

export function getShipmentByResi(resiNumber: string) {
  return request<Shipment>(`/api/v1/shipments/by-resi/${encodeURIComponent(resiNumber)}`);
}

export function updateShipmentStatus(shipmentId: string, status: ShipmentStatus) {
  return request<Shipment>(`/api/v1/shipments/${encodeURIComponent(shipmentId)}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
}
