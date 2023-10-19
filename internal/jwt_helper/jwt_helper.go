package jwt_helper

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var SECRET = []byte(viper.GetString("JWT.SECRET_KEY"))

type JWTClaims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func (m JWTClaims) Validate() error {
	if m.UserID == uuid.Nil {
		return errors.New("invalid user id")
	}
	return nil
}

func GenerateJWT(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: id,
	})

	tokenString, err := token.SignedString(SECRET)
	if err != nil {
		log.Println("ERR: [Generate JWT Failed]: ", err)
		return "", err

	}
	return tokenString, nil
}

func GetJWTUserID(tk string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tk, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	} else {
		log.Println(err)
		return uuid.Nil, err
	}
}
