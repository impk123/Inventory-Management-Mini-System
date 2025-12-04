package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/impk123/Inventory-Management-Mini-System/internal/models"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	jwtSecret string
	db        *gorm.DB
}

func NewAuthMiddleware(jwtSecret string, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		db:        db,
	}
}

func (m *AuthMiddleware) ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// มี Header มาไหม? (Authorization)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// รูปแบบถูกต้องไหม? (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Token ปลอมหรือเปล่า? (Parse)
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ตรวจสอบ claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ตรวจสอบว่า token ยังไม่หมดอายุไหม?
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		// 5. User ยังมีตัวตนอยู่จริงไหม?
		var user models.User
		userID := uint(claims["user_id"].(float64)) // แปลง ID จาก Token

		if err := m.db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// 6. ผ่านหมด! เก็บข้อมูลลง Context
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role) // เผื่อเอาไปเช็คสิทธิ์ต่อ
		c.Next()
	}
}
