package jwt_helper

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var SECRET_KEY_FILE_NAME = "secret_key"

func GenerateJWT(id pgtype.UUID) (string, error) {
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
	// fmt.Println("JWT: ", tk)
	return tk, nil
}

func ParseJWT(tk string) (*JWTClaims, error) {
	tokenString, err := decryptJWT(tk)
	if err != nil {
		log.Println("ERR: [Decrypt JWT Failed]: ", err)
		return &JWTClaims{}, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getHMACSecret()
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Println("ERR: [Parse JWT Failed]: ", err)
		return &JWTClaims{}, err
	}
}
