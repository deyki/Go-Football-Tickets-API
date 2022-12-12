package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Go-FootballTickets/deyki/v2/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {

	tokenString, errorMessage := c.Cookie("Authorization")
	if errorMessage != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var admin database.Admin

		db, err := database.ConnectDB()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		db.First(&admin, claims["sub"])

		if admin.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("admin", admin)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}