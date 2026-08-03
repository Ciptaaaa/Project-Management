# Project Management API

REST API backend for a Trello-style project management app (board → list → card). Built with Go, Fiber v3, and GORM.

## Table of Contents

- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Folder Structure](#folder-structure)
- [Installation](#installation)
- [Environment Configuration](#environment-configuration)
- [Running the App](#running-the-app)
- [Postman Collection](#postman-collection)
- [Database Schema](#database-schema)
- [API Documentation](#api-documentation)
- [Implementation Status](#implementation-status)
- [Deployment](#deployment)

---

## Tech Stack

| Layer            | Technology                                                  |
| ---------------- | ----------------------------------------------------------- |
| Language         | Go 1.25                                                     |
| HTTP Framework   | [Fiber v3](https://github.com/gofiber/fiber)                |
| ORM              | [GORM](https://gorm.io) (PostgreSQL driver)                 |
| Database         | PostgreSQL                                                  |
| Auth             | JWT (`golang-jwt/jwt/v5`), custom middleware                |
| Password Hashing | bcrypt (`golang.org/x/crypto`)                              |
| Migration        | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Env Loader       | `godotenv`                                                  |
| Struct Mapping   | `jinzhu/copier`                                             |

---

## Architecture

Layered architecture:

```
routes/          → routing & middleware
controllers/     → request parsing, response shaping
services/        → business logic
repositories/    → database access via GORM
models/          → struct/table definitions
```

Dependency injection is wired manually in `main.go`:

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

labelRepo := repositories.NewLabelRepository()
labelService := services.NewLabelService(labelRepo)
labelController := controllers.NewLabelController(labelService)

cardRepo := repositories.NewCardRepository()
cardService := services.NewCardService(cardRepo, listRepo, userRepo, labelRepo)
cardController := controllers.NewCardController(cardService)

routes.Setup(app, userController, boardController, listController, cardController, labelController)
```

---

## Folder Structure

```
Project-Management/
├── config/
│   └── config.go
├── controllers/
│   ├── user_controller.go
│   ├── board_controller.go
│   ├── list_controller.go
│   ├── card_controller.go
│   └── label_controller.go
├── database/
│   ├── migrations/
│   └── seed/
│       └── seed_admin.go
├── docs/
│   └── postman/
│       ├── Project-Management-API.postman_collection.json
│       └── Local-Dev.postman_environment.json
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
│   ├── list_position_repository.go
│   ├── card_repository.go
│   ├── label_repository.go
│   └── attachment_repository.go
├── routes/
│   └── route.go
├── services/
│   ├── user_service.go
│   ├── board_service.go
│   ├── list_service.go
│   ├── card_service.go
│   ├── label_service.go
│   └── attachment_service.go
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

## Installation

### Prerequisites

- Go 1.25 or newer
- PostgreSQL 13+
- [`golang-migrate` CLI](https://github.com/golang-migrate/migrate#cli-usage)

### Clone & install dependencies

```bash
git clone https://github.com/Ciptaaaa/Project-Management.git
cd Project-Management
go mod download
```

---

## Environment Configuration

Create a `.env` file at the project root:

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
JWT_SECRET=replace-with-a-long-random-unique-string
JWT_EXPIRY=6h
REFRESH_TOKEN_EXPIRED=24h
```

---

## Running the App

### 1. Create the database

```bash
createdb project_management
```

### 2. Run migrations

```bash
migrate -path database/migrations \
  -database "postgres://postgres:your_db_password@localhost:5432/project_management?sslmode=disable" \
  up
```

### 3. Run the server (development)

```bash
go run main.go
```

### 4. Build a production binary

```bash
go build -o bin/server main.go
./bin/server
```

---

## Postman Collection

A ready-to-use Postman collection and environment are provided under `docs/postman/`:

- `docs/postman/Project-Management-API.postman_collection.json` — all endpoints, grouped by resource (Auth, User, Board, List, Card, Label)
- `docs/postman/Local-Dev.postman_environment.json` — `base_url` and `access_token` variables

### Setup

1. Open Postman → **Import** → drag & drop both files (or **File → Import**)
2. Select the **Local Dev** environment from the dropdown in the top-right corner
3. Run **Auth → Login** — the response's `access_token` is captured automatically into the environment via a post-response script, so every other request picks it up through the `{{access_token}}` variable
4. All request URLs use `{{base_url}}` — update it in the environment to point at staging/production instead of editing every request individually

---

## Database Schema

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

cards
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ list_internal_id → FK lists, ON DELETE CASCADE
 └─ title, description, due_date, position

comments
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ card_internal_id → FK cards, ON DELETE CASCADE
 ├─ user_internal_id → FK users, ON DELETE CASCADE
 └─ message

labels
 ├─ internal_id (PK), public_id (UUID, unique)
 └─ name, color

card_labels (pivot)
 ├─ card_internal_id  → FK cards,  ON DELETE CASCADE
 └─ label_internal_id → FK labels, ON DELETE CASCADE

card_assignees (pivot)
 ├─ card_internal_id → FK cards, ON DELETE CASCADE
 └─ user_internal_id → FK users, ON DELETE CASCADE

card_attachment
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ card_internal_id → FK cards, ON DELETE CASCADE
 ├─ user_internal_id → FK users, ON DELETE CASCADE
 └─ file

card_positions
 ├─ internal_id (PK), public_id (UUID, unique)
 ├─ list_internal_id → FK lists, ON DELETE CASCADE (unique per list)
 └─ card_order UUID[]
```

---

## API Documentation

Base URL (default): `http://localhost:3030`

### Response Format

**Success**

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "...",
  "data": {}
}
```

**Success with pagination**

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "...",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_page": 10,
    "filter": "",
    "sort": ""
  }
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

`access_token` carries the claims `user_id`, `role`, `public_id`, `email`, `exp` (defaults to `JWT_EXPIRY`). `refresh_token` carries `user_id` + `exp` (defaults to `REFRESH_TOKEN_EXPIRED`).

All endpoints below require the header `Authorization: Bearer <access_token>`.

### User

#### `GET /api/v1/users/page`

| Param    | Type   | Default | Description                                               |
| -------- | ------ | ------- | --------------------------------------------------------- |
| `page`   | int    | `1`     | Page number                                               |
| `limit`  | int    | `10`    | Items per page, capped at `100`                           |
| `filter` | string | `""`    | Searches the `name` OR `email` columns (case-insensitive) |
| `sort`   | string | `""`    | Column name to sort by. Prefix `-` for descending         |

```
GET /api/v1/users/page?page=1&limit=10&filter=cipta&sort=-internal_id
```

#### `GET /api/v1/users/:id`

`:id` = `public_id`.

#### `PUT /api/v1/users/:id`

Update user data. `:id` = `public_id`.

#### `DELETE /api/v1/users/:id`

Soft delete. `:id` = `internal_id`.

### Board (protected)

#### `POST /api/v1/boards`

**Request body**

```json
{
  "title": "Website Redesign",
  "description": "...",
  "due_date": "2026-08-01T00:00:00Z"
}
```

`owner_public_id` is derived automatically from the JWT claim (`public_id`), not from the request body.

#### `PUT /api/v1/boards/:id`

Update a board. `:id` = `public_id`.

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

List of boards owned by the currently authenticated user, paginated.

Query params: `page`, `limit`, `filter`, `sort`.

#### `GET /api/v1/boards/:board_id/lists`

Get all lists belonging to a board. `:board_id` = board `public_id`.

#### `PUT /api/v1/boards/:board_id/position`

Update the order of lists within a board (drag & drop reorder).

**Request body**

```json
["list-public-id-1", "list-public-id-2", "list-public-id-3"]
```

### List (protected)

#### `POST /api/v1/lists`

**Request body**

```json
{ "board_public_id": "b3f1...uuid", "title": "To Do" }
```

#### `PUT /api/v1/lists/:id`

Update a list. `:id` = `public_id`.

#### `DELETE /api/v1/lists/:id`

Delete a list. `:id` = `public_id`.

#### `GET /api/v1/lists/:list_id/cards`

Get all cards belonging to a list. `:list_id` = list `public_id`.

#### `PUT /api/v1/lists/:list_id/positions`

Reorder cards within a single list (drag & drop, Trello-style). `:list_id` = list `public_id`. The array must contain the `public_id` of **every** card currently in the list, in the new order — the stored order is fully replaced, not merged.

**Request body**

```json
{
  "positions": ["card-public-id-1", "card-public-id-2"]
}
```

**Response `200 OK`**

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Successfully updated card position",
  "data": null
}
```

### Card (protected)

#### `POST /api/v1/cards`

**Request body**

```json
{
  "list_id": "b3f1...uuid",
  "title": "Setup CI/CD",
  "description": "...",
  "due_date": "2026-08-01T00:00:00Z",
  "position": 0
}
```

`list_id` in the body is the target list's `public_id`. On creation, the card is automatically appended to the target list's `card_positions` order.

#### `PUT /api/v1/cards/:id`

Update a card. `:id` = `public_id`. Same body as create; `list_id` is optional when moving the card to a different list. Moving a card to a different list automatically removes it from the old list's position order and appends it to the new one.

#### `DELETE /api/v1/cards/:id`

Delete a card. `:id` = `public_id`.

#### `GET /api/v1/cards/:id`

Card detail, including the `assignees`, `attachments`, and `labels` relations (when present).

#### `POST /api/v1/cards/:id/labels`

Attach a label to a card. `:id` = card `public_id`.

**Request body**

```json
{ "label_id": "b3f1...uuid" }
```

#### `DELETE /api/v1/cards/:id/labels`

Detach a label from a card. Same body as above.

### Label (protected)

#### `POST /api/v1/labels`

Create a new label. `name` and `color` are required.

**Request body**

```json
{ "name": "Urgent", "color": "#FF0000" }
```

**Response `200 OK`**

```json
{
  "status": "Success",
  "response_code": 200,
  "message": "Successfully Created Label",
  "data": {
    "internal_id": 1,
    "public_id": "b3f1...uuid",
    "name": "Urgent",
    "color": "#FF0000"
  }
}
```

#### `GET /api/v1/labels`

List every label available in the system.

#### `GET /api/v1/labels/:id`

Get a single label. `:id` = `public_id`.

#### `PUT /api/v1/labels/:id`

Partially update a label — only the fields present in the body are changed. `:id` = `public_id`.

**Request body**

```json
{ "name": "Stuck", "color": "#FF3210" }
```

#### `DELETE /api/v1/labels/:id`

Delete a label. `:id` = `public_id`. Any card currently using this label loses the association automatically (`ON DELETE CASCADE` on `card_labels`).

---

## Implementation Status

| Module                        | Model | Migration SQL | Repository | Service | Controller | Route |
| ----------------------------- | :---: | :-----------: | :--------: | :-----: | :--------: | :---: |
| User / Auth                   |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| User — Pagination/Filter/Sort |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| Board                         |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| Board Member                  |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| List                          |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| List Position (reorder)       |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| Card                          |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| Card Position (reorder)       |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| Card ↔ Label (attach/detach)  |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| **Label CRUD (standalone)**   |  ✅   |      ✅       |     ✅     |   ✅    |     ✅     |  ✅   |
| Card Attachment               |  ✅   |      ✅       |     ✅     |   ✅    |     ❌     |  ❌   |
| Card Assignee                 |  ✅   |      ✅       |     ❌     |   ❌    |     ❌     |  ❌   |
| Comment                       |  ✅   |      ✅       |     ❌     |   ❌    |     ❌     |  ❌   |
| Refresh token redeem          |   –   |       –       |     –      |   ❌    |     ❌     |  ❌   |

---

## Deployment

1. Build a binary for the target OS/arch:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bin/server main.go
   ```
2. Provision PostgreSQL (managed or self-hosted) and run all migrations in `database/migrations`.
3. Set production environment variables via your platform's secret manager (Railway/Fly.io/Docker secrets/systemd EnvironmentFile).
4. Run the binary behind a reverse proxy (Nginx/Caddy) for TLS termination, or expose it through a platform that already handles TLS (Railway, Fly.io, Render).
