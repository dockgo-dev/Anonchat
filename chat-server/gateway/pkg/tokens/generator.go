package tokens

import (
	"crypto/rand"
	"math/big"
)

const sim = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"

func Generate(n int) string {
	str := make([]byte, n)
	for i := range str {
		number, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
		str[i] = sim[number.Int64()]
	}
	return string(str)
}
