# Inventory Module

## Зміст

1. Що таке Inventory Module
2. Архітектура
3. Model
4. Repository
5. Service
6. Handler
7. Database
8. Логіка предметів
9. API
10. Взаємодія з іншими модулями
11. Поточний статус

---

## 1. Що таке Inventory Module

`Inventory Module` відповідає за інвентар гравця.

Модуль зберігає інформацію про:

* які предмети має гравець;
* кількість кожного предмета;
* до якого `Player` належить предмет.

Inventory не визначає властивості самого предмета.

Наприклад:

```text
Item:
AK-47
type: weapon
rarity: common
stackable: false

Inventory:
Player A
→ AK-47 × 1
```

Тобто `Item` описує предмет як об'єкт гри, а `Inventory` зберігає його наявність у конкретного гравця.

---

# 2. Архітектура

Inventory Module використовує стандартну backend-архітектуру:

```text
HTTP Request
     ↓
Handler
     ↓
Service
     ↓
Repository
     ↓
PostgreSQL
```

### Handler

Відповідає за HTTP-рівень:

* отримання запиту;
* отримання `item_id`;
* отримання `quantity`;
* HTTP status codes;
* JSON response.

### Service

Містить бізнес-логіку:

* перевірку кількості;
* отримання inventory конкретного гравця;
* додавання предметів;
* видалення предметів.

### Repository

Працює безпосередньо з PostgreSQL:

* SELECT;
* INSERT;
* UPDATE;
* DELETE.

### Model

Описує структуру inventory item у Go.

---

# 3. Model

Файл:

```text
internal/inventory/model.go
```

Основна структура:

```go
type InventoryItem struct {
    ID        uuid.UUID `json:"id"`
    PlayerID  uuid.UUID `json:"player_id"`
    ItemID    uuid.UUID `json:"item_id"`
    Quantity  int64     `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Поля

| Поле        | Тип       | Опис                |
| ----------- | --------- | ------------------- |
| `ID`        | UUID      | ID запису inventory |
| `PlayerID`  | UUID      | ID гравця           |
| `ItemID`    | UUID      | ID предмета         |
| `Quantity`  | int64     | Кількість предметів |
| `CreatedAt` | time.Time | Час створення       |
| `UpdatedAt` | time.Time | Час останньої зміни |

---

# 4. Repository

Файл:

```text
internal/inventory/repository.go
```

Repository відповідає за роботу з таблицею:

```text
inventory_items
```

Основні операції:

```text
GetByPlayerID
AddItem
RemoveItem
DeleteItem
```

## GetByPlayerID

Отримує весь inventory конкретного гравця.

```text
Player ID
   ↓
SELECT inventory_items
WHERE player_id = ?
```

---

## AddItem

Додає предмет до inventory.

Якщо предмета ще немає:

```text
AK-47 × 1
```

Якщо предмет вже існує:

```text
AK-47 × 1
+
AK-47 × 3
=
AK-47 × 4
```

Для цього використовується PostgreSQL:

```sql
ON CONFLICT (player_id, item_id)
DO UPDATE SET
    quantity = inventory_items.quantity + EXCLUDED.quantity
```

Таким чином, для одного гравця не створюються дублікати одного `item_id`.

---

## RemoveItem

Зменшує кількість предмета.

Наприклад:

```text
Medkit × 5
```

Після:

```text
RemoveItem × 2
```

отримуємо:

```text
Medkit × 3
```

Операція виконується тільки якщо поточна кількість достатня.

---

## DeleteItem

Повністю видаляє предмет із inventory.

Наприклад:

```text
Keycard × 1
```

після `DeleteItem`:

```text
inventory → Keycard відсутній
```

---

# 5. Service

Файл:

```text
internal/inventory/service.go
```

Service є бізнес-рівнем Inventory Module.

Він не працює напряму з HTTP або PostgreSQL.

Основні методи:

```text
GetInventory
AddItem
RemoveItem
DeleteItem
```

---

## Валідація quantity

Кількість предметів повинна бути більшою за нуль:

```text
quantity > 0
```

Наприклад:

```text
quantity = 5 → OK
quantity = 1 → OK
quantity = 0 → ERROR
quantity = -1 → ERROR
```

Для цього використовується:

```go
ErrInvalidQuantity
```

---

## Insufficient items

Service також визначає помилку:

```go
ErrInsufficientItems
```

Вона призначена для випадку, коли гравець намагається видалити більше предметів, ніж має.

Наприклад:

```text
Inventory:
Ammo × 5

Request:
Remove Ammo × 10

Result:
Insufficient items
```

---

# 6. Handler

Файл:

```text
internal/inventory/handler.go
```

Handler відповідає за HTTP API.

Inventory використовує authentication middleware.

Схема:

```text
HTTP Request
     ↓
JWT Middleware
     ↓
User ID
     ↓
Player
     ↓
Inventory
```

Гравець не передає `player_id` вручну.

Замість цього сервер отримує `user_id` з JWT:

```text
JWT
 ↓
user_id
 ↓
Player
 ↓
player_id
 ↓
Inventory
```

Це важливо для безпеки.

