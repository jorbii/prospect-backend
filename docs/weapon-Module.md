# Weapon Module

## Зміст

1. Що таке Weapon Module
2. Архітектура
3. Model
4. Repository
5. Service
6. Handler
7. Database
8. Зв'язок з Item
9. API
10. Приклад Weapon
11. Взаємодія з іншими модулями
12. Поточний статус

---

# 1. Що таке Weapon Module

`Weapon Module` відповідає за характеристики зброї в грі.

Модуль зберігає:

* damage;
* fire rate;
* magazine size;
* reload time;
* range.

Weapon не відповідає за те, чи має конкретний гравець цю зброю.

За володіння предметами відповідає `Inventory Module`.

Тому:

```text
Item = що це за предмет

Weapon = характеристики зброї

Inventory = скільки предмета має гравець
```

---

# 2. Архітектура

Weapon Module використовує стандартну backend-архітектуру:

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

Відповідає за:

* HTTP requests;
* parsing UUID;
* HTTP status codes;
* JSON responses.

### Service

Відповідає за:

* business validation;
* перевірку характеристик зброї;
* взаємодію між Handler і Repository.

### Repository

Відповідає за:

* SQL queries;
* INSERT;
* SELECT;
* роботу з PostgreSQL.

### Model

Описує структуру Weapon у Go.

---

# 3. Model

Файл:

```text
internal/weapon/model.go
```

Основна структура:

```go
type Weapon struct {
    ID           uuid.UUID `json:"id"`
    ItemID       uuid.UUID `json:"item_id"`
    Damage       int       `json:"damage"`
    FireRate     int       `json:"fire_rate"`
    MagazineSize int       `json:"magazine_size"`
    ReloadTime   float64   `json:"reload_time"`
    Range        int       `json:"range"`
    CreatedAt    time.Time `json:"created_at"`
}
```

### Поля

| Поле           | Тип       | Опис                  |
| -------------- | --------- | --------------------- |
| `ID`           | UUID      | ID конфігурації зброї |
| `ItemID`       | UUID      | ID відповідного Item  |
| `Damage`       | int       | шкода за постріл      |
| `FireRate`     | int       | швидкість стрільби    |
| `MagazineSize` | int       | місткість магазину    |
| `ReloadTime`   | float64   | час перезарядки       |
| `Range`        | int       | дальність зброї       |
| `CreatedAt`    | time.Time | час створення         |

---

# 4. Repository

Файл:

```text
internal/weapon/repository.go
```

Repository містить чотири основні операції:

```text
Create
GetByID
GetByItemID
GetAll
```

## Create

Створює Weapon configuration у database.

Використовується для seed/admin/internal backend operations.

Public endpoint для створення зброї клієнтом не передбачений.

---

## GetByID

Отримує Weapon за його власним UUID:

```text
weapon_id
    ↓
weapons.id
    ↓
Weapon
```

---

## GetByItemID

Отримує Weapon через `item_id`.

Це важлива операція для зв'язку:

```text
Item
 ↓
item_id
 ↓
Weapon
```

---

## GetAll

Повертає всі Weapon configurations.

Предмети сортуються за `created_at`.

---

# 5. Service

Файл:

```text
internal/weapon/service.go
```

Service перевіряє характеристики перед записом у database.

### Item ID

`ItemID` не може бути порожнім UUID:

```text
ItemID != uuid.Nil
```

### Damage

```text
damage > 0
```

### Fire Rate

```text
fire_rate > 0
```

### Magazine Size

```text
magazine_size > 0
```

### Reload Time

```text
reload_time > 0
```

### Range

```text
range > 0
```

Таким чином некоректні характеристики не проходять через Service до Repository.

---

# 6. Handler

Файл:

```text
internal/weapon/handler.go
```

Handler реалізує read API для Weapon.

Доступні операції:

```text
GetAll
GetByID
GetByItemID
```

Handler також перевіряє UUID, отримані з URL.

Невалідний UUID:

```text
400 Bad Request
```

Weapon, якого не існує:

```text
404 Not Found
```

---

# 7. Database

Migration:

