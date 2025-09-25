package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}
var nextUserID = 3

func main() {
	r := gin.Default()
	
	r.GET("/health", healthHandler)
	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/users", getUsersHandler)
	r.GET("/users/:id", getUserHandler)
	
	port := getEnv("PORT", "8080")
	r.Run(":" + port)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "user-service",
		"time":    time.Now().Unix(),
	})
}

func registerHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.ID = nextUserID
	nextUserID++
	
	users = append(users, user)
	user.Password = "" // Don't return password
	
	c.JSON(http.StatusCreated, user)
}

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate JWT token (simplified)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("secret"))
	
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  gin.H{"email": req.Email},
	})
}

func getUsersHandler(c *gin.Context) {
	safeUsers := make([]User, len(users))
	for i, user := range users {
		safeUsers[i] = User{ID: user.ID, Name: user.Name, Email: user.Email}
	}
	c.JSON(http.StatusOK, safeUsers)
}

func getUserHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, user := range users {
		if user.ID == id {
			safeUser := User{ID: user.ID, Name: user.Name, Email: user.Email}
			c.JSON(http.StatusOK, safeUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}