package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"math/rand"
	"time"
)

func CreateToken(userID int64, seed string) string {
	// Convert userID to byte slice
	userIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(userIDBytes, uint64(userID))
	// Concatenate userID bytes and seed
	data := append(userIDBytes, []byte(seed)...)
    // Generate MD5 hash of data
	hash := md5.Sum(data)
	// Convert hash to hex string
	token := hex.EncodeToString(hash[:])
	return token
}

func ExtractUserIDFromToken(token string, seed string) (int64, error) {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
	  return 0, err
	}
	// Extract hash from token (last 32 bytes)
	hashBytes := tokenBytes[len(tokenBytes)-32:]
	// Extract userID bytes from token (first 8 bytes)
	userIDBytes := tokenBytes[:8]
	// Concatenate userID bytes and seed
	data := append(userIDBytes, []byte(seed)...)
	// Generate MD5 hash of data
	hash := md5.Sum(data)
	// Compare calculated hash with hash from token
	if !bytes.Equal(hash[:], hashBytes) {
	  return 0, errors.New("invalid token")
	}
	// Convert userID bytes to int64
	userID := int64(binary.BigEndian.Uint64(userIDBytes))
	return userID, nil
  }

// func ExtractUserIDFromToken(token string) (int64, error) {
// 	tokenBytes, err := hex.DecodeString(token)
// 	if err != nil {
// 		return 0, err
// 	}
// 	// Extract first 8 bytes of token (representing userID)
// 	userIDBytes := tokenBytes[:8]
// 	// Convert userID bytes to int64
//  	userID := int64(binary.BigEndian.Uint64(userIDBytes))
// 	return userID, nil
// }
// 	hasher.Write([]byte(seed))
// 	idBytes := hasher.Sum(decoded)
// 	return int64(idBytes[0]), nil
// }

func GenerateToken(algorithm string, lifetime int) (string, error) {
	rand.Seed(time.Now().UnixNano())
	guid := fmt.Sprintf("%s", rand.New(rand.NewSource(time.Now().UnixNano())).Uint64())
	var hasher hash.Hash
	switch algorithm {
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	default:
		return "", fmt.Errorf("Unsupported algorithm %s", algorithm)
	}
	hasher.Write([]byte(guid))
	token := hex.EncodeToString(hasher.Sum(nil))

	// Check lifetime
	if lifetime > 0 {
		now := time.Now().Unix()
		expiration := now + int64(lifetime)
		return fmt.Sprintf("%s|%d", token, expiration), nil
	}

	return token, nil
}



// func CreateToken(userID int64, seed string) string {
// 	// Convert userID to byte slice
// 	userIDBytes := make([]byte, 8)
// 	binary.BigEndian.PutUint64(userIDBytes, uint64(userID))

// 	// Concatenate seed and userID bytes
// 	data := append([]byte(seed), userIDBytes...)

// 	// Generate MD5 hash of data
// 	hash := md5.Sum(data)

// 	// Convert hash to hex string
// 	return hex.EncodeToString(hash[:])
// }

// func ExtractUserIDFromToken(token, seed string) (int64, error) {
// 	// Decode token from hex string to byte slice
// 	tokenBytes, err := hex.DecodeString(token)
// 	if err != nil {
// 		return 0, err
// 	}

// 	lenSeedBytes := len([]byte(seed))
// 	// Extract first 8 bytes of token (representing userID)
// 	userIDBytes := tokenBytes[lenSeedBytes : lenSeedBytes+8]

// 	// Convert userID bytes to int64
// 	userID := int64(binary.BigEndian.Uint64(userIDBytes))

// 	return userID, nil
// }

// func CreateToken(seed string, userID int64) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(seed))
// 	hasher.Write([]byte(fmt.Sprintf("%d", userID)))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// func ExtractUserIDFromToken(token string, seed string) (int64, error) {
// 	decoded, err := hex.DecodeString(token)
// 	if err != nil {
// 		return 0, err
// 	}

// 	hasher := md5.New()