# Item Module

## 1. Що таке Item Module

Item Module відповідає за каталог усіх предметів, які існують у грі Prospect.

Предмет є базовою сутністю, яка використовується іншими модулями:

* Inventory
* Weapon
* Loot
* Rewards
* Match

Item Module зберігає загальні характеристики предмета, але не відповідає за його використання в грі.

---

## 2. Архітектура

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

### Model

`model.go` описує структуру Item.

```text
Item
├── ID
├── Name
├── Type
├── Stackable
├── MaxStack
├── Rarity
└── CreatedAt
```

### Repository

`repository.go` відповідає за роботу з PostgreSQL.

Основні операції:

* `Create`
* `GetByID`
* `GetAll`

Repository містить SQL-запити та не відповідає за бізнес-правила.

### Service

`service.go` містить бізнес-логіку Item Module.

Відповідає за:

* перевірку назви;
* перевірку типу;
* перевірку `max_stack`;
* перевірку сумісності `stackable` і `max_stack`;
* встановлення `common` як стандартної рідкості.

### Handler

`handler.go` відповідає за HTTP API.

Доступні клієнту операції:

```text
GET /api/items
GET /api/items/{id}
```

Метод `Create` існує в Service/Repository для внутрішнього використання, seed або майбутньої адміністративної системи, але не відкритий через публічний API.

---

## 3. Database

Migration:

```text
internal/database/migrations/004_create_items.sql
```

Створює таблицю `items`.

Основні поля:

| Поле         | Тип         | Опис                       |
| ------------ | ----------- | -------------------------- |
| `id`         | UUID        | Унікальний ID предмета     |
| `name`       | TEXT        | Назва предмета             |
| `type`       | TEXT        | Тип предмета               |
| `stackable`  | BOOLEAN     | Чи можна складати предмети |
| `max_stack`  | INTEGER     | Максимальний розмір stack  |
| `rarity`     | TEXT        | Рідкість предмета          |
| `created_at` | TIMESTAMPTZ | Час створення              |

Migration:

```text
internal/database/migrations/005_add_inventory_item_fk.sql
```

Додає зв'язок:

```text
inventory_items.item_id
        │
        ▼
     items.id
```

Це гарантує, що Inventory не може посилатися на неіснуючий Item.

---

## 4. Item Types

Приклади типів:

```text
weapon
ammo
medical
armor
key
```

Тип визначає загальну категорію предмета.

Специфічні характеристики зброї не зберігаються в Item Module. Вони будуть реалізовані в Weapon Module.

---

## 5. Stacking

Для предметів, які можна складати:

```text
stackable = true
max_stack > 1
```

Наприклад:

```text
5.56 Ammo
stackable = true
max_stack = 120
```

Для предметів, які не можна складати:

```text
stackable = false
max_stack = 1
```

Наприклад:

```text
AK-47
stackable = false
max_stack = 1
```

---

## 6. Rarity

Поточна система підтримує текстове поле `rarity`.

Приклади:

```text
common
rare
epic
legendary
```

Рідкість може використовуватися майбутніми Loot і Rewards модулями.

---

## 7. API

### Get All Items

```http
GET /api/items
```

Повертає список усіх предметів.

Приклад:

```json
[
  {
    "id": "...",
    "name": "AK-47",
    "type": "weapon",
    "stackable": false,
    "max_stack": 1,
    "rarity": "common",
    "created_at": "..."
  }
]
```

### Get Item By ID

```http
GET /api/items/{id}
```

Повертає конкретний предмет.

Можливі відповіді:

```text
200 OK
400 Bad Request
404 Not Found
```

---

## 8. Приклад даних

Поточний тестовий каталог:

```text
AK-47
5.56 Ammo
Medkit
Armor Vest
Keycard
```

---

## 9. Взаємодія з іншими модулями

### Inventory

Inventory зберігає:

```text
player_id
item_id
quantity
```

При цьому характеристики предмета знаходяться в `items`.

```text
Player
   │
   ▼
Inventory
   │
   └── item_id ──► Item
```

### Weapon

Weapon Module буде використовувати Item як базову сутність.

```text
Item
  │
  ▼
Weapon
├── Damage
├── FireRate
├── Magazine
├── AmmoType
└── ...
```

### Loot

Loot Module визначатиме, які Items можуть з'являтися на карті.

```text
Loot
  │
  └── item_id ──► Item
```

---

## 10. Поточний стан

Item Module реалізовано та протестовано.

```text
✅ Model
✅ Repository
✅ Service
✅ Handler
✅ PostgreSQL migration
✅ Inventory foreign key
✅ GET /api/items
✅ GET /api/items/{id}
✅ Validation
```

Наступний модуль:

```text
Weapon Module
```
