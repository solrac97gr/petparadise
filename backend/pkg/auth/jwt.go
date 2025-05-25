package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/solrac97gr/petparadise/internal/users/domain/models"
	"github.com/solrac97gr/petparadise/pkg/config"
)

// Configuration constants
const (
	AccessTokenExpiration  = time.Minute * 15   // 15 minutes
	RefreshTokenExpiration = time.Hour * 24 * 7 // 7 days
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

	// Start the cleanup ticker to periodically remove expired tokens
}

// AccessClaims represents the JWT claims for access tokens
type AccessClaims struct {
	UserID string      `json:"user_id"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// RefreshClaims represents the JWT claims for refresh tokens
type RefreshClaims struct {
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"` // Unique ID for this refresh token
	jwt.RegisteredClaims
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // Seconds until access token expires
}

// GenerateTokenPair creates a new pair of access and refresh tokens for a user
func GenerateTokenPair(user *models.User) (*TokenPair, error) {
	// Generate unique token ID for the refresh token
	tokenID := fmt.Sprintf("%s-%d", user.ID, time.Now().UnixNano())

	// Create access token
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, err := generateRefreshToken(user.ID, tokenID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenExpiration.Seconds()),
	}, nil
}

// generateAccessToken creates a new access token for a user
func generateAccessToken(user *models.User) (string, error) {
	// Set claims
	claims := &AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
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

// generateRefreshToken creates a new refresh token for a user
func generateRefreshToken(userID, tokenID string) (string, error) {
	// Set claims
	claims := &RefreshClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
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

// ValidateAccessToken validates a JWT access token and returns the claims
func ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Check for token expiration error specifically
	if err != nil && strings.Contains(err.Error(), "token is expired") {
		return nil, ErrExpiredToken
	} else if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		// Double check if token has expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, ErrExpiredToken
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateRefreshToken validates a JWT refresh token and returns the claims
func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Check for token expiration error specifically
	if err != nil && strings.Contains(err.Error(), "token is expired") {
		return nil, ErrExpiredToken
	} else if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		// Double check if token has expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, ErrExpiredToken
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshTokens uses a valid refresh token to generate a new token pair
func RefreshTokens(refreshToken string) (*TokenPair, error) {
	// Validate the refresh token
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// In a real application, you would verify the user exists and is still active
	// For simplicity, we'll create a minimal user object with the ID from the claims
	user := &models.User{
		ID: claims.UserID,
		// In a real implementation, you would fetch other user data from the database
		Role: models.RoleUser, // Default role - in real implementation, get from DB
	}

	// Generate new token pair
	tokenPair, err := GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}
