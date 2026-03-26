package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Token hashes a token using SHA256 for O(1) database lookups
func SHA256Token(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
