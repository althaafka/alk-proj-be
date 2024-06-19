package helpers

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaim struct {
	UserID uint
	jwt.StandardClaims
}

func GenerateToken(userID uint) (string, error) {
	expTime := time.Now().Add(5 * time.Hour)

	claims := &JWTClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
