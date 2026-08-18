package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"prospect/internal/auth"
	"prospect/internal/database"
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
	service := auth.NewService(repository, tokenService)

	fmt.Println("Database connected")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/auth/register", service.RegisterHandler)
	http.HandleFunc("/api/auth/login", service.LoginHandler)

	http.Handle(
		"/api/auth/test",
		tokenService.AuthMiddleware(http.HandlerFunc(auth.TestAuthHandler)),
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
