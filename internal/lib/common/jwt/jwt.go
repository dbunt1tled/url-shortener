package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// time.Now().Add(time.Second * time.Duration(exp)).Unix(),
func encode(payload map[string]interface{}) string {
	var claims jwt.MapClaims
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	pem.
	return token.SignedString([]byte("secret"))
}
