package xcrypto

import (
	"fmt"

	"github.com/spaolacci/murmur3"
)

func MurMur3Hash(plaintext string) string {
	h1, h2 := murmur3.Sum128([]byte(plaintext))
	return fmt.Sprintf("%d-%d", h1, h2)
}
