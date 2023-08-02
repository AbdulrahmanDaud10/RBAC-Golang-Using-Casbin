package api

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password *string) {
	bytePass := []byte(*password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	*password = string(hashedPassword)
}

// ComparePassword -> compares database password to the provided passord
func ComparePassword(DbPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(DbPassword), []byte(password)) == nil
}

// GenerateToken -> generates token
func GenerateToken(userID uint) string {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour * 3).Unix(),
		"iat":    time.Now().Unix(),
		"userID": userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return ""
	}

	return secretToken
}

// ValidateToken --> validate the given token
func ValidateToken(token string) (*jwt.Token, error) {

	// function return secret key after checking if the signing method is HMAC and returned key is used by 'Parse' to decode the token)
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
