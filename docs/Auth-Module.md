# Auth Module

## Зміст

1. [Що таке Auth Module](#s1)
2. [Архітектура Auth Module](#s2)
   - 2.1 [Handler](#s2-1)
   - 2.2 [Service](#s2-2)
   - 2.3 [Repository](#s2-3)
   - 2.4 [TokenService](#s2-4)
   - 2.5 [Middleware](#s2-5)
3. [Registration](#s3)
   - 3.1 [Request](#s3-1)
   - 3.2 [Validation](#s3-2)
   - 3.3 [Password Hashing](#s3-3)
   - 3.4 [Creating User](#s3-4)
   - 3.5 [Saving User](#s3-5)
   - 3.6 [Duplicate Users](#s3-6)
   - 3.7 [Successful Registration](#s3-7)
4. [Login](#s4)
   - 4.1 [Request](#s4-1)
   - 4.2 [Finding User](#s4-2)
   - 4.3 [Password Verification](#s4-3)
   - 4.4 [JWT Generation](#s4-4)
5. [JWT](#s5)
   - 5.1 [JWT Claims](#s5-1)
   - 5.2 [JWT Secret](#s5-2)
   - 5.3 [HS256](#s5-3)
   - 5.4 [Token Expiration](#s5-4)
6. [Authentication Middleware](#s6)
   - 6.1 [Authorization Header](#s6-1)
   - 6.2 [JWT Verification](#s6-2)
   - 6.3 [Getting User ID](#s6-3)
   - 6.4 [Context](#s6-4)
   - 6.5 [Protected Endpoint](#s6-5)
7. [Database](#s7)
   - 7.1 [Users Table](#s7-1)
   - 7.2 [Constraints](#s7-2)
8. [Error Handling](#s8)
   - 8.1 [HTTP 400](#s8-1)
   - 8.2 [HTTP 401](#s8-2)
   - 8.3 [HTTP 409](#s8-3)
   - 8.4 [HTTP 500](#s8-4)
9. [Project Files](#s9)
10. [Complete Auth Flow](#s10)
11. [Current Status](#s11)
12. [Next Module — Player](#s12)

---

<a id="s1"></a>
## 1. Що таке Auth Module

Auth Module відповідає за акаунт користувача та його автентифікацію.

Він виконує такі задачі:

- реєстрація користувача;
- перевірка даних реєстрації;
- хешування пароля;
- збереження користувача в PostgreSQL;
- логін;
- перевірка пароля;
- створення JWT access token;
- перевірка JWT;
- визначення, який саме користувач робить запит.

**Головна ідея**

Користувач один раз логіниться:

```
Email + Password
      ↓
    Login
      ↓
    JWT Token
```

Після цього клієнт відправляє JWT з наступними запитами:

```
Request
   +
Bearer Token
   ↓
Server
   ↓
"Я знаю, який це користувач"
```

Тобто пароль не потрібно відправляти з кожним запитом.

---

<a id="s2"></a>
## 2. Архітектура Auth Module

Auth Module розділений на декілька шарів.

Основна схема:

```
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

Для JWT є окремий потік:

```
Login
  ↓
TokenService
  ↓
JWT
  ↓
Client
  ↓
Auth Middleware
  ↓
Protected Handler
```

Таке розділення потрібне для того, щоб кожна частина системи мала одну чітку відповідальність.

<a id="s2-1"></a>
### 2.1 Handler

Handler працює з HTTP.

Він:

- отримує HTTP request;
- читає JSON;
- викликає Service;
- обробляє помилки;
- повертає HTTP response.

Наприклад:

```
POST /api/auth/login
        ↓
LoginHandler
        ↓
LoginUser()
```

Handler не повинен містити складну бізнес-логіку.

<a id="s2-2"></a>
### 2.2 Service

Service містить основну бізнес-логіку Auth.

Наприклад, під час реєстрації Service:

- перевіряє дані;
- хешує пароль;
- створює User;
- передає User у Repository.

Service не повинен сам виконувати SQL.

<a id="s2-3"></a>
### 2.3 Repository

Repository відповідає за роботу з PostgreSQL.

Наприклад:

```
CreateUser()
GetUserByEmail()
```

Repository знає SQL, але не повинен знати про HTTP.

Наприклад, Repository не повинен вирішувати:

> "Який HTTP status повернути?"

Це відповідальність Handler.

<a id="s2-4"></a>
### 2.4 TokenService

TokenService відповідає за створення JWT.

Він отримує ID користувача:

```
User ID
   ↓
TokenService
   ↓
JWT
```

TokenService не займається PostgreSQL і не працює з HTTP.

<a id="s2-5"></a>
### 2.5 Middleware

Middleware перевіряє JWT перед доступом до захищеного endpoint.

Наприклад:

```
GET /api/player/profile
        ↓
Auth Middleware
        ↓
JWT valid?
        ↓
YES → Handler
NO  → 401
```

Middleware також дістає userID з JWT і передає його далі через context.

---

<a id="s3"></a>
## 3. Registration

Endpoint:

```
POST /api/auth/register
```

Реєстрація створює нового користувача в базі.

<a id="s3-1"></a>
### 3.1 Request

Клієнт відправляє:

```json
{
  "username": "joligo",
  "email": "joligo@test.com",
  "password": "12345678"
}
```

Handler перетворює JSON у Go-структуру:

```go
var req RegisterRequest

err := json.NewDecoder(r.Body).Decode(&req)
```

Після цього:

```
req.Username = "joligo"
req.Email = "joligo@test.com"
req.Password = "12345678"
```

<a id="s3-2"></a>
### 3.2 Validation

Перед створенням користувача перевіряємо:

- username не порожній;
- email не порожній;
- password не порожній;
- password має мінімум 8 символів.

Наприклад:

```go
if len(req.Password) < 8 {
    return errors.New("password is too short")
}
```

Якщо validation не пройшла:

```
HTTP 400 Bad Request
```

<a id="s3-3"></a>
### 3.3 Password Hashing

Пароль не можна зберігати у відкритому вигляді.

Неправильно:

```
password = 12345678
```

Якщо база даних буде скомпрометована, пароль стане відомий.

Тому використовується bcrypt.

```go
passwordHash, err := bcrypt.GenerateFromPassword(
    []byte(req.Password),
    bcrypt.DefaultCost,
)
```

Результат приблизно такий:

```
12345678
    ↓
  bcrypt
    ↓
$2a$10$...
```

У PostgreSQL зберігається тільки hash.

Оригінальний пароль у базі не зберігається.

<a id="s3-4"></a>
### 3.4 Creating User

Після створення password hash створюється User:

```go
user := User{
    ID:           uuid.New(),
    Username:     req.Username,
    Email:        req.Email,
    PasswordHash: string(passwordHash),
    CreatedAt:    time.Now(),
}
```

Для кожного користувача генерується унікальний UUID:

```
24b03630-de9e-4a3f-b136-3cf4280ca091
```

Цей ID буде головним ідентифікатором користувача на backend.

<a id="s3-5"></a>
### 3.5 Saving User

Repository виконує SQL:

```sql
INSERT INTO users (
    id,
    username,
    email,
    password_hash,
    created_at
)
VALUES ($1, $2, $3, $4, $5)
```

Після успішного виконання користувач зберігається в PostgreSQL.

<a id="s3-6"></a>
### 3.6 Duplicate Users

Username та email мають бути унікальними.

Якщо користувач уже існує, PostgreSQL повертає:

```
SQLSTATE 23505
```

Repository перевіряє constraint:

```
users_username_key
users_email_key
```

Потім database error перетворюється на зрозумілу для Service помилку:

```
ErrUsernameExists
ErrEmailExists
```

Handler повертає:

```
409 Conflict
```

Наприклад:

```
username already exists
```

<a id="s3-7"></a>
### 3.7 Successful Registration

Якщо все пройшло успішно:

```json
{
  "message": "registration successful",
  "username": "joligo"
}
```

Пароль у response не повертається.

---

<a id="s4"></a>
## 4. Login

Login потрібен для перевірки, що користувач знає правильний email і пароль.

Endpoint:

```
POST /api/auth/login
```

<a id="s4-1"></a>
### 4.1 Request

Клієнт відправляє:

```json
{
  "email": "joligo@test.com",
  "password": "12345678"
}
```

Повний flow:

```
Client
  ↓
LoginHandler
  ↓
LoginUser
  ↓
GetUserByEmail
  ↓
PostgreSQL
  ↓
bcrypt.CompareHashAndPassword
  ↓
User.ID
  ↓
TokenService
  ↓
JWT
  ↓
Client
```

<a id="s4-2"></a>
### 4.2 Finding User

Repository шукає користувача по email:

```sql
SELECT
    id,
    username,
    email,
    password_hash,
    created_at,
    last_login_at
FROM users
WHERE email = $1
```

Якщо користувача немає:

```
ErrInvalidCredentials
```

Handler повертає:

```
401 Unauthorized
```

Ми спеціально не говоримо:

```
email does not exist
```

або:

```
password is wrong
```

Замість цього повертаємо одну загальну помилку:

```
invalid email or password
```

Це не дає атакуючому легко визначити, чи існує конкретний email у системі.

<a id="s4-3"></a>
### 4.3 Password Verification

Під час login ми не порівнюємо password hash як звичайні строки.

Використовується:

```go
bcrypt.CompareHashAndPassword(
    []byte(user.PasswordHash),
    []byte(req.Password),
)
```

Bcrypt перевіряє, чи відповідає введений пароль hash із бази.

Якщо пароль неправильний:

```
ErrInvalidCredentials
```

<a id="s4-4"></a>
### 4.4 JWT Generation

Якщо email і password правильні:

```
User authenticated
        ↓
TokenService
        ↓
JWT Access Token
```

Клієнт отримує:

```json
{
  "access_token": "eyJ..."
}
```

Цей token клієнт використовує для наступних захищених запитів.

---

<a id="s5"></a>
## 5. JWT

JWT — це токен, який сервер використовує для підтвердження того, що користувач уже пройшов login.

JWT має приблизно такий формат:

```
xxxxx.yyyyy.zzzzz
```

Три частини:

```
Header.Payload.Signature
```

<a id="s5-1"></a>
### 5.1 JWT Claims

Наш JWT містить основні claims:

```
sub
iat
exp
```

**sub**

sub означає Subject.

У нашому випадку це ID користувача:

```
sub:
24b03630-de9e-4a3f-b136-3cf4280ca091
```

Це найважливіше поле для нашого backend.

**iat**

iat означає:

```
Issued At
```

Тобто час створення token.

**exp**

exp означає:

```
Expiration Time
```

Тобто час, після якого token перестає бути дійсним.

У нашій реалізації token живе 24 години.

<a id="s5-2"></a>
### 5.2 JWT Secret

JWT підписується секретним ключем:

```
JWT_SECRET=...
```

Secret зберігається у `.env`.

Він не повинен бути захардкоджений у source code.

Під час запуску:

```go
os.Getenv("JWT_SECRET")
```

отримує secret.

Після цього створюється TokenService.

<a id="s5-3"></a>
### 5.3 HS256

Для JWT використовується:

```
HS256
```

Це HMAC SHA-256.

Спрощено процес виглядає так:

```
JWT + Secret
     ↓
 Signature
```

При перевірці:

```
JWT + Secret
     ↓
Verify Signature
```

Якщо signature неправильна — token відхиляється.

Для HS256 secret під час перевірки передається як `[]byte`.

<a id="s5-4"></a>
### 5.4 Token Expiration

JWT має expiration time.

Після завершення exp token більше не повинен прийматися сервером.

Це означає:

```
Token created
      ↓
   24 hours
      ↓
Token expired
```

Користувачу після цього потрібно отримати новий token через login або майбутню refresh-token систему.

---

<a id="s6"></a>
## 6. Authentication Middleware

Middleware захищає endpoints, які потребують авторизації.

Клієнт відправляє:

```
Authorization: Bearer eyJ...
```

Повний flow:

```
HTTP Request
     ↓
Authorization Header
     ↓
Bearer Token
     ↓
JWT Parse
     ↓
Check Algorithm
     ↓
Check Signature
     ↓
Check Token Validity
     ↓
Read Claims
     ↓
Get User ID
     ↓
Context
     ↓
Handler
```

<a id="s6-1"></a>
### 6.1 Authorization Header

Middleware отримує header:

```go
authHeader := r.Header.Get("Authorization")
```

Приклад:

```
Bearer eyJhbGciOiJIUzI1NiIs...
```

Слово Bearer означає, що після нього передається access token.

Якщо header відсутній:

```
401 Unauthorized
authorization required
```

Якщо формат неправильний:

```
401 Unauthorized
invalid authorization header
```

<a id="s6-2"></a>
### 6.2 JWT Verification

Middleware використовує:

```go
jwt.Parse(...)
```

Під час перевірки:

- перевіряється структура JWT;
- перевіряється алгоритм;
- перевіряється signature;
- перевіряється validity token;
- перевіряються claims, включно з expiration.

Ми очікуємо:

```
HS256
```

Якщо використовується інший алгоритм або signature неправильна:

```
401 Unauthorized
invalid token
```

Secret для перевірки передається як:

```go
[]byte(s.secretKey)
```

<a id="s6-3"></a>
### 6.3 Getting User ID

Після успішної перевірки JWT middleware отримує claims.

З них беремо:

```go
claims["sub"]
```

Наприклад:

```
24b03630-de9e-4a3f-b136-3cf4280ca091
```

Це userID.

Тепер сервер знає, який користувач робить запит.

<a id="s6-4"></a>
### 6.4 Context

User ID передається далі через Go context.

Спрощено:

```go
ctx := context.WithValue(
    r.Context(),
    userIDKey,
    userID,
)

r = r.WithContext(ctx)
```

Після цього наступний Handler може отримати user ID:

```go
userID := r.Context().Value(userIDKey)
```

Handler не повинен сам розбирати JWT.

Цим займається Middleware.

<a id="s6-5"></a>
### 6.5 Protected Endpoint

Наприклад:

```
GET /api/auth/test
```

Без token:

```
401 Unauthorized
```

З правильним token:

```
authenticated user:
24b03630-de9e-4a3f-b136-3cf4280ca091
```

Це підтверджує, що весь authentication flow працює:

```
JWT
 ↓
Middleware
 ↓
User ID
 ↓
Context
 ↓
Handler
```

---

<a id="s7"></a>
## 7. Database

Auth Module використовує PostgreSQL.

Основна таблиця:

```
users
```

<a id="s7-1"></a>
### 7.1 Users Table

Поточна структура:

| Column | Type |
|---|---|
| id | uuid |
| username | varchar(32) |
| email | varchar(255) |
| password_hash | text |
| created_at | timestamptz |
| last_login_at | timestamptz |

**id** — Унікальний UUID користувача.

**username** — Ім'я користувача. Максимум 32 символи.

**email** — Email користувача. Максимум 255 символів.

**password_hash** — Bcrypt hash пароля. Сам пароль тут не зберігається.

**created_at** — Дата створення акаунта.

**last_login_at** — Дата останнього login.

<a id="s7-2"></a>
### 7.2 Constraints

Username та email мають бути унікальними.

Основні constraints:

```
users_username_key
users_email_key
```

Це означає, що два користувачі не можуть мати однаковий username або email.

Також id, username, email, password_hash та created_at не можуть бути NULL.

---

<a id="s8"></a>
## 8. Error Handling

Auth Module використовує HTTP status codes для повідомлення клієнту про результат операції.

<a id="s8-1"></a>
### 8.1 HTTP 400

400 Bad Request означає, що клієнт відправив неправильний request.

Наприклад:

```
email is required
password is required
password is too short
invalid request
```

<a id="s8-2"></a>
### 8.2 HTTP 401

401 Unauthorized означає, що authentication не пройшла.

Наприклад:

```
authorization required
invalid authorization header
invalid token
invalid email or password
```

<a id="s8-3"></a>
### 8.3 HTTP 409

409 Conflict використовується, коли дані конфліктують з уже існуючими даними.

Наприклад:

```
username already exists
email already exists
```

<a id="s8-4"></a>
### 8.4 HTTP 500

500 Internal Server Error означає, що на сервері сталася непередбачена помилка.

Наприклад:

```
database error
JWT generation error
```

Клієнту не потрібно показувати внутрішню інформацію про сервер.

---

<a id="s9"></a>
## 9. Project Files

**handler.go** — Відповідає за HTTP endpoints:

```
RegisterHandler
LoginHandler
```

**service.go** — Відповідає за бізнес-логіку:

```
RegisterUser
LoginUser
ValidateRegister
```

**repository.go** — Відповідає за PostgreSQL:

```
CreateUser
GetUserByEmail
```

**token.go** — Відповідає за створення JWT:

```
TokenService
GenerateToken
```

**middleware.go** — Відповідає за перевірку JWT:

```
AuthMiddleware
```

**model.go** — Містить структури даних:

```
User
RegisterRequest
RegisterResponse
LoginRequest
LoginResponse
```

---

<a id="s10"></a>
## 10. Complete Auth Flow

**Registration**

```
Client
   ↓
POST /api/auth/register
   ↓
RegisterHandler
   ↓
ValidateRegister
   ↓
bcrypt password hash
   ↓
Create User
   ↓
Repository
   ↓
PostgreSQL
   ↓
Registration successful
```

**Login**

```
Client
   ↓
POST /api/auth/login
   ↓
LoginHandler
   ↓
GetUserByEmail
   ↓
PostgreSQL
   ↓
bcrypt password verification
   ↓
TokenService
   ↓
JWT
   ↓
Client
```

**Protected Request**

```
Client
   │
   │ Authorization: Bearer JWT
   ▼
Auth Middleware
   │
   ├── Check JWT
   ├── Check signature
   ├── Check expiration
   └── Get userID
   │
   ▼
Context
   │
   ▼
Protected Handler
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

---

<a id="s11"></a>
## 11. Current Status

Auth Module зараз реалізує:

**Registration**
- Request validation
- Password hashing with bcrypt
- PostgreSQL persistence
- Unique username
- Unique email

**Login**
- Password verification
- JWT generation
- JWT expiration
- JWT signature verification

**Bearer authentication**
- User ID extraction
- Context-based user identification
- Protected endpoints

Auth Module готовий як фундамент для наступних ігрових систем.

---

<a id="s12"></a>
## 12. Next Module — Player

Наступним буде Player Module.

Важливо розділити акаунт і ігрового персонажа.

`users` зберігає інформацію про акаунт:

```
users
    ↓
Account
    │
    ├── username
    ├── email
    └── password_hash
```

`players` буде зберігати ігрові дані:

```
players
    ↓
Game Data
    │
    ├── level
    ├── experience
    ├── currency
    └── progression
```

JWT дає нам:

```
JWT
 ↓
userID
```

Після цього Player Module зможе знайти дані цього користувача:

```
JWT
 ↓
userID
 ↓
Player
 ↓
Game Data
```

Це стане основою для наступних ігрових систем Prospect.
