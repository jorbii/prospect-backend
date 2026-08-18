
````markdown
# Prospect Backend

Backend server for **Prospect**, a multiplayer mobile game.

The server is responsible for user authentication, authorization, database communication, and API handling.

## Tech Stack

- **Go** — backend programming language
- **PostgreSQL** — database
- **JWT** — authentication
- **REST API** — communication between client and server
- **Git** — version control

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
│   ├── user/
│   └── ...
│
├── migrations/
│
├── .env
├── .gitignore
├── go.mod
└── README.md
````

## Authentication

The backend uses **JWT (JSON Web Token)** for user authentication.

The basic authentication flow:

```text
Client
  │
  │ Login / Register
  ▼
Backend
  │
  │ Validate user
  ▼
PostgreSQL
  │
  │ User data
  ▼
Backend
  │
  │ JWT
  ▼
Client
```

After successful authentication, the server returns a JWT token.

The client sends this token with protected API requests.

```http
Authorization: Bearer <token>
```

The backend validates the token before allowing access to protected endpoints.

## Environment Variables

The project uses environment variables for configuration and sensitive data.

Create a `.env` file in the project root:

```env
DATABASE_URL=your_database_url
JWT_SECRET=your_secret_key
PORT=8080
```

> Never commit `.env` to GitHub.

The `.env` file is included in `.gitignore`.

## Installation

Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/prospect-backend.git
cd prospect-backend
```

Install Go dependencies:

```bash
go mod download
```

Create your `.env` file and configure the required variables.

## Run the Server

Start the server with:

```bash
go run ./cmd/server
```

If everything is configured correctly, the server will start on the configured port.

## Development

Run the project:

```bash
go run ./cmd/server
```

Format Go code:

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

## API

The backend provides REST API endpoints for communication with the game client.

Current functionality includes:

* User registration
* User login
* JWT authentication
* Protected API routes
* User data management
* PostgreSQL integration

More game-related services and endpoints will be added as development continues.

## Roadmap

* [x] User registration
* [x] User login
* [x] JWT authentication
* [x] PostgreSQL integration
* [ ] User profile
* [ ] Game session management
* [ ] Matchmaking
* [ ] Game server
* [ ] Inventory
* [ ] Skins
* [ ] Player statistics
* [ ] Real-time communication
* [ ] Production deployment

## Project Status

🚧 **In development**

Prospect Backend is currently under active development. The architecture and API may change as new game systems are implemented.

## License

This project is currently private and is not licensed for redistribution.

```
```
