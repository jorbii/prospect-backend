# Player Module

## Зміст

1. [Що таке Player Module](#s1)
2. [Архітектура Player Module](#s2)
   - 2.1 [Model](#s2-1)
   - 2.2 [Service](#s2-2)
   - 2.3 [Repository](#s2-3)
3. [Player Creation](#s3)
   - 3.1 [Створення Player](#s3-1)
   - 3.2 [Initial Values](#s3-2)
   - 3.3 [Transaction](#s3-3)
4. [Player Data](#s4)
   - 4.1 [XP](#s4-1)
   - 4.2 [Level](#s4-2)
   - 4.3 [Cash](#s4-3)
   - 4.4 [Premium Currency](#s4-4)
5. [Player Retrieval](#s5)
   - 5.1 [Get By User ID](#s5-1)
   - 5.2 [JWT → User ID → Player](#s5-2)
6. [XP and Level System](#s6)
   - 6.1 [Adding XP](#s6-1)
   - 6.2 [Level Calculation](#s6-2)
   - 6.3 [Atomic Progression Update](#s6-3)
7. [Cash System](#s7)
   - 7.1 [Adding Cash](#s7-1)
   - 7.2 [Removing Cash](#s7-2)
   - 7.3 [Insufficient Cash](#s7-3)
8. [Premium Currency](#s8)
   - 8.1 [Adding Premium Currency](#s8-1)
   - 8.2 [Removing Premium Currency](#s8-2)
   - 8.3 [Insufficient Premium Currency](#s8-3)
9. [Atomic Operations](#s9)
10. [Database](#s10)
    - 10.1 [Players Table](#s10-1)
    - 10.2 [Constraints](#s10-2)
11. [Error Handling](#s11)
12. [Project Files](#s12)
13. [Complete Player Flow](#s13)
14. [Current Status](#s14)
15. [Next Module](#s15)

---

<a id="s1"></a>
## 1. Що таке Player Module

Player Module відповідає за ігрові дані користувача.

Він відділений від Auth Module.

Auth Module відповідає за:

```text
Account
    ↓
User
    ├── username
    ├── email
    └── password

Player Module відповідає за:

Game Data
    ↓
Player
    ├── XP
    ├── Level
    ├── Cash
    └── Premium Currency

Головна ідея:

User
  │
  │ userID
  ▼
Player
  │
  ├── progression
  ├── currency
  └── game data

Один User має одного Player.

<a id="s2"></a>

2. Архітектура Player Module

Основна схема:

HTTP Handler
      ↓
Player Service
      ↓
Player Repository
      ↓
PostgreSQL

Player Module не займається:

authentication;
password hashing;
JWT generation;
JWT validation.

Authentication виконується Auth Middleware.

Player Module отримує вже визначений userID.

<a id="s2-1"></a>

2.1 Model

Основна структура Player:

type Player struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    XP              int64
    Level           int
    Cash            int64
    PremiumCurrency int64
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
Поля

ID

Унікальний ID ігрового персонажа.

UUID

UserID

ID користувача з таблиці users.

players.user_id → users.id

XP

Досвід гравця.

int64

Level

Поточний рівень гравця.

int

Cash

Основна ігрова валюта.

int64

PremiumCurrency

Premium валюта.

int64

CreatedAt

Час створення Player.

UpdatedAt

Час останньої зміни Player.

<a id="s2-2"></a>

2.2 Service

Service містить бізнес-логіку Player Module.

Основні операції:

CreatePlayerTx
GetByUserID


AddXP


AddCash
RemoveCash


AddPremiumCurrency
RemovePremiumCurrency

Service не виконує SQL напряму.

Він використовує Repository.

<a id="s2-3"></a>

2.3 Repository

Repository відповідає за PostgreSQL.

Основні операції:

Create
CreateTx
GetByUserID
GetByID


AddXP


AddCash
RemoveCash


AddPremiumCurrency
RemovePremiumCurrency

Repository не повинен містити HTTP logic.

<a id="s3"></a>

3. Player Creation

Player автоматично створюється під час registration.

Flow:

POST /api/auth/register
        ↓
Create User
        ↓
Create Player

Player не створюється клієнтом окремим запитом.

<a id="s3-1"></a>

3.1 Створення Player

Auth Service створює User і викликає:

CreatePlayerTx(
    ctx,
    tx,
    user.ID,
)

Player Service створює новий Player:

player := Player{
    ID:              uuid.New(),
    UserID:          userID,
    XP:              0,
    Level:           1,
    Cash:            0,
    PremiumCurrency: 0,
    CreatedAt:       now,
    UpdatedAt:       now,
}

<a id="s3-2"></a>

3.2 Initial Values

Новий Player отримує:

XP                = 0
Level             = 1
Cash              = 0
Premium Currency  = 0

Таким чином кожен новий користувач починає гру з однакового стану.

<a id="s3-3"></a>

3.3 Transaction

Створення User і Player виконується в одній PostgreSQL transaction.

Схема:

BEGIN
  │
  ├── INSERT User
  │
  ├── INSERT Player
  │
COMMIT

Якщо створення Player завершується помилкою:

ROLLBACK

У результаті User також не повинен залишитися без Player.

Це гарантує цілісність між:

users
  ↓
players

<a id="s4"></a>

4. Player Data

Player зараз містить базову progression та currency систему.

Player
│
├── XP
├── Level
├── Cash
└── Premium Currency

<a id="s4-1"></a>

4.1 XP

XP — досвід гравця.

XP збільшується через:

AddXP(
    ctx,
    playerID,
    amount,
)

XP не може збільшуватися на 0 або негативне значення.

Наприклад:

AddXP(500)


0 XP
 ↓
500 XP

<a id="s4-2"></a>

4.2 Level

Level визначає progression гравця.

Поточна система:

XP	Level
0–499	1
500–999	2
1000–1999	3
2000–2999	4
3000–3999	5
4000–4999	6
5000–6499	7
6500–7999	8
8000–9999	9
10000+	10

Максимальний Level у поточній реалізації:

10

Ця система може бути розширена в майбутньому.

<a id="s4-3"></a>

4.3 Cash

Cash — основна ігрова валюта.

Додавання:

AddCash(
    ctx,
    playerID,
    amount,
)

Витрачання:

RemoveCash(
    ctx,
    playerID,
    amount,
)

Amount повинен бути позитивним:

amount > 0

<a id="s4-4"></a>

4.4 Premium Currency

Premium Currency — окрема валюта для premium-механік.

Додавання:

AddPremiumCurrency(
    ctx,
    playerID,
    amount,
)

Витрачання:

RemovePremiumCurrency(
    ctx,
    playerID,
    amount,
)

Premium Currency зберігається окремо від Cash.

Це дозволяє в майбутньому використовувати її для:

Premium Shop
Premium Items
Cosmetics
Battle Pass

<a id="s5"></a>

5. Player Retrieval

Player можна отримати через userID.

Основний метод:

GetByUserID(
    ctx,
    userID,
)

<a id="s5-1"></a>

5.1 Get By User ID

Repository виконує:

SELECT
    id,
    user_id,
    xp,
    level,
    cash,
    premium_currency,
    created_at,
    updated_at
FROM players
WHERE user_id = $1

Таким чином Player знаходиться через User ID.

<a id="s5-2"></a>

5.2 JWT → User ID → Player

Захищений endpoint використовує JWT.

Flow:

Client
   ↓
Authorization: Bearer JWT
   ↓
Auth Middleware
   ↓
JWT verification
   ↓
userID
   ↓
Context
   ↓
Player Service
   ↓
GetByUserID
   ↓
Player

JWT не містить Player ID.

JWT містить:

sub = userID

Тому Player Module спочатку отримує userID.

<a id="s6"></a>

6. XP and Level System

XP і Level є частиною progression системи.

<a id="s6-1"></a>

6.1 Adding XP

Service перевіряє amount:

if amount <= 0 {
    return Player{}, errors.New(
        "xp amount must be positive",
    )
}

Після цього Repository атомарно додає XP.

<a id="s6-2"></a>

6.2 Level Calculation

Level визначається на основі нового XP.

Наприклад:

XP = 450
    ↓
Level 1


+100 XP
    ↓
XP = 550
    ↓
Level 2

При досягненні нового threshold Level автоматично збільшується.

<a id="s6-3"></a>

6.3 Atomic Progression Update

XP і Level змінюються однією SQL операцією.

Приклад:

UPDATE players
SET
    xp = xp + $1,
    level = CASE
        WHEN xp + $1 >= 10000 THEN 10
        WHEN xp + $1 >= 8000 THEN 9
        WHEN xp + $1 >= 6500 THEN 8
        WHEN xp + $1 >= 5000 THEN 7
        WHEN xp + $1 >= 4000 THEN 6
        WHEN xp + $1 >= 3000 THEN 5
        WHEN xp + $1 >= 2000 THEN 4
        WHEN xp + $1 >= 1000 THEN 3
        WHEN xp + $1 >= 500 THEN 2
        ELSE 1
    END,
    updated_at = NOW()
WHERE id = $2
RETURNING ...

Це означає, що XP і Level змінюються атомарно.

<a id="s7"></a>

7. Cash System

Cash використовується як основна ігрова валюта.

<a id="s7-1"></a>

7.1 Adding Cash

Cash додається через:

AddCash(
    ctx,
    playerID,
    amount,
)

SQL:

UPDATE players
SET
    cash = cash + $1,
    updated_at = NOW()
WHERE id = $2

Amount повинен бути:

> 0

<a id="s7-2"></a>

7.2 Removing Cash

Cash витрачається через:

RemoveCash(
    ctx,
    playerID,
    amount,
)

SQL:

UPDATE players
SET
    cash = cash - $1,
    updated_at = NOW()
WHERE id = $2
  AND cash >= $1

<a id="s7-3"></a>

7.3 Insufficient Cash

Якщо:

cash < amount

UPDATE не змінює жодного рядка.

Repository повертає:

false

Service перетворює це на:

insufficient cash

Таким чином Cash не може стати негативним через звичайну операцію RemoveCash.

<a id="s8"></a>

8. Premium Currency

Premium Currency працює аналогічно Cash.

<a id="s8-1"></a>

8.1 Adding Premium Currency
AddPremiumCurrency(
    ctx,
    playerID,
    amount,
)

SQL використовує:

premium_currency = premium_currency + $1

<a id="s8-2"></a>

8.2 Removing Premium Currency
RemovePremiumCurrency(
    ctx,
    playerID,
    amount,
)

SQL:

UPDATE players
SET
    premium_currency = premium_currency - $1,
    updated_at = NOW()
WHERE id = $2
  AND premium_currency >= $1

<a id="s8-3"></a>

8.3 Insufficient Premium Currency

Якщо premium balance недостатній:

premium_currency < amount

операція не виконується.

Service повертає:

insufficient premium currency

<a id="s9"></a>

9. Atomic Operations

Currency operations використовують атомарні SQL UPDATE.

Неправильний підхід:

SELECT balance
      ↓
balance + amount
      ↓
UPDATE

Такий підхід може створити race condition.

Правильний підхід:

UPDATE players
SET cash = cash + $1
WHERE id = $2

А для списання:

UPDATE players
SET cash = cash - $1
WHERE id = $2
  AND cash >= $1

PostgreSQL гарантує атомарність UPDATE.

Це особливо важливо для ігрової економіки, де одночасно можуть відбуватися:

Reward
Purchase
Loot
Quest reward
Match reward

<a id="s10"></a>

10. Database

Player Module використовує PostgreSQL.

Основна таблиця:

players

<a id="s10-1"></a>

10.1 Players Table

Поточна структура:

Column	Type
id	uuid
user_id	uuid
xp	bigint
level	integer
cash	bigint
premium_currency	bigint
created_at	timestamptz
updated_at	timestamptz

<a id="s10-2"></a>

10.2 Constraints

players.id є primary key.

players.user_id повинен посилатися на:

users.id

Один User має один Player.

Тому user_id повинен бути:

UNIQUE

Currency і XP не повинні приймати негативні значення.

Рекомендовані database constraints:

CHECK (xp >= 0)


CHECK (level >= 1)


CHECK (cash >= 0)


CHECK (premium_currency >= 0)

Ці constraints є додатковим рівнем захисту.

Навіть якщо помилка з'явиться в application code, база не дозволить записати некоректний баланс.

<a id="s11"></a>

11. Error Handling

Player Service використовує бізнес-помилки.

Invalid amount
xp amount must be positive
cash amount must be positive
premium currency amount must be positive

Це означає, що клієнт або внутрішня система передала:

0

або негативне значення.

Insufficient Cash
insufficient cash

Гравець не має достатньо Cash.

Insufficient Premium Currency
insufficient premium currency

Гравець не має достатньо Premium Currency.

Database Error

Будь-яка неочікувана помилка PostgreSQL повертається вище по application stack.

<a id="s12"></a>

12. Project Files
model.go

Містить:

Player
PlayerResponse
service.go

Містить:

CreatePlayerTx
GetByUserID
AddXP


AddCash
RemoveCash


AddPremiumCurrency
RemovePremiumCurrency
repository.go

Містить database operations:

Create
CreateTx
GetByID
GetByUserID


AddXP


AddCash
RemoveCash


AddPremiumCurrency
RemovePremiumCurrency

<a id="s13"></a>

13. Complete Player Flow
Registration
Client
   ↓
POST /api/auth/register
   ↓
Auth Service
   ↓
Create User
   ↓
Create Player
   ↓
PostgreSQL Transaction
   ↓
COMMIT
Get Player
Client
   ↓
Authorization: Bearer JWT
   ↓
Auth Middleware
   ↓
JWT Verification
   ↓
User ID
   ↓
Context
   ↓
Player Service
   ↓
GetByUserID
   ↓
PostgreSQL
   ↓
Player
Add XP
Game System
   ↓
PlayerService.AddXP()
   ↓
PlayerRepository.AddXP()
   ↓
PostgreSQL
   ↓
XP + amount
   ↓
Calculate Level
   ↓
Player updated
Add Cash
Game System
   ↓
AddCash()
   ↓
Repository
   ↓
cash = cash + amount
Spend Cash
Game System
   ↓
RemoveCash()
   ↓
Repository
   ↓
cash >= amount?
   │
   ├── YES → subtract
   │
   └── NO  → insufficient cash

<a id="s14"></a>

14. Current Status

Player Module реалізує:

Player Creation
автоматичне створення Player;
UUID;
початковий Level 1;
початковий XP 0;
початковий Cash 0;
початкова Premium Currency 0;
transaction разом зі створенням User.
Player Retrieval
Get Player by User ID;
інтеграція з JWT userID;
context-based user identification.
Progression
XP;
Level;
automatic Level calculation;
atomic XP + Level update.
Economy
Cash;
Add Cash;
Remove Cash;
insufficient balance protection;
atomic currency operations.
Premium Economy
Premium Currency;
Add Premium Currency;
Remove Premium Currency;
insufficient balance protection;
atomic currency operations.

Player Module готовий як базовий v1 фундамент для ігрових систем Prospect.

<a id="s15"></a>

15. Next Module

Наступні ігрові системи можуть використовувати Player Module:

Player
  │
  ├── Match
  │
  ├── Inventory
  │
  ├── Loot
  │
  ├── Weapons
  │
  ├── Quests
  │
  ├── Rewards
  │
  └── Shop

Наприклад:

Match
  ↓
Player survives
  ↓
Reward
  ├── XP
  └── Cash
  ↓
PlayerService

Або:

Shop
  ↓
Purchase
  ↓
PlayerService.RemoveCash()
  ↓
Inventory

Player Module таким чином стає базовим шаром для progression та economy всієї гри Prospect.