// Command server boots the HardwarePOS REST API: applies migrations,
// wires each module's service/handler pair, and starts the Gin engine.
package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"hardwarepos/backend/internal/auth"
	"hardwarepos/backend/internal/categories"
	"hardwarepos/backend/internal/config"
	"hardwarepos/backend/internal/customers"
	"hardwarepos/backend/internal/dashboard"
	"hardwarepos/backend/internal/database"
	"hardwarepos/backend/internal/inventory"
	"hardwarepos/backend/internal/middleware"
	"hardwarepos/backend/internal/products"
	"hardwarepos/backend/internal/purchases"
	"hardwarepos/backend/internal/reports"
	"hardwarepos/backend/internal/returns"
	"hardwarepos/backend/internal/sales"
	"hardwarepos/backend/internal/settings"
	"hardwarepos/backend/internal/suppliers"
	"hardwarepos/backend/internal/users"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DBPath, cfg.Migrations)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))

	api := router.Group("/api")

	authSvc := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	authed := api.Group("")
	authed.Use(middleware.Auth(cfg.JWTSecret))

	authHandler.RegisterRoutes(api, authed)

	categoriesHandler := categories.NewHandler(categories.NewService(db))
	categoriesHandler.RegisterRoutes(authed)

	suppliersHandler := suppliers.NewHandler(suppliers.NewService(db))
	suppliersHandler.RegisterRoutes(authed)

	productsSvc := products.NewService(db)
	productsHandler := products.NewHandler(productsSvc)
	productsHandler.RegisterRoutes(authed)

	inventoryHandler := inventory.NewHandler(inventory.NewService(db))
	inventoryHandler.RegisterRoutes(authed)

	purchasesHandler := purchases.NewHandler(purchases.NewService(db))
	purchasesHandler.RegisterRoutes(authed)

	salesSvc := sales.NewService(db)
	salesHandler := sales.NewHandler(salesSvc)
	salesHandler.RegisterRoutes(authed)

	returnsHandler := returns.NewHandler(returns.NewService(db))
	returnsHandler.RegisterRoutes(authed)

	customersHandler := customers.NewHandler(customers.NewService(db))
	customersHandler.RegisterRoutes(authed)

	usersHandler := users.NewHandler(users.NewService(db))
	usersHandler.RegisterRoutes(authed)

	dashboardHandler := dashboard.NewHandler(dashboard.NewService(db, productsSvc, salesSvc))
	dashboardHandler.RegisterRoutes(authed)

	reportsHandler := reports.NewHandler(reports.NewService(db))
	reportsHandler.RegisterRoutes(authed)

	settingsHandler := settings.NewHandler(settings.NewService(db))
	settingsHandler.RegisterRoutes(authed)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("HardwarePOS API listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
