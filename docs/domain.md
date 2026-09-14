# Domain model (MVP)

```
Customer (perusahaan)
  └── Shipment / resi
        ├── Assignment (Driver + Vehicle, ownership OWNED|PARTNER)
        ├── TrackingEvent (lat/lng/time)
        └── POD (foto, geo, nama penerima)
```

- **Vehicle / Driver** membawa `ownership_type` supaya KPI & aset tidak tertukar milik vs partner.
- Kafka topics (rencana): `shipment.events`, `tracking.events`, `pod.events`
- Redis cache (rencana): resi lookup, last position per vehicle