Клієнт не може просто передати:

```json
{
    "player_id": "чужий-player-id"
}
```

і отримати чужий inventory.

---

# 7. Database

Migration:

```text
internal/database/migrations/003_create_inventory_items.sql
```

Таблиця:

```text
inventory_items
```

Структура:

```text
id
player_id
item_id
quantity
created_at
updated_at
```

## Primary Key

```sql
id UUID PRIMARY KEY
```

Кожен запис inventory має власний ID.

---

## Player relation

```sql
player_id UUID NOT NULL
REFERENCES players(id)
ON DELETE CASCADE
```

Це означає:

```text
Player
  ↓
Inventory Items
```

Якщо Player видаляється, його inventory також видаляється.

---

## Item relation

Inventory також пов'язаний з `items`.

Migration:

```text
internal/database/migrations/005_add_inventory_item_fk.sql
```

Використовується:

```sql
FOREIGN KEY (item_id)
REFERENCES items(id)
ON DELETE RESTRICT
```

Тобто inventory може містити тільки предмет, який існує в `items`.

При цьому Item не можна видалити з каталогу, якщо він вже використовується inventory.

---

# 8. Unique Constraint

Для inventory встановлено:

```sql
UNIQUE (player_id, item_id)
```

Це означає, що один гравець не може мати два окремих записи одного предмета.

Погано:

```text
Player A
├── AK-47 × 2
└── AK-47 × 3
```

Правильно:

```text
Player A
└── AK-47 × 5
```

Саме тому `AddItem` збільшує `quantity`, а не створює новий запис.

---

# 9. Quantity Constraint

У таблиці встановлено:

```sql
CHECK (quantity > 0)
```

Тому в database неможливо зберегти:

```text
quantity = 0
quantity = -5
```

Inventory завжди містить тільки позитивну кількість предметів.

---

# 10. API

Inventory endpoints захищені JWT authentication middleware.

## GET /api/inventory

Отримує inventory поточного гравця.

```http
GET /api/inventory
Authorization: Bearer <JWT>
```

Приклад response:

```json
[
    {
        "id": "...",
        "player_id": "...",
        "item_id": "...",
        "quantity": 5,
        "created_at": "...",
        "updated_at": "..."
    }
]
```

---

## POST /api/inventory/items

Додає предмет до inventory.

```http
POST /api/inventory/items
Authorization: Bearer <JWT>
Content-Type: application/json
```

Request:

```json
{
    "item_id": "...",
    "quantity": 5
}
```

Якщо предмет вже є:

```text
old quantity + new quantity
```

Наприклад:

```text
5 + 3 = 8
```

---

## DELETE /api/inventory/items/{itemID}

Повністю видаляє предмет із inventory.

```http
DELETE /api/inventory/items/{itemID}
Authorization: Bearer <JWT>
```

При успішному видаленні:

```text
204 No Content
```

---

# 11. Взаємодія з іншими модулями

Inventory Module не існує ізольовано.

Основні залежності:

```text
Auth
 ↓
Player
 ↓
Inventory
 ↓
Item
```

### Auth

Auth визначає користувача через JWT.

```text
JWT → user_id
```

### Player

Через `user_id` Inventory знаходить Player.

```text
user_id → player_id
```

### Item

`item_id` посилається на Item Catalog.

```text
Inventory.item_id
        ↓
Items.id
```

---

# 12. Приклад повного flow

Гравець отримує 5 аптечок.

```text
Game Event
    ↓
Item ID = Medkit
Quantity = 5
    ↓
Inventory Service
    ↓
Inventory Repository
    ↓
PostgreSQL
```

У database:

```text
player_id | item_id | quantity
----------|---------|---------
Player A  | Medkit  | 5
```

Гравець отримує ще 3:

```text
5 + 3 = 8
```

У database:

```text
player_id | item_id | quantity
----------|---------|---------
Player A  | Medkit  | 8
```

---

# 13. Inventory vs Item

Це два різні поняття.

### Item

Описує предмет гри:

```text
AK-47
type: weapon
rarity: common
stackable: false
```

### Inventory

Описує володіння предметом:

```text
Player A
AK-47 × 1
```

Тобто:

```text
Item = що це за предмет

Inventory = скільки цього предмета має гравець
```

---

# 14. Поточний статус

Inventory Module реалізований і протестований.

Реалізовано:

* Inventory Model
* Repository
* Service
* Handler
* PostgreSQL migration
* Player → Inventory relation
* Item → Inventory relation
* Unique `(player_id, item_id)`
* Positive quantity constraint
* Add Item
* Remove Item
* Delete Item
* Get Inventory
* JWT protection
* автоматичне визначення Player через authenticated User
* документацію модуля

### Поточна структура

```text
internal/
├── inventory/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
└── database/
    └── migrations/
        └── 003_create_inventory_items.sql
```

Inventory Module готовий до використання іншими ігровими системами.

Наступні модулі можуть використовувати його для:

```text
Loot
  ↓
Inventory

Rewards
  ↓
Inventory

Extraction
  ↓
Inventory

Match
  ↓
Loot
  ↓
Inventory
```
