package main

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/routes"
	"strings"

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

	// Custom template functions tambah
	r.SetFuncMap(template.FuncMap{
		"no": func(a, b int) int {
			return a + b
		},
		"baseURL": func(path string) string {
			base := strings.TrimRight(os.Getenv("BASE_URL"), "/")
			p := "/" + strings.TrimLeft(path, "/")
			return base + p
		},
	})

	// Templates & static files
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/assets", "./assets")

	useSecureCookie := strings.ToLower(os.Getenv("APP_SECURE_COOKIE")) == "true"

	// Register custom session payload for gob encoder used by cookie store.
	gob.Register(models.SessionUser{})

	// SESSION - must be registered BEFORE routes that use sessions
	store := cookie.NewStore([]byte("secret-key"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 8, // 8 jam
		HttpOnly: true,
		// Secure harus false saat akses lokal HTTP; aktifkan otomatis jika APP_ENV=production atau APP_SECURE_COOKIE=true.
		Secure:   useSecureCookie,
		SameSite: http.SameSiteLaxMode,
	})
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
