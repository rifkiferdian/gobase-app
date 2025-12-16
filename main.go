package main

import (
	"log"
	"os"
	"stok-hadiah/config"
	"stok-hadiah/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to OS environment")
	}

	// Initialize database / config
	config.Connect()

	// Initialize Gin engine
	r := gin.Default()

	// Templates & static files
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/assets", "./assets")

	// SESSION - must be registered BEFORE routes that use sessions
	store := cookie.NewStore([]byte("secret-key"))
	r.Use(sessions.Sessions("mysession", store))

	// Register application routes
	routes.RegisterWebRoutes(r)

	// Determine port (fallback to 8080 if APP_PORT is not set)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server and log fatal on error
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
