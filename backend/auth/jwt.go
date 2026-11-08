package auth

import (
	"errors"
	"fmt"
	"hrm-backend/config"
	"hrm-backend/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const claimsContextKey = "authClaims"

type Claims struct {
	EmployeeDBID uint   `json:"employee_db_id"`
	EmployeeID   string `json:"employee_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte
var tokenExpiry = 8 * time.Hour

func Initialize() error {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return errors.New("JWT_SECRET must be set and at least 32 characters long")
	}

	expiryHours := int64(8)
	if rawExpiry := strings.TrimSpace(os.Getenv("JWT_EXPIRY_HOURS")); rawExpiry != "" {
		parsedExpiry, err := strconv.ParseInt(rawExpiry, 10, 64)
		maxExpiryHours := int64(^uint64(0)>>1) / int64(time.Hour)
		if err != nil || parsedExpiry <= 0 || parsedExpiry > maxExpiryHours {
			return errors.New("JWT_EXPIRY_HOURS must be a positive integer")
		}
		expiryHours = parsedExpiry
	}

	jwtSecret = []byte(secret)
	tokenExpiry = time.Duration(expiryHours) * time.Hour
	return nil
}

func GenerateToken(employeeDBID uint, employeeID, email, role string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT configuration is not initialized")
	}
	now := time.Now()
	claims := Claims{
		EmployeeDBID: employeeDBID,
		EmployeeID:   employeeID,
		Email:        email,
		Role:         role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, prefix))
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || token == nil || !token.Valid || claims.EmployeeDBID == 0 || claims.EmployeeID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		var employee models.Employee
		if err := config.DB.Select("id", "employee_id", "email", "role").First(&employee, claims.EmployeeDBID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to authenticate"})
			return
		}
		claims.EmployeeID = employee.EmployeeID
		claims.Email = employee.Email
		claims.Role = employee.Role
		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		claims, ok := c.Get(claimsContextKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		userClaims, ok := claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication"})
			return
		}
		if _, ok := allowed[userClaims.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}

func CurrentClaims(c *gin.Context) (*Claims, bool) {
	claims, ok := c.Get(claimsContextKey)
	if !ok {
		return nil, false
	}
	userClaims, ok := claims.(*Claims)
	return userClaims, ok
}
