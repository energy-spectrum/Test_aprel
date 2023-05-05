package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"app/config"

	"github.com/google/uuid"
)

// GenerateUniqueToken returns a unique token based on UUID.
// The token is generated using the md5.
func GenerateUniqueToken() string {
	newGUID := uuid.New()
	hash := md5.Sum([]byte(newGUID.String()))
	return hex.EncodeToString(hash[:])
}

func GetExpirationTimeForToken(env *config.Env) time.Time {
	lifetime := time.Hour * time.Duration(env.TokenExpiryHour)
	expirationTime := time.Now().Add(lifetime)
	expirationTime = expirationTime.Truncate(100 * time.Microsecond)

	return expirationTime
}
