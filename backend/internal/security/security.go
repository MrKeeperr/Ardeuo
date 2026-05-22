package security

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const defaultTokenTTL = 24 * time.Hour

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) (string, error) {
	return hashPassword(password)
}

func CheckPasswordHash(password, hash string) bool {
	return checkPasswordHash(password, hash)
}

type Claims struct {
	UserID int    `json:"user_id"`
	RoleID int    `json:"role_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, roleID int, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	ttl := defaultTokenTTL
	if ttlHours := os.Getenv("JWT_TTL_HOURS"); ttlHours != "" {
		hours, err := strconv.Atoi(ttlHours)
		if err != nil || hours <= 0 {
			return "", errors.New("invalid JWT_TTL_HOURS")
		}
		ttl = time.Duration(hours) * time.Hour
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
