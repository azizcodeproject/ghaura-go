export type ShipmentStatus = "DRAFT" | "BOOKED" | "IN_TRANSIT" | "DELIVERED" | "CANCELLED";

export type Customer = {
  id: string;
  name: string;
  npwp?: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
};

export type Shipment = {
  id: string;
  resiNumber: string;
  customerId: string;
  customerName?: string;
  status: ShipmentStatus;
  origin: string;
  destination: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
  cache?: "HIT" | "MISS";
};

export type ListResponse<T> = {
  items: T[];
};

export const SHIPMENT_STATUSES: ShipmentStatus[] = [
  "DRAFT",
  "BOOKED",
  "IN_TRANSIT",
  "DELIVERED",
  "CANCELLED",
];

export const STATUS_LABELS: Record<ShipmentStatus, string> = {
  DRAFT: "Draf",
  BOOKED: "Dipesan",
  IN_TRANSIT: "Dalam perjalanan",
  DELIVERED: "Terkirim",
  CANCELLED: "Dibatalkan",
};
