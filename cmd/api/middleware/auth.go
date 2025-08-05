package middleware

import (
	"fmt"
	"go-library-service/cmd/api/entity"
	"go-library-service/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
type Claims struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtKey = []byte(utils.RequiredEnv("JWT_SECRET"))

// AuthMiddleware verifies the JWT token and sets the user ID in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ResponseError{Error: "Authorization header is required", Code: http.StatusUnauthorized})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ResponseError{Error: "Bearer token not found", Code: http.StatusUnauthorized})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ResponseError{Error: "Invalid or expired token", Code: http.StatusUnauthorized})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// RoleMiddleware checks user has the required role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "User role not found in context", Code: http.StatusInternalServerError})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, entity.ResponseError{Error: "Invalid user role type", Code: http.StatusInternalServerError})
			return
		}

		for _, role := range roles {
			if role == roleStr {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, entity.ResponseError{Error: "Invalid permissions", Code: http.StatusForbidden})
	}
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