```text
internal/database/migrations/006_create_weapons.sql
```

Створюється таблиця:

```text
weapons
```

Структура:

```text
id
item_id
damage
fire_rate
magazine_size
reload_time
range
created_at
```

---

# 8. Зв'язок з Item

Weapon напряму пов'язаний з Item:

```sql
item_id UUID NOT NULL UNIQUE
REFERENCES items(id)
ON DELETE CASCADE
```

Це означає, що кожна Weapon configuration належить одному Item.

Наприклад:

```text
Item
├── id: ...
├── name: AK-47
└── type: weapon

Weapon
├── item_id: ...
├── damage: 32
├── fire_rate: 600
├── magazine_size: 30
├── reload_time: 2.4
└── range: 50
```

---

# 9. Unique Item

`item_id` має `UNIQUE` constraint.

Тому один Item може мати тільки одну Weapon configuration.

Правильно:

```text
AK-47
   ↓
Weapon
```

Неможливо:

```text
AK-47
 ├── Weapon configuration 1
 └── Weapon configuration 2
```

---

# 10. API

Weapon endpoints є read-only для клієнта.

## GET /api/weapons

Отримує всі Weapon configurations.

```http
GET /api/weapons
```

Response:

```json
[
    {
        "id": "...",
        "item_id": "...",
        "damage": 32,
        "fire_rate": 600,
        "magazine_size": 30,
        "reload_time": 2.4,
        "range": 50,
        "created_at": "..."
    }
]
```

---

## GET /api/weapons/{id}

Отримує Weapon за її ID.

```http
GET /api/weapons/{weaponID}
```

Можливі responses:

```text
200 OK
400 Bad Request
404 Not Found
```

---

## GET /api/weapons/item/{itemID}

Отримує Weapon configuration через Item ID.

```http
GET /api/weapons/item/{itemID}
```

Наприклад:

```text
Item ID
   ↓
Weapon configuration
```

Можливі responses:

```text
200 OK
400 Bad Request
404 Not Found
```

---

# 11. Приклад Weapon

Для `AK-47` використовується:

```text
Damage:        32
Fire Rate:     600
Magazine Size: 30
Reload Time:   2.4
Range:         50
```

У грі це може бути представлено:

```text
Item
AK-47
    ↓
Weapon
    ├── Damage: 32
    ├── Fire Rate: 600
    ├── Magazine: 30
    ├── Reload: 2.4s
    └── Range: 50
```

---

# 12. Взаємодія з іншими модулями

Weapon Module пов'язаний з кількома системами.

## Item Module

Item визначає сам предмет:

```text
Item
├── name
├── type
├── rarity
└── stackable
```

Weapon визначає його характеристики:

```text
Weapon
├── damage
├── fire_rate
├── magazine_size
├── reload_time
└── range
```

---

## Inventory Module

Inventory зберігає предмет гравця:

```text
Player
 ↓
Inventory
 ↓
item_id
 ↓
Item
 ↓
Weapon
```

Наприклад:

```text
Player A
    ↓
Inventory
    ↓
AK-47 × 1
    ↓
Item: AK-47
    ↓
Weapon:
Damage = 32
Fire Rate = 600
Magazine = 30
```

---

## Combat Module

У майбутньому `Combat Module` використовуватиме Weapon для розрахунку бойових параметрів:

```text
Weapon
 ↓
Damage
 ↓
Combat
 ↓
Target HP
```

---

# 13. Поточний статус

Weapon Module реалізований і протестований.

Реалізовано:

* Weapon Model
* Weapon Repository
* Weapon Service
* Weapon Handler
* PostgreSQL migration
* Item → Weapon relation
* `UNIQUE(item_id)`
* database validation constraints
* Weapon creation на рівні Service/Repository
* Get All
* Get By ID
* Get By Item ID
* UUID validation
* `400 Bad Request`
* `404 Not Found`
* API integration
* документацію модуля

Поточна структура:

```text
internal/
├── weapon/
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
└── database/
    └── migrations/
        └── 006_create_weapons.sql
```

Weapon Module готовий для використання Combat, Loot та інших ігрових систем.
