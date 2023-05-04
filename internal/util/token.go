package util

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/google/uuid"
)

// generateUniqueToken returns a unique token based on a given GUID
// The token is generated using the md5
func GenerateUniqueToken() string {
	// Generate a new GUID
    newGUID := uuid.New()
    hash := md5.Sum([]byte(newGUID.String()))
    return hex.EncodeToString(hash[:])
}


// type JwtCustomClaims struct {
// 	UserID string `json:"userID"`
// 	jwt.StandardClaims
// }

// func CreateAccessToken(userID int64, secret string, expiry int) (string, error) {
// 	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
// 	claims := &JwtCustomClaims{
// 		UserID: strconv.FormatInt(userID, 10),
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: exp,
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	accessToken, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", err
// 	}

// 	return accessToken, err
// }

// func ExtractIDFromToken(requestToken string, secret string) (string, error) {
// 	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(secret), nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok && !token.Valid {
// 		return "", fmt.Errorf("invalid Token")
// 	}

// 	userID, ok := claims["userID"].(string)
// 	if !ok {
//         return "", fmt.Errorf("invalid Token")
//     }
// 	return userID, nil
// }
