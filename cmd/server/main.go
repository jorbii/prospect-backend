package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"prospect/internal/auth"
	"prospect/internal/database"
	"prospect/internal/inventory"
	"prospect/internal/player"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	db, err := database.NewPostgresPool()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer db.Close()

	repository := auth.NewRepository(db)
	tokenService := auth.NewTokenService(os.Getenv("JWT_SECRET"))

	playerRepository := player.NewRepository(db)
	playerService := player.NewService(playerRepository)
	playerHandler := player.NewHandler(playerService)

	inventoryRepository := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepository)
	inventoryHandler := inventory.NewHandler(
		inventoryService,
		playerService,
	)

	service := auth.NewService(
		db,
		repository,
		tokenService,
		playerService,
	)

	fmt.Println("Database connected")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/auth/register", service.RegisterHandler)
	http.HandleFunc("/api/auth/login", service.LoginHandler)

	http.Handle(
		"/api/auth/test",
		tokenService.AuthMiddleware(
			http.HandlerFunc(auth.TestAuthHandler),
		),
	)

	http.Handle(
		"/api/player/me",
		tokenService.AuthMiddleware(
			http.HandlerFunc(playerHandler.GetMe),
		),
	)

	http.Handle(
		"/api/inventory",
		tokenService.AuthMiddleware(
			http.HandlerFunc(inventoryHandler.GetInventory),
		),
	)

	http.Handle(
		"/api/inventory/items",
		tokenService.AuthMiddleware(
			http.HandlerFunc(inventoryHandler.AddItem),
		),
	)

	http.Handle(
		"/api/inventory/items/",
		tokenService.AuthMiddleware(
			http.HandlerFunc(inventoryHandler.DeleteItem),
		),
	)

	fmt.Println("Prospect server started on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
