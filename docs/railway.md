# Deploy Ghaura-Go di Railway

Sprint 1 di-deploy sebagai **dua service aplikasi** (bukan `docker-compose` sebagai satu unit), plus plugin Postgres dan Redis. Kafka / Redpanda **tidak diperlukan**.

`apps/api/railway.json` dan `apps/web/railway.json` mengunci builder ke Dockerfile yang sudah ada.

## 1. Buat project dari GitHub

1. Buka [Railway](https://railway.com) dan buat project baru.
2. Pilih **Deploy from GitHub repo** → `azizcodeproject/ghaura-go` (branch `main`).
3. Jangan biarkan Railway men-deploy seluruh monorepo sebagai satu service. Hapus service otomatis jika terlanjur dibuat, lalu lanjut langkah di bawah.

## 2. Tambah Postgres dan Redis

Di project yang sama, tambah plugin:

- **Postgres** (PostgreSQL)
- **Redis**

Tidak perlu menambah Kafka / Redpanda untuk Sprint 1.

## 3. Service `api` — root `apps/api`

Tambah service dari repo yang sama, lalu set:

| Setting | Nilai |
|---------|--------|
| Root Directory | `apps/api` |
| Builder | Dockerfile (`apps/api/Dockerfile`) |

Railway akan meng-inject `PORT`. API memakai urutan: `PORT` → `HTTP_ADDR` → `:8080`.

### Variabel service `api`

| Variabel | Wajib | Sumber / catatan |
|----------|-------|------------------|
| `DATABASE_URL` | Ya | Referensi plugin Postgres, misalnya `${{Postgres.DATABASE_URL}}` |
| `REDIS_URL` | Ya (salah satu) | Referensi plugin Redis, misalnya `${{Redis.REDIS_URL}}` (`redis://` atau `rediss://`) |
| `REDIS_ADDR` | Alternatif | `host:port` jika tidak memakai `REDIS_URL`. `REDIS_URL` menang jika keduanya terisi |
| `APP_ENV` | Disarankan | `production` |
| `HTTP_ADDR` | Tidak | Hanya fallback lokal; di Railway `PORT` sudah cukup |
| `KAFKA_BROKERS` | Tidak | Opsional / skip Sprint 1 |

Nama service plugin di interpolasi harus sama persis dengan nama di dashboard (sering `Postgres` dan `Redis`).

**Migrasi berjalan otomatis saat API start** (`customers`, `shipments`). Tidak ada job migrasi terpisah.

Generate **public domain** untuk `api` (Settings → Networking → Generate Domain). Simpan URL publik, misalnya `https://ghaura-api.up.railway.app`.

## 4. Service `web` — root `apps/web`

Tambah service kedua dari repo yang sama:

| Setting | Nilai |
|---------|--------|
| Root Directory | `apps/web` |
| Builder | Dockerfile (`apps/web/Dockerfile`) |

`NEXT_PUBLIC_API_URL` di-inline saat **build**. Set sebelum build/deploy web, lalu redeploy web jika URL API berubah.

### Variabel service `web`

| Variabel | Wajib | Catatan |
|----------|-------|---------|
| `NEXT_PUBLIC_API_URL` | Ya | URL publik API, contoh `https://ghaura-api.up.railway.app` (tanpa slash di akhir) |

Railway meng-inject `PORT` untuk proses Next.js standalone. Dockerfile web bind ke `0.0.0.0`.

Generate **public domain** untuk `web`. Browser memanggil API lewat `NEXT_PUBLIC_API_URL`, jadi domain API harus publik.

## 5. Urutan yang aman

1. Postgres + Redis siap.
2. Deploy `api` dengan `DATABASE_URL` dan `REDIS_URL`.
3. Cek `https://<api-domain>/health`.
4. Set `NEXT_PUBLIC_API_URL` di `web` ke URL publik API.
5. Deploy / redeploy `web`.
6. Buka domain web → `/shipments`.

## 6. Kafka

`KAFKA_BROKERS` masih terbaca di config API tetapi **belum dipakai** di Sprint 1. Jangan provision Kafka di Railway untuk sekarang.

## 7. Lokal tetap Docker Compose

Pengembangan lokal tidak berubah:

```bash
docker compose up --build
```

Lihat [README](../README.md) dan [sprint-1.md](sprint-1.md). Compose tetap memakai `HTTP_ADDR` dan `REDIS_ADDR`.
