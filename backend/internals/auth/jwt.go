package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Subject string    `json:"sub"`
	Issued  time.Time `json:"isa"`
	Expires time.Time `json:"exp"`
	jwt.RegisteredClaims
}

func Verifyer(token, secret string) (claims Claims, err error) {
	var c Claims
	jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	return c, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Yoily",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}).SignedString([]byte(tokenSecret))
}
