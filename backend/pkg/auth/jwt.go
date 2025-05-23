package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/pkg/config"
)

// Configuration constants
const (
	TokenExpiration = time.Hour * 24 // 24 hours
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWT secret key - will be set from configuration
var jwtSecret []byte

// InitJWTSecret initializes the JWT secret from config
func InitJWTSecret(cfg *config.Config) {
	jwtSecret = []byte(cfg.JWTSecret)
}

// Claims represents the JWT claims
type Claims struct {
	UserID string      `json:"user_id"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(user *models.User) (string, error) {
	// Set claims
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token has expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, ErrExpiredToken
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// SetJWTSecret sets the JWT secret key (useful for tests or configuration)
func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}
