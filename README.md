# Prospect Backend

Backend server for **Prospect**, a multiplayer mobile tactical shooter.

The server provides authentication, player management, item and inventory systems, loot, weapon configurations, and the foundation for multiplayer match management.

> 🚧 **Project status:** In active development

---

## Tech Stack

* **Go 1.26+** — backend programming language
* **PostgreSQL 18+** — relational database
* **pgx/v5** — PostgreSQL driver and connection pool
* **JWT** — authentication and authorization
* **REST API** — client-server communication
* **bcrypt** — password hashing
* **Git** — version control

---

## Architecture

The backend follows a layered architecture:

```text
HTTP Request
     │
     ▼
 Handler
     │
     ▼
 Service
     │
     ▼
 Repository
     │
     ▼
 PostgreSQL
```

### Handler

Responsible for:

* HTTP requests and responses
* request parsing
* UUID validation
* JSON encoding/decoding
* HTTP status codes

### Service

Responsible for:

* business logic
* validation
* interaction between modules
* enforcing application rules

### Repository

Responsible for:

* PostgreSQL queries
* database operations
* persistence

### Model

Contains domain structures used by the corresponding module.

---

## Project Structure

```text
prospect-backend/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── player/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── item/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── inventory/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── loot/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── weapon/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   └── database/
│       ├── database.go
│       └── migrations/
│
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

# Authentication

Prospect uses **JWT-based authentication**.

The authentication flow:

```text
Client
  │
  │ Register / Login
  ▼
Auth Handler
  │
  ▼
Auth Service
  │
  ▼
PostgreSQL
  │
  ▼
JWT Token
  │
  ▼
Client
```

Protected endpoints require:

```http
Authorization: Bearer <token>
```

The authentication middleware validates the JWT and places the authenticated user's ID into the request context.

---

# Database

The project uses **PostgreSQL** with `pgx/v5` and `pgxpool`.

Database migrations are stored in:

```text
internal/database/migrations/
```

Current database systems include:

```text
Users
  │
  └── Players
        │
        ├── Inventory
        │     │
        │     └── Items
        │
        └── Game systems
```

Foreign keys and database constraints are used to maintain data integrity.

---

# Modules

## Auth Module

Responsible for:

* user registration
* user login
* password hashing
* JWT generation
* JWT validation
* authentication middleware

---

## Player Module

Responsible for player-related game data.

Current player data includes:

* XP
* level
* cash
* premium currency

The player is linked to the authenticated user.

---

## Item Module

Defines the base item configuration used by game systems.

An item represents **what an object is**.

Examples:

```text
AK-47
Medkit
Ammo
Armor
```

Item properties include concepts such as:

```text
name
type
rarity
stackable
```

---

## Inventory Module

Responsible for player-owned items.

The separation is:

```text
Item
  ↓
What is this?

Inventory
  ↓
What does this player own?
```

For example:

```text
Player
  ↓
Inventory
  ↓
AK-47 × 1
```

---

## Loot Module

Responsible for game loot and the relationship between loot and available items.

The loot system is intended to become one of the foundations of the in-game economy and extraction gameplay.

---

## Weapon Module

Responsible for weapon configurations.

A weapon extends an `Item` with combat-related characteristics:

```text
Item
  ↓
Weapon
```

Current weapon properties include:

* damage
* fire rate
* magazine size
* reload time
* range

Example:

```text
AK-47

Damage:        32
Fire Rate:     600
Magazine Size: 30
Reload Time:   2.4
Range:         50
```

The `Weapon Module` does not determine whether a player owns the weapon.

Ownership belongs to the `Inventory Module`.

Detailed documentation:

```text
internal/weapon/README.md
```

---

# Module Relationships

The current game architecture can be represented as:

```text
                    ┌─────────────┐
                    │    Auth     │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Player    │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Inventory  │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │    Item     │
                    └──────┬──────┘
                           │
                ┌──────────┴──────────┐
                ▼                     ▼
          ┌───────────┐         ┌───────────┐
          │  Weapon   │         │    Loot   │
          └───────────┘         └───────────┘
```

This separation allows game systems to evolve independently.

For example:

```text
Item
  │
  ├── Weapon
  ├── Armor
  ├── Consumable
  └── Ammunition
```

without putting ownership logic directly into the item definitions.

---

# API

Current API functionality includes:

### Authentication

```http
POST /api/auth/register
POST /api/auth/login
```

### Player

```http
GET /api/player/me
```

### Weapons

```http
GET /api/weapons
GET /api/weapons/{id}
GET /api/weapons/item/{itemID}
```

Additional endpoints are implemented by the corresponding modules as development continues.

---

# Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=your_database_url
JWT_SECRET=your_secret_key
PORT=8080
```

Never commit `.env` to Git.

The file is included in `.gitignore`.

---

# Installation

Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/prospect-backend.git
cd prospect-backend
```

Download dependencies:

```bash
go mod download
```

Configure the `.env` file.

---

# Run the Server

Start the development server:

```bash
go run ./cmd/server
```

The default development port is:

```text
:8080
```

---

# Development Commands

Format the project:

```bash
gofmt -w .
```

Run tests:

```bash
go test ./...
```

Build the project:

```bash
go build ./...
```

Run the server:

```bash
go run ./cmd/server
```

---

# Roadmap

## Core Backend

* [x] User registration
* [x] User login
* [x] JWT authentication
* [x] PostgreSQL integration
* [x] Player module
* [x] Item module
* [x] Inventory module
* [x] Loot module
* [x] Weapon module
* [x] Database migrations

## Match System

* [ ] Match module
* [ ] Match creation
* [ ] Player joining
* [ ] Maximum player limit
* [ ] Match states
* [ ] Match start
* [ ] Match completion
* [ ] Match rewards
* [ ] Concurrent player joining protection

## Gameplay

* [ ] Combat system
* [ ] Weapon firing
* [ ] Damage calculation
* [ ] Player death
* [ ] Loot interaction
* [ ] Extraction
* [ ] Match rewards
* [ ] Player progression

## Real-Time Backend

* [ ] WebSocket communication
* [ ] Match goroutines
* [ ] Match event channels
* [ ] Concurrent game state
* [ ] Graceful match shutdown

## Future Systems

* [ ] Player statistics
* [ ] Skins
* [ ] Economy
* [ ] Matchmaking
* [ ] Game server
* [ ] Production deployment
* [ ] Monitoring and observability

---

# Development Philosophy

Prospect is being developed as a modular Go backend.

The project focuses on:

* clear separation of responsibilities
* explicit business logic
* database integrity
* secure authentication
* maintainable Go code
* modular game systems
* concurrency-safe multiplayer architecture

The goal is not only to build the game backend, but also to use Prospect as a practical backend engineering project covering real-world Go development patterns.

---

# Project Status

🚧 **In development**

The core backend foundation is implemented.

Current systems include:

```text
Authentication
      ↓
Player
      ↓
Inventory
      ↓
Items
      ↓
Loot / Weapons
```

The next major development stage is the **Match System**, which will introduce match lifecycle management, player participation, concurrency control, and eventually real-time game communication.

---

# License

This project is currently private and is not licensed for redistribution.
