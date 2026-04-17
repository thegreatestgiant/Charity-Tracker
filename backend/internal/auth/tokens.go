package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Verifyer(tokenString string, secret []byte) (claims jwt.Claims, err error) {
	c := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		log.Printf("Cannot Parse Token: %v", err)
		return nil, err
	}
	if !token.Valid {
		log.Println("Invalid token:")
		return nil, errors.New("Invalid Token")
	}
	return c, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret []byte, expiresIn time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Yoily",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}).SignedString(tokenSecret)
}
