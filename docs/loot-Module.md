# Loot Module

## Зміст

1. [Що таке Loot Module](#s1)
2. [Архітектура Loot Module](#s2)

   * 2.1 [Model](#s2-1)
   * 2.2 [Service](#s2-2)
   * 2.3 [Repository](#s2-3)
   * 2.4 [Handler](#s2-4)
3. [Loot Model](#s3)
4. [Створення Loot](#s4)
5. [Отримання Loot](#s5)

   * 5.1 [Отримати весь Loot](#s5-1)
   * 5.2 [Отримати Loot за ID](#s5-2)
   * 5.3 [Невалідний ID](#s5-3)
   * 5.4 [Loot не знайдений](#s5-4)
6. [Видалення Loot](#s6)
7. [Pickup Loot](#s7)

   * 7.1 [Pickup Flow](#s7-1)
   * 7.2 [Player Resolution](#s7-2)
   * 7.3 [Transaction](#s7-3)
   * 7.4 [Concurrency Protection](#s7-4)
8. [Database](#s8)

   * 8.1 [Таблиця loot](#s8-1)
   * 8.2 [Foreign Key](#s8-2)
   * 8.3 [Constraints](#s8-3)
9. [Error Handling](#s9)
10. [API Endpoints](#s10)
11. [Структура файлів](#s11)
12. [Повний Flow](#s12)
13. [Current Status](#s13)
14. [Next Module](#s14)

---

<a name="s1"></a>

# 1. Що таке Loot Module

Loot Module відповідає за предмети, які знаходяться безпосередньо на карті гри.

На відміну від Inventory, Loot не належить конкретному Player.

Inventory:

```text
Player
  │
  └── Inventory
        │
        ├── AK-47
        ├── Ammo
        └── Medkit
```

Loot:

```text
Map
 │
 ├── Loot
 │     ├── Ammo
 │     ├── Medkit
 │     └── Weapon
 │
 └── Players
```

Loot має:

* `item_id` — який Item лежить на карті
* `quantity` — кількість предметів
* `position_x` — X координата
* `position_y` — Y координата

Наприклад:

```json
{
    "id": "7c1c7b4f-6d7f-4c3d-9d7e-1a8e5b6c1234",
    "item_id": "2f4d8c91-7e13-4f3d-a812-5e4d7b9c3210",
    "quantity": 60,
    "position_x": 210.2,
    "position_y": 94.7,
    "created_at": "2026-09-05T18:00:00Z"
}
```

Loot Module відповідає за:

```text
CREATE
GET
DELETE
PICKUP
```

Pickup переміщує Loot у Inventory Player.

---

<a name="s2"></a>

# 2. Архітектура Loot Module

Loot Module використовує стандартну для проєкту структуру:

```text
HTTP Request
     │
     ▼
┌──────────────┐
│   Handler    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Service    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  PostgreSQL  │
└──────────────┘
```

Кожен шар має свою відповідальність.

```text
Handler

↓

HTTP

Service

↓

Business Logic

Repository

↓

Database

Model

↓

Data Structure
```

Pickup має додаткову взаємодію з Player:

```text
JWT
 │
 ▼
user_id
 │
 ▼
Player Repository
 │
 ▼
player_id
 │
 ▼
Loot Repository
 │
 ▼
Inventory
```

---

<a name="s2-1"></a>

## 2.1 Model

Файл:

```text
internal/loot/model.go
```

Model описує структуру Loot.

```go
package loot

import (
    "time"

    "github.com/google/uuid"
)

type Loot struct {
    ID        uuid.UUID `json:"id"`
    ItemID    uuid.UUID `json:"item_id"`
    Quantity  int64     `json:"quantity"`
    PositionX float32   `json:"position_x"`
    PositionY float32   `json:"position_y"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Поля

| Поле      | Тип       | Опис                        |
| --------- | --------- | --------------------------- |
| ID        | uuid.UUID | Унікальний ID Loot          |
| ItemID    | uuid.UUID | ID предмета з таблиці items |
| Quantity  | int64     | Кількість предметів         |
| PositionX | float32   | X координата на карті       |
| PositionY | float32   | Y координата на карті       |
| CreatedAt | time.Time | Час створення               |

---

<a name="s2-2"></a>

## 2.2 Service

Файл:

```text
internal/loot/service.go
```

Service відповідає за business logic.

Основні операції:

```text
Create
GetByID
GetAll
Delete
Pickup
```

Service перевіряє вхідні дані.

Наприклад:

```go
if loot.ItemID == uuid.Nil {
    return Loot{}, ErrInvalidItemID
}
```

та:

```go
if loot.Quantity <= 0 {
    return Loot{}, ErrInvalidQuantity
}
```

Для Pickup Service:

1. перевіряє `userID`
2. перевіряє `lootID`
3. знаходить Player через `userID`
4. отримує `player.ID`
5. передає `player.ID` та `lootID` у Repository

```go
func (s *Service) Pickup(
    ctx context.Context,
    userID uuid.UUID,
    lootID uuid.UUID,
) error {
    if userID == uuid.Nil {
        return ErrInvalidPlayerID
    }

    if lootID == uuid.Nil {
        return ErrInvalidLootID
    }

    player, err := s.playerRepository.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }

    return s.repository.PickupTx(
        ctx,
        player.ID,
        lootID,
    )
}
```

Таким чином Handler не працює безпосередньо з Player ID.

Правильний flow:

```text
Handler
   │
   ▼
Service
   │
   ├── Validate
   │
   ├── Resolve Player
   │
   ▼
Repository
```

---

<a name="s2-3"></a>

## 2.3 Repository

Файл:

```text
internal/loot/repository.go
```

Repository відповідає за роботу з PostgreSQL.

Методи:

```go
Create()
GetByID()
GetAll()
Delete()
PickupTx()
```

Repository не повинен містити game logic.

Наприклад:

```go
func (r *Repository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (Loot, error)
```

Він просто виконує SQL:

```sql
SELECT
    id,
    item_id,
    quantity,
    position_x,
    position_y,
    created_at
FROM loot
WHERE id = $1;
```

### PickupTx

Pickup виконується через database transaction:

```go
func (r *Repository) PickupTx(
    ctx context.Context,
    playerID uuid.UUID,
    lootID uuid.UUID,
) error
```

Операція виконує:

```text
BEGIN
   │
   ├── Find Loot
   │
   ├── Lock Loot
   │
   ├── Add Item to Inventory
   │
   ├── Delete Loot
   │
   ▼
 COMMIT
```

Якщо будь-яка операція завершується помилкою:

```text
ROLLBACK
```

---

<a name="s2-4"></a>

## 2.4 Handler

Файл:

```text
internal/loot/handler.go
```

Handler відповідає за HTTP.

Його завдання:

```text
HTTP Request
     ↓
Parse URL / Body
     ↓
Validate HTTP input
     ↓
Call Service
     ↓
HTTP Response
```

Для Create Handler отримує JSON:

```json
{
    "item_id": "2f4d8c91-7e13-4f3d-a812-5e4d7b9c3210",
    "quantity": 60,
    "position_x": 210.2,
    "position_y": 94.7
}
```

Для Pickup Handler отримує `lootID` з URL:

```text
POST /api/loot/pickup/{lootID}
```

та `userID` з JWT context:

```go
userIDString, ok := auth.UserIDFromContext(r)
```

Після цього Handler передає дані Service.

---

<a name="s3"></a>

# 3. Loot Model

Loot пов'язаний з Item через `item_id`.

Приклад:

```text
items

┌──────────────────────────────────┐
│ id = abc                         │
│ name = 5.56 Ammo                 │
└──────────────────────────────────┘
                ▲
                │
                │ item_id
                │
loot            │
┌──────────────────────────────────┐
│ id = xyz                         │
│ item_id = abc                    │
│ quantity = 60                    │
│ position_x = 210.2               │
│ position_y = 94.7                │
└──────────────────────────────────┘
```

Тобто Loot не дублює інформацію про Item.

Loot зберігає тільки:

```text
item_id
```

А інформація про сам Item знаходиться в:

```text
items
```

Це дозволяє розділити відповідальність:

```text
Item

↓

Що це за предмет?

Loot

↓

Де він знаходиться і скільки його?
```

---

<a name="s4"></a>

# 4. Створення Loot

Створення Loot виконується через:

```http
POST /api/loot
```

Request body:

```json
{
    "item_id": "2f4d8c91-7e13-4f3d-a812-5e4d7b9c3210",
    "quantity": 60,
    "position_x": 210.2,
    "position_y": 94.7
}
```

Service перевіряє:

```text
ItemID != uuid.Nil
Quantity > 0
```

Flow:

```text
POST /api/loot
       │
       ▼
Parse JSON
       │
       ▼
Service.Create()
       │
       ├── Validate ItemID
       │
       ├── Validate Quantity
       │
       ▼
Repository.Create()
       │
       ▼
INSERT INTO loot
       │
       ▼
PostgreSQL
```

Repository виконує:

```sql
INSERT INTO loot (
    item_id,
    quantity,
    position_x,
    position_y
)
VALUES ($1, $2, $3, $4)
RETURNING
    id,
    item_id,
    quantity,
    position_x,
    position_y,
    created_at;
```

Успішне створення повертає:

```http
201 Created
```

та створений Loot у JSON.

---

<a name="s5"></a>

# 5. Отримання Loot

Loot можна отримати двома способами:

```text
GET /api/loot
```

або:

```text
GET /api/loot/{id}
```

---

<a name="s5-1"></a>

## 5.1 Отримати весь Loot

Endpoint:

```http
GET /api/loot

Authorization: Bearer <token>
```

Handler викликає:

```go
loots, err := h.service.GetAll(r.Context())
```

Service:

```go
func (s *Service) GetAll(
    ctx context.Context,
) ([]Loot, error) {
    return s.repository.GetAll(ctx)
}
```

Repository:

```sql
SELECT
    id,
    item_id,
    quantity,
    position_x,
    position_y,
    created_at
FROM loot
ORDER BY created_at ASC;
```

Response:

```json
[
    {
        "id": "7c1c7b4f-6d7f-4c3d-9d7e-1a8e5b6c1234",
        "item_id": "2f4d8c91-7e13-4f3d-a812-5e4d7b9c3210",
        "quantity": 60,
        "position_x": 210.2,
        "position_y": 94.7,
        "created_at": "2026-09-05T18:00:00Z"
    }
]
```

При відсутності Loot Repository повертає порожній slice:

```json
[]
```

---

<a name="s5-2"></a>

## 5.2 Отримати Loot за ID

Endpoint:

```http
GET /api/loot/{lootID}

Authorization: Bearer <token>
```

Handler отримує ID:

```go
idString := strings.TrimPrefix(
    r.URL.Path,
    "/api/loot/",
)
```

Після цього UUID парситься:

```go
id, err := uuid.Parse(idString)
```

Якщо UUID валідний:

```text
Handler
   ↓
Service.GetByID()
   ↓
Repository.GetByID()
   ↓
PostgreSQL
```

Repository:

```sql
SELECT
    id,
    item_id,
    quantity,
    position_x,
    position_y,
    created_at
FROM loot
WHERE id = $1;
```

Response:

```json
{
    "id": "7c1c7b4f-6d7f-4c3d-9d7e-1a8e5b6c1234",
    "item_id": "2f4d8c91-7e13-4f3d-a812-5e4d7b9c3210",
    "quantity": 60,
    "position_x": 210.2,
    "position_y": 94.7,
    "created_at": "2026-09-05T18:00:00Z"
}
```

---

<a name="s5-3"></a>

## 5.3 Невалідний ID

Наприклад:

```http
GET /api/loot/hello
```

UUID parse завершиться помилкою:

```go
id, err := uuid.Parse(idString)
```

Handler повертає:

```http
400 Bad Request
```

Response:

```text
invalid loot id
```

---

<a name="s5-4"></a>

## 5.4 Loot не знайдений

Якщо UUID валідний, але такого Loot немає:

```http
GET /api/loot/7c1c7b4f-6d7f-4c3d-9d7e-1a8e5b6c1234
```

Repository поверне помилку `pgx.ErrNoRows`.

Handler повертає:

```http
404 Not Found
```

Response:

```text
loot not found
```

Таким чином:

```text
Invalid UUID
    ↓
400 Bad Request

Valid UUID
    ↓
Loot exists
    ↓
200 OK

Valid UUID
    ↓
Loot doesn't exist
    ↓
404 Not Found
```

---

<a name="s6"></a>

# 6. Видалення Loot

Repository підтримує видалення:

```go
func (r *Repository) Delete(
    ctx context.Context,
    id uuid.UUID,
) error
```

SQL:

```sql
DELETE FROM loot
WHERE id = $1;
```

Service:

```go
func (s *Service) Delete(
    ctx context.Context,
    id uuid.UUID,
) error {
    return s.repository.Delete(ctx, id)
}
```

На поточному етапі окремий HTTP endpoint для Delete не реалізований.

Loot видаляється через game logic, зокрема під час Pickup.

---

<a name="s7"></a>

# 7. Pickup Loot

Pickup дозволяє Player забрати Loot з карти та додати його до Inventory.

Endpoint:

```http
POST /api/loot/pickup/{lootID}

Authorization: Bearer <token>
```

У Pickup використовуються два різних ID:

```text
user_id
loot_id
```

`user_id` береться з JWT.

`loot_id` береться з URL.

Player ID безпосередньо з JWT не передається.

---

<a name="s7-1"></a>

## 7.1 Pickup Flow

Повний flow:

```text
Client
  │
  │ POST /api/loot/pickup/{lootID}
  │ Authorization: Bearer JWT
  ▼
Auth Middleware
  │
  │ user_id
  ▼
Loot Handler
  │
  ├── Parse loot_id
  │
  ├── Get user_id from context
  │
  ▼
Loot Service
  │
  ├── Validate user_id
  │
  ├── Validate loot_id
  │
  ├── Find Player
  │
  ▼
Player Repository
  │
  │ GetByUserID()
  ▼
player_id
  │
  ▼
Loot Repository
  │
  ▼
PickupTx()
  │
  ├── Get Loot
  ├── Add Item to Inventory
  └── Delete Loot
  │
  ▼
COMMIT
  │
  ▼
Handler
  │
  ▼
200 OK
```

Успішна відповідь:

```http
200 OK
```

```json
{
    "message": "loot picked up successfully"
}
```

---

<a name="s7-2"></a>

## 7.2 Player Resolution

JWT містить:

```text
user_id
```

А Inventory використовує:

```text
player_id
```

Тому Pickup використовує наступний flow:

```text
JWT
 │
 ▼
user_id
 │
 ▼
Player.GetByUserID()
 │
 ▼
player.id
 │
 ▼
Inventory
```

Це важливо, оскільки `user_id` та `player_id` є різними сутностями.

Структура:

```text
User
 │
 └── Player
       │
       └── Inventory
```

---

<a name="s7-3"></a>

## 7.3 Transaction

Pickup виконується в одній database transaction.

```text
BEGIN
   │
   ├── SELECT Loot
   │
   ├── Add Item to Inventory
   │
   ├── DELETE Loot
   │
   ▼
 COMMIT
```

Якщо додавання Item до Inventory не вдалося:

```text
ROLLBACK
```

Якщо видалення Loot не вдалося:

```text
ROLLBACK
```

Це запобігає ситуації:

```text
Item added to Inventory
        ↓
Delete Loot failed
        ↓
Loot still exists
```

Такий сценарій міг би створити duplication exploit.

Правильна поведінка:

```text
BEGIN

Add Item
   ↓
Success

Delete Loot
   ↓
Success

COMMIT
```

або:

```text
BEGIN

Operation failed

ROLLBACK
```

---

<a name="s7-4"></a>

## 7.4 Concurrency Protection

Перед отриманням Loot використовується:

```sql
SELECT item_id, quantity
FROM loot
WHERE id = $1
FOR UPDATE;
```

`FOR UPDATE` блокує конкретний Loot row до завершення transaction.

Це захищає від одночасного Pickup одного Loot двома запитами.

Без блокування потенційно можливий сценарій:

```text
Request A → Find Loot
Request B → Find Loot

Request A → Add Item
Request B → Add Item

Request A → Delete Loot
Request B → Delete Loot
```

З `FOR UPDATE`:

```text
Request A
   │
   ▼
LOCK Loot
   │
   ├── Add Item
   ├── Delete Loot
   └── COMMIT
          │
          ▼
       Unlock

Request B
   │
   ▼
Cannot use already processed Loot
```

Таким чином Pickup захищений від race condition на рівні database transaction.

---

<a name="s8"></a>

# 8. Database

Файл migration:

```text
internal/database/migrations/007_create_loot.sql
```

---

<a name="s8-1"></a>

## 8.1 Таблиця loot

```sql
CREATE TABLE loot (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    item_id UUID NOT NULL
        REFERENCES items(id)
        ON DELETE RESTRICT,

    quantity BIGINT NOT NULL DEFAULT 1,

    position_x REAL NOT NULL,
    position_y REAL NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT loot_quantity_positive
        CHECK (quantity > 0)
);
```

Структура:

| Column     | Type        | Constraint    |
| ---------- | ----------- | ------------- |
| id         | UUID        | PRIMARY KEY   |
| item_id    | UUID        | NOT NULL, FK  |
| quantity   | BIGINT      | NOT NULL, > 0 |
| position_x | REAL        | NOT NULL      |
| position_y | REAL        | NOT NULL      |
| created_at | TIMESTAMPTZ | NOT NULL      |

---

<a name="s8-2"></a>

## 8.2 Foreign Key

`item_id` посилається на:

```text
items.id
```

```sql
item_id UUID NOT NULL
    REFERENCES items(id)
    ON DELETE RESTRICT
```

Це означає:

```text
Item
 │
 └── Loot
```

Не можна видалити Item, якщо він використовується Loot.

Це захищає database integrity.

---

<a name="s8-3"></a>

## 8.3 Constraints

### Quantity

```sql
CHECK (quantity > 0)
```

Не можна створити:

```text
quantity = 0
```

або:

```text
quantity = -10
```

### Item

```sql
item_id UUID NOT NULL
```

Loot завжди повинен посилатися на Item.

### Coordinates

```sql
position_x REAL NOT NULL
position_y REAL NOT NULL
```

Loot завжди має позицію на карті.

---

<a name="s9"></a>

# 9. Error Handling

Loot Module використовує validation на Service рівні.

Основні errors:

```go
var (
    ErrInvalidItemID   = errors.New("item id is required")
    ErrInvalidQuantity = errors.New("quantity must be greater than zero")
    ErrInvalidLootID   = errors.New("loot id is required")
    ErrInvalidPlayerID = errors.New("player id is required")
)
```

### Invalid Item ID

```text
ItemID == uuid.Nil
```

Результат:

```text
ErrInvalidItemID
```

HTTP:

```text
400 Bad Request
```

### Invalid Quantity

```text
Quantity <= 0
```

Результат:

```text
ErrInvalidQuantity
```

HTTP:

```text
400 Bad Request
```

### Invalid Loot ID

Невалідний UUID у URL:

```text
ErrInvalidLootID
```

HTTP:

```text
400 Bad Request
```

### Unauthorized

Якщо JWT context не містить `user_id`:

```text
401 Unauthorized
```

### Loot Not Found

Якщо Loot не існує:

```text
404 Not Found
```

### Database Error

Непередбачена помилка database:

```text
500 Internal Server Error
```

### HTTP Statuses

| Situation         | Status |
| ----------------- | -----: |
| Invalid UUID      |    400 |
| Invalid Item ID   |    400 |
| Invalid Quantity  |    400 |
| Unauthorized      |    401 |
| Loot not found    |    404 |
| Database error    |    500 |
| Successful GET    |    200 |
| Successful Pickup |    200 |
| Successful Create |    201 |

---

<a name="s10"></a>

# 10. API Endpoints

Loot endpoints захищені JWT middleware.

### Get all Loot

```http
GET /api/loot

Authorization: Bearer <token>
```

Response:

```http
200 OK
```

---

### Create Loot

```http
POST /api/loot

Authorization: Bearer <token>
Content-Type: application/json
```

Request:

```json
{
    "item_id": "2f4d8c91-7e13-4f3d-a812-5e4d7b9c3210",
    "quantity": 60,
    "position_x": 210.2,
    "position_y": 94.7
}
```

Response:

```http
201 Created
```

---

### Get Loot by ID

```http
GET /api/loot/{lootID}

Authorization: Bearer <token>
```

Response:

```http
200 OK
```

---

### Invalid ID

```http
GET /api/loot/hello

Authorization: Bearer <token>
```

Response:

```http
400 Bad Request
```

---

### Nonexistent Loot

```http
GET /api/loot/{nonexistentUUID}

Authorization: Bearer <token>
```

Response:

```http
404 Not Found
```

---

### Pickup Loot

```http
POST /api/loot/pickup/{lootID}

Authorization: Bearer <token>
```

Response:

```http
200 OK
```

```json
{
    "message": "loot picked up successfully"
}
```

Pickup:

```text
Loot
  │
  ├── item_id
  └── quantity
        │
        ▼
    Inventory
        │
        ▼
    Loot DELETE
```

---

### Поточні endpoints

```text
GET  /api/loot
POST /api/loot
GET  /api/loot/{id}
POST /api/loot/pickup/{id}
```

Окремого client-facing endpoint:

```text
DELETE /api/loot/{id}
```

поки немає.

Видалення Loot використовується внутрішньо під час Pickup.

---

<a name="s11"></a>

# 11. Структура файлів

Поточна структура:

```text
internal/

├── auth/
│   ├── handler.go
│   ├── middleware.go
│   ├── model.go
│   ├── repository.go
│   ├── service.go
│   ├── token.go
│   └── context.go
│
├── player/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── inventory/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── item/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── weapon/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── loot/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
└── database/
    ├── postgres.go
    └── migrations/
        ├── 002_create_players.sql
        ├── 003_create_inventory_items.sql
        ├── 004_create_items.sql
        ├── 005_add_inventory_item_fk.sql
        ├── 006_create_weapons.sql
        └── 007_create_loot.sql
```

---

<a name="s12"></a>

# 12. Повний Flow

## GET /api/loot

```text
Client
  │
  │ GET /api/loot
  │ Authorization: Bearer JWT
  ▼
Auth Middleware
  │
  │ Token valid
  ▼
Loot Handler
  │
  ▼
Loot Service
  │
  ▼
Loot Repository
  │
  │ SELECT ...
  ▼
PostgreSQL
  │
  │ []Loot
  ▼
Repository
  │
  ▼
Service
  │
  ▼
Handler
  │
  │ JSON
  ▼
Client
```

---

## POST /api/loot

```text
Client
  │
  │ POST /api/loot
  │ JSON
  ▼
Auth Middleware
  │
  ▼
Loot Handler
  │
  │ Parse JSON
  ▼
Loot Service
  │
  ├── Validate ItemID
  ├── Validate Quantity
  │
  ▼
Loot Repository
  │
  │ INSERT
  ▼
PostgreSQL
  │
  │ Created Loot
  ▼
Handler
  │
  │ 201 Created
  ▼
Client
```

---

## GET /api/loot/{id}

```text
Client
  │
  │ GET /api/loot/{id}
  ▼
Auth Middleware
  │
  ▼
Handler
  │
  │ Parse UUID
  ▼
Service
  │
  ▼
Repository
  │
  │ SELECT ... WHERE id = $1
  ▼
PostgreSQL
  │
  ├── Found
  │     ↓
  │   Loot
  │
  └── Not Found
        ↓
       Error
  ▼
Handler
  │
  ├── 200 OK
  │
  └── 404 Not Found
```

---

## POST /api/loot/pickup/{id}

```text
Client
  │
  │ POST /api/loot/pickup/{lootID}
  │ Authorization: Bearer JWT
  ▼
Auth Middleware
  │
  │ user_id
  ▼
Loot Handler
  │
  ├── Parse lootID
  │
  └── Get userID
  │
  ▼
Loot Service
  │
  ├── Validate IDs
  │
  └── Get Player by userID
  │
  ▼
Player Repository
  │
  │ GetByUserID()
  ▼
player_id
  │
  ▼
Loot Repository
  │
  │ BEGIN
  ▼
SELECT Loot
FOR UPDATE
  │
  ▼
Get item_id + quantity
  │
  ▼
INSERT / UPDATE Inventory
  │
  ▼
DELETE Loot
  │
  ▼
COMMIT
  │
  ▼
Handler
  │
  │ 200 OK
  ▼
Client
```

---

<a name="s13"></a>

# 13. Current Status

Loot Module реалізований.

### Database

```text
✓ loot table
✓ item_id foreign key
✓ quantity constraint
✓ position_x
✓ position_y
✓ created_at
```

### Model

```text
✓ Loot struct
```

### Repository

```text
✓ Create
✓ GetByID
✓ GetAll
✓ Delete
✓ PickupTx
✓ Database transaction
✓ FOR UPDATE locking
```

### Service

```text
✓ Create
✓ GetByID
✓ GetAll
✓ Delete
✓ Pickup
✓ ItemID validation
✓ Quantity validation
✓ Loot ID validation
✓ Player ID validation
✓ Player resolution через user_id
```

### Handler

```text
✓ GetAll
✓ GetByID
✓ Create
✓ Pickup
✓ UUID validation
✓ Request body parsing
✓ 404 handling
✓ JWT context integration
```

### Authentication

```text
✓ JWT protected routes
✓ user_id extracted from JWT
```

### Pickup

```text
✓ Loot → Inventory
✓ Player resolution
✓ Inventory stacking
✓ Loot deletion
✓ Transaction
✓ Rollback
✓ FOR UPDATE protection
✓ End-to-end tested
```

Поточний flow:

```text
Item
  │
  ▼
Loot
  │
  ▼
Map
  │
  │ Pickup
  ▼
Player
  │
  ▼
Inventory
```

Pickup вже протестований end-to-end:

```text
POST /api/loot
        ↓
Loot created
        ↓
POST /api/loot/pickup/{loot_id}
        ↓
Loot moved to Inventory
        ↓
Loot deleted
        ↓
GET /api/inventory
        ↓
Item quantity updated
```

---

<a name="s14"></a>

# 14. Next Module

Loot Module завершений на поточному етапі.

Наступний логічний модуль:

```text
Loadout
```

Loadout буде відповідати за спорядження Player перед грою.

Очікуваний flow:

```text
Player
  │
  ▼
Loadout
  │
  ├── Weapon
  ├── Armor
  ├── Equipment
  └── Items
```

У подальшому Loadout буде пов'язаний з Match:

```text
Player
  │
  ▼
Loadout
  │
  ▼
Match
  │
  ▼
Spawn
```

Після реалізації Loadout backend поступово переходить від CRUD-модулів до основної gameplay logic Prospect.
