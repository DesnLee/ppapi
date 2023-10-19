package jwt_helper

import (
	"crypto/rand"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})
	secret, err := generateHMACSecret()
	if err != nil {
		log.Println("ERR: [Generate HMAC Secret Failed]: ", err)
		return "", err
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println("ERR: [Generate JWT Failed]: ", err)
		return "", err

	}
	return tokenString, nil
}

func generateHMACSecret() ([]byte, error) {
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
