package jwt_helper

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var SECRET_KEY_FILE_NAME = "secret_key"

func GenerateJWT(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: id,
	})

	secret, err := getHMACSecret()
	if err != nil {
		log.Println("ERR: [Generate JWT Failed]: ", err)
		return "", err
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println("ERR: [Generate JWT Failed]: ", err)
		return "", err
	}

	tk, err := encryptJWT(tokenString)
	if err != nil {
		log.Println("ERR: [Encrypt JWT Failed]: ", err)
		return "", err
	}
	return tk, nil
}

func GetJWTUserID(tk string) (uuid.UUID, error) {
	tokenString, err := decryptJWT(tk)
	if err != nil {
		log.Fatalln("ERR: [Decrypt JWT Failed]: ", err)
		return uuid.Nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getHMACSecret()
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	} else {
		log.Println("ERR: [Parse JWT Failed]: ", err)
		return uuid.Nil, err
	}
}
