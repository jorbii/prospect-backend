````markdown
# Loot Module

## Зміст

1. [Що таке Loot Module](#s1)
2. [Архітектура Loot Module](#s2)
   - 2.1 [Model](#s2-1)
   - 2.2 [Service](#s2-2)
   - 2.3 [Repository](#s2-3)
   - 2.4 [Handler](#s2-4)
3. [Loot Model](#s3)
4. [Створення Loot](#s4)
5. [Отримання Loot](#s5)
   - 5.1 [Отримати весь Loot](#s5-1)
   - 5.2 [Отримати Loot за ID](#s5-2)
   - 5.3 [Невалідний ID](#s5-3)
   - 5.4 [Loot не знайдений](#s5-4)
6. [Видалення Loot](#s6)
7. [Database](#s7)
   - 7.1 [Таблиця loot](#s7-1)
   - 7.2 [Foreign Key](#s7-2)
   - 7.3 [Constraints](#s7-3)
8. [Error Handling](#s8)
9. [API Endpoints](#s9)
10. [Структура файлів](#s10)
11. [Повний Flow](#s11)
12. [Майбутній Pickup Flow](#s12)
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
````

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

Loot Module поки що відповідає тільки за:

```text
CREATE
GET
DELETE
```

Pickup механіка буде реалізована пізніше.

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
```

Також Service перевіряє вхідні дані.

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

Тобто Handler не повинен самостійно вирішувати business rules.

Правильний flow:

```text
Handler
   │
   ▼
Service
   │
   ├── Validate
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

Repository відповідає тільки за роботу з PostgreSQL.

Методи:

```go
Create()
GetByID()
GetAll()
Delete()
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
Parse URL
     ↓
Call Service
     ↓
HTTP Response
```

Наприклад:

```text
GET /api/loot/UUID
```

Handler:

1. отримує ID з URL
2. перевіряє UUID
3. викликає Service
4. повертає JSON
5. повертає відповідний HTTP status

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
│ item_id = abc                   │
│ quantity = 60                   │
│ position_x = 210.2              │
│ position_y = 94.7               │
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

Створення Loot виконується через Service.

Метод:

```go
func (s *Service) Create(
    ctx context.Context,
    loot Loot,
) (Loot, error)
```

Перед створенням перевіряється `ItemID`.

```go
if loot.ItemID == uuid.Nil {
    return Loot{}, ErrInvalidItemID
}
```

Також перевіряється Quantity:

```go
if loot.Quantity <= 0 {
    return Loot{}, ErrInvalidQuantity
}
```

Після validation Service передає дані Repository:

```text
Create Loot
     │
     ▼
Validate ItemID
     │
     ▼
Validate Quantity
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

Repository поверне помилку PostgreSQL.

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

На поточному етапі Delete існує на рівні Repository і Service.

Окремий HTTP endpoint для клієнта поки не реалізований.

Причина:

Loot повинен видалятися не просто через звичайний DELETE, а в контексті game logic.

Наприклад, майбутній Pickup:

```text
Player
   │
   ▼
Pickup Loot
   │
   ▼
Check Loot
   │
   ▼
Add Item to Inventory
   │
   ▼
Delete Loot
```

Ця операція повинна бути transaction-safe.

---

<a name="s7"></a>

# 7. Database

Файл migration:

```text
internal/database/migrations/007_create_loot.sql
```

---

<a name="s7-1"></a>

## 7.1 Таблиця loot

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

<a name="s7-2"></a>

## 7.2 Foreign Key

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

Наприклад:

```text
items
   │
   └── AK-47
          │
          └── loot
```

Якщо Loot існує:

```sql
DELETE FROM items
WHERE name = 'AK-47';
```

PostgreSQL не дозволить це зробити через:

```text
ON DELETE RESTRICT
```

Це захищає database integrity.

---

<a name="s7-3"></a>

## 7.3 Constraints

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

<a name="s8"></a>

# 8. Error Handling

Loot Module використовує validation на Service рівні.

Errors:

```go
var (
    ErrInvalidItemID   = errors.New("item id is required")
    ErrInvalidQuantity = errors.New("quantity must be greater than zero")
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

### Invalid Quantity

```text
Quantity <= 0
```

Результат:

```text
ErrInvalidQuantity
```

### HTTP errors

| Situation      | Status |
| -------------- | -----: |
| Invalid UUID   |    400 |
| Loot not found |    404 |
| Database error |    500 |
| Successful GET |    200 |

---

<a name="s9"></a>

# 9. API Endpoints

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

### Поточні client-facing endpoints

```text
GET     /api/loot
GET     /api/loot/{id}
```

Поки що немає:

```text
POST    /api/loot
DELETE  /api/loot/{id}
POST    /api/loot/{id}/pickup
```

Це зроблено навмисно.

Loot generation та pickup будуть частиною подальшої game logic.

---

<a name="s10"></a>

# 10. Структура файлів

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

<a name="s11"></a>

# 11. Повний Flow

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

<a name="s12"></a>

# 12. Майбутній Pickup Flow

Pickup ще не реалізований.

Але архітектура Loot вже підготовлена під нього.

Коли Player підбирає Loot:

```text
Player
  │
  │ Pickup Loot
  ▼
Loot Service
  │
  ├── Find Loot
  │
  ├── Validate Loot
  │
  ├── Add Item to Inventory
  │
  └── Delete Loot
  │
  ▼
COMMIT
```

Ключовий момент:

```text
Add Item
    +
Delete Loot
```

повинні бути однією database transaction.

Неправильний варіант:

```text
Add Item
   ↓
Success

Delete Loot
   ↓
Failed
```

У такому випадку Player отримав Item, але Loot залишився на карті.

Це може створити duplication exploit.

Правильний варіант:

```text
BEGIN
   │
   ├── Add Item to Inventory
   │
   ├── Delete Loot
   │
   ▼
COMMIT
```

Якщо будь-яка операція не вдалася:

```text
ROLLBACK
```

Тоді стан повертається назад.

---

<a name="s13"></a>

# 13. Current Status

Loot Module реалізований на базовому рівні.

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
```

### Service

```text
✓ Create
✓ GetByID
✓ GetAll
✓ Delete
✓ ItemID validation
✓ Quantity validation
```

### Handler

```text
✓ GetAll
✓ GetByID
✓ UUID validation
✓ 404 handling
```

### Authentication

```text
✓ JWT protected routes
```

### Pickup

```text
⏳ Not implemented yet
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
```

Майбутній flow:

```text
Loot
  │
  │ Pickup
  ▼
Inventory
```

---

<a name="s14"></a>

# 14. Next Module

Після Loot Module наступним логічним етапом є реалізація механіки взаємодії Loot з Player та Inventory.

Основний майбутній flow:

```text
Player
   │
   ▼
Loot
   │
   ▼
Pickup
   │
   ▼
Transaction
   │
   ├── Inventory + Item
   │
   └── Loot DELETE
   │
   ▼
COMMIT
```

Після цього backend матиме вже основний цикл:

```text
Player
   │
   ├── Inventory
   │
   └── Loot
          │
          ▼
       Pickup
          │
          ▼
       Inventory
```

Це є основою для подальшої ігрової логіки Prospect.

```
```
