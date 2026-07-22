# Project Management API

REST API backend untuk aplikasi manajemen proyek bergaya Trello (board → list → card). Dibangun dengan Go, Fiber v3, dan GORM.

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
| Migration | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Env Loader | `godotenv` |
| Struct Mapping | `jinzhu/copier` |

---

## Arsitektur

Layered architecture:

```
routes/          → routing & middleware
controllers/     → parsing request, membentuk response
services/        → business logic
repositories/    → akses database via GORM
models/          → definisi struct/tabel
```

Dependency injection dirakit manual di `main.go`:

```go
userRepo := repositories.NewUserRepository()
userService := services.NewUserService(userRepo)
userController := controllers.NewUserController(userService)

boardRepo := repositories.NewBoardRepository()
boardMemberRepo := repositories.NewBoardMemberRepository()
boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
boardController := controllers.NewBoardController(boardService)

listPosRepo := repositories.NewListPositionRepository()
listRepo := repositories.NewListRepository()
listService := services.NewListService(listRepo, boardRepo, listPosRepo)
listController := controllers.NewListController(listService)

routes.Setup(app, userController, boardController, listController)
```

---

## Struktur Folder

```
Project-Management/
├── config/
│   └── config.go
├── controllers/
│   ├── user_controller.go
│   ├── board_controller.go
│   └── list_controller.go
├── database/
│   ├── migrations/
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
│   ├── user_repository.go
│   ├── board_repository.go
│   ├── board_member_repository.go
│   ├── list_repository.go
│   └── list_position_repository.go
├── routes/
│   └── route.go
├── services/
│   ├── user_service.go
│   ├── board_service.go
│   └── list_service.go
├── utils/
│   ├── jwt.go
│   ├── password.go
│   ├── response.go
│   └── sorting_list_position.go
├── go.mod
├── go.sum
└── main.go
```

---

## Instalasi

### Prasyarat

- Go 1.25 atau lebih baru
- PostgreSQL 13+
- [`golang-migrate` CLI](https://github.com/golang-migrate/migrate#cli-usage)

### Clone & install dependency

```bash
git clone https://github.com/Ciptaaaa/Project-Management.git
cd Project-Management
go mod download
```

---

## Konfigurasi Environment

Buat file `.env` di root project:

```bash
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
JWT_EXPIRY=6h
REFRESH_TOKEN_EXPIRED=24h
```

---

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

### 3. Jalankan server (development)

```bash
go run main.go
```

### 4. Build binary untuk production

```bash
go build -o bin/server main.go
./bin/server
```

---

## Skema Database

```
users
 ├─ internal_id (PK, bigserial)
 ├─ public_id   (UUID, unique)
 ├─ name, email (unique), password (bcrypt hash), role
 └─ created_at, updated_at, deleted_at

boards
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ title, description, due_date
 └─ owner_internal_id / owner_public_id → FK users, ON DELETE CASCADE

board_members
 ├─ board_internal_id → FK boards, ON DELETE CASCADE
 └─ user_internal_id  → FK users,  ON DELETE CASCADE

lists
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ board_internal_id / board_public_id → FK boards, ON DELETE CASCADE
 └─ title, position

list_positions
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ board_internal_id → FK boards, ON DELETE CASCADE (unique per board)
 └─ list_order UUID[]
```

Model `Card`, `CardAssignee`, `CardAttachment`, `CardLabel`, `CardPosition`, `Comment`, `Label` sudah ada di `models/`, tapi belum ada migration untuk tabelnya.

---

## Dokumentasi API

Base URL (default): `http://localhost:3030`

### Format Response

**Success**
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "...",
  "data": { }
}
```

**Success dengan pagination**
```json
{
  "status": "Success",
  "response_code": 200,
  "message": "...",
  "data": [ ],
  "meta": { "page": 1, "limit": 10, "total": 100, "total_page": 10, "filter": "", "sort": "" }
}
```

**Error**
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

`access_token` membawa claim `user_id`, `role`, `public_id`, `email`, `exp` (default `JWT_EXPIRY`). `refresh_token` membawa `user_id` + `exp` (default `REFRESH_TOKEN_EXPIRED`).

Semua endpoint di bawah ini butuh header `Authorization: Bearer <access_token>`.

### User

#### `GET /api/v1/users/page`

| Param | Tipe | Default | Keterangan |
|---|---|---|---|
| `page` | int | `1` | Halaman ke berapa |
| `limit` | int | `10` | Item per halaman, di-cap maksimal `100` |
| `filter` | string | `""` | Cari di kolom `name` ATAU `email` (case-insensitive) |
| `sort` | string | `""` | Nama kolom untuk sorting. Prefix `-` = descending |

```
GET /api/v1/users/page?page=1&limit=10&filter=cipta&sort=-internal_id
```

#### `GET /api/v1/users/:id`

`:id` = `public_id`.

#### `PUT /api/v1/users/:id`

Update data user. `:id` = `public_id`.

#### `DELETE /api/v1/users/:id`

Soft delete. `:id` = `internal_id`.

### Board (protected)

#### `POST /api/v1/boards`

**Request body**
```json
{ "title": "Website Redesign", "description": "...", "due_date": "2026-08-01T00:00:00Z" }
```

#### `PUT /api/v1/boards/:id`

Update board. `:id` = `public_id`.

#### `POST /api/v1/boards/:id/members`

**Request body**
```json
["b3f1...uuid", "a2e2...uuid"]
```

#### `DELETE /api/v1/boards/:id/members`

**Request body**
```json
["b3f1...uuid", "a2e2...uuid"]
```

#### `GET /api/v1/boards/my`

Query params: `page`, `limit`, `filter`, `sort`.

### List (protected)

#### `POST /api/v1/lists`

**Request body**
```json
{ "board_public_id": "b3f1...uuid", "title": "To Do" }
```

#### `PUT /api/v1/lists/:id`

Update list. `:id` = `public_id`.

---

## Status Implementasi

| Modul | Model | Migration SQL | Repository | Service | Controller | Route |
|---|:---:|:---:|:---:|:---:|:---:|:---:|
| User / Auth | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| User — Pagination/Filter/Sort | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Board | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Board Member | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| List | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| List Position (reorder) | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Card | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Card Assignee / Position | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Label / Comment / Attachment | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Refresh token redeem | – | – | – | ❌ | ❌ | ❌ |

---

## Deployment

1. Build binary untuk target OS/arch:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bin/server main.go
   ```
2. Siapkan PostgreSQL (managed atau self-hosted) dan jalankan seluruh migration di `database/migrations`.
3. Set environment variable production lewat secret manager platform (Railway/Fly.io/Docker secrets/systemd EnvironmentFile).
4. Jalankan binary di belakang reverse proxy (Nginx/Caddy) untuk TLS termination, atau expose lewat platform yang sudah handle TLS (Railway, Fly.io, Render).