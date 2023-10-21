package jwt_helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

// getFilePath 获取密钥文件路径
func getFilePath() string {
	keyPath := viper.GetString("JWT.SECRET_KEY_PATH")
	if keyPath == "" {
		keyPath = "./"
	}
	return path.Join(keyPath, SECRET_KEY_FILE_NAME)
}

// generateHMACSecret 生成 HMAC 密钥
func generateHMACSecret() ([]byte, error) {
	filePath := getFilePath()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Println("ERR: [Generate JWT Secret Key Failed]: ", err)
		return nil, err
	}
	err = os.WriteFile(filePath, key, 0644)
	if err != nil {
		log.Println("ERR: [Write JWT Secret Key File Failed]: ", err)
		return nil, err
	}

	return key, nil
}

// getHMACSecret 获取/生成 JWT 的密钥
func getHMACSecret() ([]byte, error) {
	filePath := getFilePath()
	keyByte, err := os.ReadFile(filePath)

	if err == nil {
		return keyByte, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		log.Println("ERR: [Read JWT Secret Key File Failed]: ", err)
		return nil, err
	}

	return generateHMACSecret()
}

// encryptJWT 加密 JWT
func encryptJWT(tk string) (string, error) {
	secret, err := getHMACSecret()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(tk))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(tk))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decryptJWT 解密 JWT
func decryptJWT(tk string) (string, error) {
	secret, err := getHMACSecret()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(tk)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext is too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
