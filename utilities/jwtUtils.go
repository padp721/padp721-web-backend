package utilities

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(privateKey *rsa.PrivateKey, username string) (string, time.Time, error) {
	now := time.Now()
	expires := now.Add(24 * time.Hour)

	jwtClaims := jwt.MapClaims{
		"iss": "admin",
		"sub": username,
		"aud": "users",
		"exp": expires.Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"jti": uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
	jwtString, err := token.SignedString(privateKey)
	if err != nil {
		return "", now, err
	}

	return jwtString, expires, nil
}

func ParseJWT(privateKey *rsa.PrivateKey, jwtTokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return LoadPublicKey(privateKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
