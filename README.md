# Project Management API
---

## Daftar Isi

- [Tech Stack](#tech-stack)
- [Arsitektur](#arsitektur)
- [Struktur Folder](#struktur-folder)
- [Instalasi](#instalasi)
- [Konfigurasi Environment](#konfigurasi-environment)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Skema Database](#skema-database)
- [Dokumentasi API](#dokumentasi-api)
- [Status Implementasi](#status-implementasi)
- [Deployment](#deployment)


---

## Tech Stack

| Layer | Teknologi |
|---|---|
| Bahasa | Go 1.25 |
| HTTP Framework | [Fiber v3](https://github.com/gofiber/fiber) |
| ORM | [GORM](https://gorm.io) (driver PostgreSQL) |
| Database | PostgreSQL |
| Auth | JWT (`golang-jwt/jwt/v5`), custom middleware |
| Hashing Password | bcrypt (`golang.org/x/crypto`) |
| Migration | [golang-migrate](https://github.com/golang-migrate/migrate) (file SQL manual di `database/migrations`) |
| Env Loader | `godotenv` |
| Struct Mapping | `jinzhu/copier` (mapping `User` → `UserResponse`) |

---

## Arsitektur

Layered architecture — konsepnya sama seperti kamu memisahkan `components/` dari `services/api/` di Next.js, supaya route handler tidak berisi semua logic sekaligus.

```
HTTP Request
   │
   ▼
routes/          → routing & pasang middleware (mirip App Router route.ts, tapi manual)
   │
   ▼
controllers/     → terima request, parsing query/body, bentuk response (≈ handler di route.ts)
   │
   ▼
services/        → business logic (≈ service layer di /lib atau /services Next.js kamu)
   │
   ▼
repositories/    → akses database via GORM, termasuk logic query dinamis (≈ Prisma client call yang dibungkus)
   │
   ▼
models/          → definisi struct/tabel (≈ schema.prisma)
   │
   ▼
PostgreSQL
```

```go
userRepo := repositories.NewUserRepository()
userService := services.NewUserService(userRepo)
userController := controllers.NewUserController(userService)
routes.Setup(app, userController)
```
## Struktur Folder

```
Project-Management/
├── config/
│   └── config.go              
├── controllers/
│   └── user_controller.go     
├── database/
│   ├── migrations/
│   │   ├── 000001_create_user_table.up.sql
│   │   └── 000001_create_user_table.down.sql
│   └── seed/
│       └── seed_admin.go      
├── middleware/
│   └── jwt.go                 
├── models/
│   ├── user.go                  
│   ├── board.go                 
│   ├── board_member.go          
│   ├── list.go                  
│   ├── list_position.go         
│   ├── card.go                  
│   ├── card_assignee.go         
│   ├── card_attachment.go       
│   ├── card_label.go            
│   ├── card_position.go         
│   ├── comment.go               
│   ├── label.go                 
│   └── types/
│       └── uuid_array.go       
├── repositories/
│   └── user_repository.go      
├── routes/
│   └── route.go                
├── services/
│   └── user_service.go         
├── utils/
│   ├── jwt.go                  
│   ├── password.go             
│   └── response.go             
├── .gitignore                  │
├── go.mod
├── go.sum
└── main.go                     
```

---

## Instalasi

### Prasyarat

- Go **1.25** atau lebih baru
- PostgreSQL **13+**
- [`golang-migrate` CLI](https://github.com/golang-migrate/migrate#cli-usage) (opsional, untuk menjalankan file migration SQL)

### Clone & install dependency

```bash
git clone https://github.com/Ciptaaaa/Project-Management.git
cd Project-Management
go mod download
```


# Server
PORT=3030

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_db_password
DB_NAME=project_management

# JWT
JWT_SECRET=ganti-dengan-random-string-panjang-dan-unik
JWT_EXPIRY=1h
REFRESH_TOKEN_EXPIRED=24h

# Seed admin
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123
ADMIN_ROLE=admin
```

## Menjalankan Aplikasi

### 1. Buat database

```bash
createdb project_management
```

### 2. Jalankan migration

```bash
migrate -path database/migrations \
  -database "postgres://postgres:your_db_password@localhost:5432/project_management?sslmode=disable" \
  up
```

### 3. Jalankan server

```bash
go run main.go
```
### Build binary untuk production

```bash
go build -o bin/server main.go
./bin/server
```

---

## Skema Database

Migration yang **sudah ada** (`000001_create_user_table.up.sql`):

```sql
CREATE TABLE users (
    internal_id BIGSERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password text NOT NULL,
    role varchar(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    public_id UUID NOT NULL DEFAULT gen_random_uuid(),
    CONSTRAINT user_public_id_unique UNIQUE(public_id)
);
```


## Dokumentasi API

Base URL: `http://localhost:3030`

Format response standar (`utils/response.go`):

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "...",
  "data": { }
}
```

Response error:
```json
{
  "status": "Error Bad Request",
  "response_code": 400,
  "message": "...",
  "error": "detail error"
}
```

### Auth

#### `POST /v1/auth/register`

Registrasi user baru. Role otomatis `"user"` — tidak bisa diset dari client (mencegah privilege escalation lewat body request).

**Request body**
```json
{ "name": "Cipta", "email": "cipta@example.com", "password": "secret123" }
```

**Response `200 OK`**
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Register Success!",
  "data": {
    "public_id": "b3f1...uuid",
    "name": "Cipta",
    "email": "cipta@example.com",
    "role": "user",
    "created_at": "2026-07-06T10:00:00Z",
    "updated_at": "2026-07-06T10:00:00Z"
  }
}
```

Error `400` jika email sudah terdaftar.

#### `POST /v1/auth/login`

**Request body**
```json
{ "email": "cipta@example.com", "password": "secret123" }
```

**Response `200 OK`**
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Login Successfully!",
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "user": {
      "public_id": "b3f1...uuid",
      "name": "Cipta",
      "email": "cipta@example.com",
      "role": "user"
    }
  }
}
```

`access_token` berisi claim `user_id`, `role`, `public_id`, `email`, `exp` (default 6 jam dari `JWT_EXPIRY`). `refresh_token` berisi `user_id` + `exp` (default 24 jam) — endpoint `POST /v1/auth/refresh` untuk redeem token ini 

### User (protected — butuh `Authorization: Bearer <access_token>`)

#### `GET /api/v1/users/page` — 🆕 pagination, filter, sort

Endpoint baru untuk list user dengan pagination.

**Query params**

| Param | Tipe | Default | Keterangan |
|---|---|---|---|
| `page` | int | `1` | Halaman ke berapa |
| `limit` | int | `10` | Jumlah item per halaman, di-cap maksimal `100` |
| `filter` | string | `""` | Cari di kolom `name` ATAU `email` (case-insensitive, `ILIKE`) |
| `sort` | string | `""` | Nama kolom untuk sorting. Prefix `-` = descending. Kolom yang diizinkan: `id`, `name`, `email` (whitelist di `allowedSortFields`) |

**Contoh request**
```
GET /api/v1/users/page?page=1&limit=10&filter=cipta&sort=-id
```

**Response `200 OK`**
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Data found",
  "data": [
    {
      "public_id": "b3f1...uuid",
      "name": "Cipta",
      "email": "cipta@example.com",
      "role": "user",
      "created_at": "...",
      "updated_at": "..."
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_page": 1,
    "filter": "cipta",
    "sort": "-id"
  }
}
```


#### `GET /api/v1/users/:id`

`:id` adalah `public_id` (UUID), bukan `internal_id`.

**Response `200 OK`**
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Data Found!",
  "data": {
    "public_id": "b3f1...uuid",
    "name": "Cipta",
    "email": "cipta@example.com",
    "role": "user",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

## Status Implementasi

| Modul | Model | Migration SQL | Repository | Service | Controller | Route |
|---|:---:|:---:|:---:|:---:|:---:|:---:|
| User / Auth | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| User — Pagination/Filter/Sort | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Board | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| List | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Card | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Label / Comment / Attachment | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Refresh token endpoint | – | – | – | ❌ | ❌ | ❌ |

