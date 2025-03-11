package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"strings"
	"time"
)

func RequireAuth(c *gin.Context) {
	//Get the jwt
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check if the Authorization header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract the token part from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	//Decode/Validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//Find the user with that token
		var user models.User

		uuidStr := claims["sub"].(string)

		//This is the uuid(id = uuid), i dint use uuid name because uuid already exist
		id, err := uuid.Parse(uuidStr)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		initializers.DB.First(&user, "id = ?", id)

		if user.ID == uuid.Nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//Check the password version
		if claims["pwd_version"].(string) != user.PasswordVersion.String() {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//Attach to req
		c.Set("user", user)

		//Continue

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
