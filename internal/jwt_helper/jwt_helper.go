package jwt_helper

import (
	"crypto/rand"
	"errors"
	"log"
	"os"
	"path"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var SECRET_KEY_FILE_NAME = "secret_key"

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
	return tokenString, nil
}

func GetJWTUserID(tk string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tk, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getHMACSecret()
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	} else {
		log.Println(err)
		return uuid.Nil, err
	}
}

func getHMACSecret() ([]byte, error) {
	keyPath := viper.GetString("JWT.SECRET_KEY_PATH")
	if keyPath == "" {
		keyPath = "./"
	}

	// 读取文件
	filePath := path.Join(keyPath, SECRET_KEY_FILE_NAME)
	keyByte, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			key := make([]byte, 64)
			_, err := rand.Read(key)
			if err != nil {
				log.Fatalln("ERR: [Generate JWT Secret Key Failed]: ", err)
				return nil, err
			}
			err = os.WriteFile(filePath, key, 0644)
			if err != nil {
				log.Fatalln("ERR: [Write JWT Secret Key File Failed]: ", err)
				return nil, err
			}
			return getHMACSecret()
		} else {
			log.Fatalln("ERR: [Read JWT Secret Key File Failed]: ", err)
			return nil, err
		}
	}

	return keyByte, nil
}
