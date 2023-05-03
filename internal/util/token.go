package util

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func CreateAccessToken(userID int64, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &JwtCustomClaims{
		UserID: strconv.FormatInt(userID, 10),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessToken, err
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid Token")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
        return "", fmt.Errorf("invalid Token")
    }
	return userID, nil
}
