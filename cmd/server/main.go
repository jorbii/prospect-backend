package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"prospect/internal/auth"
	"prospect/internal/database"
	"prospect/internal/inventory"
	"prospect/internal/item"
	"prospect/internal/loadout"
	"prospect/internal/loot"
	"prospect/internal/player"
	"prospect/internal/weapon"
)

func main() {

	// ========================================
	// Environment
	// ========================================

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// ========================================
	// Database
	// ========================================

	db, err := database.NewPostgresPool()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer db.Close()

	// ========================================
	// Auth Module
	// ========================================

	repository := auth.NewRepository(db)
	tokenService := auth.NewTokenService(os.Getenv("JWT_SECRET"))

	service := auth.NewService(
		db,
		repository,
		tokenService,
		nil,
	)

	// ========================================
	// Player Module
	// ========================================

	playerRepository := player.NewRepository(db)
	playerService := player.NewService(playerRepository)
	playerHandler := player.NewHandler(playerService)

	// ========================================
	// Auth Service
	// ========================================

	service = auth.NewService(
		db,
		repository,
		tokenService,
		playerService,
	)

	// ========================================
	// Inventory Module
	// ========================================

	inventoryRepository := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepository)
	inventoryHandler := inventory.NewHandler(
		inventoryService,
		playerService,
	)

	// ========================================
	// Item Module
	// ========================================

	itemRepository := item.NewRepository(db)
	itemService := item.NewService(itemRepository)
	itemHandler := item.NewHandler(itemService)

	fmt.Println("Database connected")

	// ========================================
	// weapon Module
	// ========================================

	weaponRepository := weapon.NewRepository(db)
	weaponService := weapon.NewService(weaponRepository)
	weaponHandler := weapon.NewHandler(weaponService)

	// ========================================
	// Health
	// ========================================

	http.HandleFunc("/health", healthHandler)

	// ========================================
	// loot Module
	// ========================================

	lootRepository := loot.NewRepository(db)

	lootService := loot.NewService(
		lootRepository,
		playerRepository,
	)

	lootHandler := loot.NewHandler(lootService)

	// ========================================
	// loadout Module
	// ========================================

	loadoutRepository := loadout.NewRepository(db)

	loadoutService := loadout.NewService(
		loadoutRepository,
		playerRepository,
	)

	loadoutHandler := loadout.NewHandler(loadoutService)

	// ========================================
	// Auth Routes
	// ========================================

	http.HandleFunc(
		"/api/auth/register",
		service.RegisterHandler,
	)

	http.HandleFunc(
		"/api/auth/login",
		service.LoginHandler,
	)

	http.Handle(
		"/api/auth/test",
		tokenService.AuthMiddleware(
			http.HandlerFunc(auth.TestAuthHandler),
		),
	)

	// ========================================
	// Player Routes
	// ========================================

	http.Handle(
		"/api/player/me",
		tokenService.AuthMiddleware(
			http.HandlerFunc(playerHandler.GetMe),
		),
	)

	// ========================================
	// Inventory Routes
	// ========================================

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

	// ========================================
	// Item Routes
	// ========================================

	http.HandleFunc(
		"/api/items",
		itemHandler.GetAll,
	)

	http.HandleFunc(
		"/api/items/",
		itemHandler.GetByID,
	)

	// ========================================
	// weapon Routes
	// ========================================

	http.HandleFunc(
		"/api/weapons",
		weaponHandler.GetAll,
	)

	http.HandleFunc(
		"/api/weapons/",
		weaponHandler.GetByID,
	)

	// ========================================
	// loot Route
	// ========================================

	http.Handle(
		"/api/loot/",
		tokenService.AuthMiddleware(
			http.HandlerFunc(lootHandler.GetByID),
		),
	)

	http.Handle(
		"/api/loot/pickup/",
		tokenService.AuthMiddleware(
			http.HandlerFunc(lootHandler.Pickup),
		),
	)

	http.Handle(
		"/api/loot",
		tokenService.AuthMiddleware(
			http.HandlerFunc(lootHandler.Create),
		),
	)

	// ========================================
	// loadout Route
	// ========================================

	http.Handle(
		"/api/loadout",
		tokenService.AuthMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					loadoutHandler.Get(w, r)

				case http.MethodPost:
					loadoutHandler.Create(w, r)

				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}
			}),
		),
	)

	// ========================================
	// Start Server
	// ========================================

	fmt.Println("Prospect server started on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
