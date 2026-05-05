package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key_change_me_in_prod") // In a real app, use environment variables

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// HashPassword generates a bcrypt hash from a password string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT creates a new token for a user
func GenerateJWT(userID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, expirationTime, err
}

// SetAuthCookie sets an HttpOnly cookie with the JWT
func SetAuthCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

// GetUserIDFromToken extracts the user ID from a JWT token string
func GetUserIDFromToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}
