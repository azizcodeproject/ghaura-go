# Sprint 1 — Order + resi

Sprint ini menambahkan master pelanggan, penerbitan resi, lookup resi (dengan cache Redis), dan halaman web untuk daftar/detail pengiriman.

## Jalankan lokal dengan Docker Compose

Dari root repository:

```bash
docker compose up --build
```

Layanan yang naik:

| Layanan | URL / port |
|---------|------------|
| API | http://localhost:8080/health |
| Web | http://localhost:3000 |
| Postgres | `localhost:5432` (user/pass/db: `ghaura`) |
| Redis | `localhost:6379` |
| Redpanda | `localhost:19092` (belum dipakai di Sprint 1) |

Migrasi `customers` dan `shipments` dijalankan otomatis saat API start.

### Infra saja + API native

```bash
docker compose up -d postgres redis
cd apps/api
go run ./cmd/server
```

Frontend:

```bash
cd apps/web
npm install
npm run dev
```

Set `NEXT_PUBLIC_API_URL=http://localhost:8080` bila perlu (sudah default).

## Contoh alur API (curl)

### 1. Buat pelanggan

```bash
curl -sS -X POST http://localhost:8080/api/v1/customers \
  -H 'Content-Type: application/json' \
  -d '{"name":"PT Maju Jaya","npwp":"10.0.0.1-000.000"}'
```

Simpan `id` dari response sebagai `CUSTOMER_ID`.

### 2. Buat pengiriman / resi

```bash
curl -sS -X POST http://localhost:8080/api/v1/shipments \
  -H 'Content-Type: application/json' \
  -d '{"customerId":"CUSTOMER_ID","origin":"Jakarta","destination":"Surabaya"}'
```

Response berisi `resiNumber` unik (contoh `GHR-260914-A3K7MP`) dan `id` pengiriman.

### 3. Ambil resi (cache Redis)

Request pertama biasanya `X-Cache: MISS` (lalu disimpan 5 menit di key `shipment:resi:{resiNumber}`):

```bash
curl -sS -D - http://localhost:8080/api/v1/shipments/by-resi/GHR-260914-A3K7MP
```

Request berikutnya dalam TTL harus `X-Cache: HIT`. Body juga punya field `cache`.

### Endpoint lain

```bash
# daftar pelanggan
curl -sS http://localhost:8080/api/v1/customers

# daftar resi, opsional filter status
curl -sS 'http://localhost:8080/api/v1/shipments?status=DRAFT'

# detail by id
curl -sS http://localhost:8080/api/v1/shipments/SHIPMENT_ID

# ubah status — cache resi di-invalidate
curl -sS -X PATCH http://localhost:8080/api/v1/shipments/SHIPMENT_ID/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"BOOKED"}'
```

Status yang valid: `DRAFT`, `BOOKED`, `IN_TRANSIT`, `DELIVERED`, `CANCELLED`.

## Frontend

- `/shipments` — form pelanggan, form resi, pencarian resi, dan tabel daftar
- `/shipments/[id]` — detail + ubah status
- `/shipments/resi/[resiNumber]` — detail via endpoint cache Redis
