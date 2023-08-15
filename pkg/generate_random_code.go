package pkg

import (
	"crypto/rand"
)

func GenerateRandomCode(l int) (code string, error error) {
	// 伪随机 math/rand
	// 真随机 crypto/rand

	r := make([]byte, l)
	digs := make([]byte, l)

	_, err := rand.Read(r)
	if err != nil {
		code = ""
		error = err
		return
	}

	for i, v := range r {
		digs[i] = v%10 + 48
	}

	code = string(digs)
	error = nil
	return
}
