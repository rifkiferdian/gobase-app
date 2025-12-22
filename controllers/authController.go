package controllers

import (
	"database/sql"
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const userModelType = "Models\\User"

func LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil {
		c.Redirect(302, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Login User",
	})
}

func LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
	// 	return
	// }

	// fmt.Println("DEBUG:", string(hashedPassword))

	var (
		userID int
		dbUser string
		dbName string
		dbPass string
		dbNip  sql.NullString
		dbRole sql.NullString
	)
	err := config.DB.QueryRow(`
		SELECT 
			u.id,
			u.username,
			u.name,
			u.password,
			COALESCE(u.nip, '') AS nip,
			COALESCE(r.name, '') AS role
		FROM users u
		LEFT JOIN model_has_roles mhr ON mhr.model_id = u.id 
		LEFT JOIN roles r ON r.id = mhr.role_id
		WHERE u.username = ?
	`, username).
		Scan(&userID, &dbUser, &dbName, &dbPass, &dbNip, &dbRole)

	if err == sql.ErrNoRows {
		c.HTML(200, "login.html", gin.H{
			"Title": "Login User",
			"Error": "Username tidak ditemukan",
		})
		return
	} else if err != nil {
		c.HTML(500, "login.html", gin.H{
			"Title": "Login User",
			"Error": "Terjadi kesalahan saat mengambil data user",
		})
		return
	}

	// cek password
	if bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password)) != nil {
		c.HTML(200, "login.html", gin.H{
			"Title": "Login User",
			"Error": "Password salah",
		})
		return
	}

	// simpan session
	session := sessions.Default(c)
	session.Set("user", models.SessionUser{
		UserID:          userID,
		NIP:             dbNip.String,
		Name:            dbName,
		Username:        dbUser,
		Role:            dbRole.String,
		IsAuthenticated: true,
	})
	if err := session.Save(); err != nil {
		c.HTML(500, "login.html", gin.H{
			"Title": "Login User",
			"Error": "Gagal menyimpan sesi: " + err.Error(),
		})
		return
	}

	c.Redirect(302, "/dashboard")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(302, "/")
}

func CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if username already exists
	var existingUser string
	err := config.DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert new user
	_, err = config.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
