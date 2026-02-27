package xcrypto

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
)

func SHA3Hash(plaintext string) string {
	input := "secure_data_to_hash"

	hash := sha3.Sum256([]byte(input))
	fmt.Printf("SHA3-256 Hash: %x\n", hash)

	hash512 := sha3.Sum512([]byte(input))

	return hex.EncodeToString(hash512[:])
}
