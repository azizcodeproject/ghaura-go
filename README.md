# ghaura-go

Platform ekspedisi B2B kargo (armada milik + partner).

## Stack

| Layer | Tech |
|-------|------|
| API | Go + Gin |
| DB | PostgreSQL 16 |
| Cache | Redis 7 |
| Broker | Redpanda (Kafka-compatible) |
| Web | Next.js + TanStack Query |

## Struktur

```
apps/api   — backend Gin
apps/web   — frontend Next.js
deploy/    — ops scripts
docs/      — arsitektur & domain
```

## Quick start (Docker)

```bash
docker compose up --build
```

- API: http://localhost:8080/health
- Web: http://localhost:3000 (daftar resi: http://localhost:3000/shipments)
- Postgres: `localhost:5432` (ghaura/ghaura)
- Redis: `localhost:6379`
- Kafka (Redpanda): `localhost:19092` (belum dipakai Sprint 1)

Panduan Sprint 1 (migrasi, cache Redis, contoh curl create customer → create shipment → get by resi): [docs/sprint-1.md](docs/sprint-1.md)

## MVP roadmap

1. Order + resi (Customer → Shipment)
2. Tracking + Assignment (driver/vehicle, ownership_type)
3. Driver app / POD
4. KPI driver + aset

## Local API tanpa Docker image

```bash
# infra only
docker compose up -d postgres redis redpanda

cd apps/api
go run ./cmd/server
```
