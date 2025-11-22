package tokens

import (
	crand "crypto/rand"
	"encoding/base64"
	rand "math/rand"
)

const alphanum = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789"

func GenerateRefresh() string {
	token := make([]byte, 32)
	if _, err := crand.Read(token); err != nil {
		for i := range token {
			number := rand.Intn(len(alphanum))
			token[i] = alphanum[number]
		}
	}
	return base64.RawURLEncoding.EncodeToString(token)
}
