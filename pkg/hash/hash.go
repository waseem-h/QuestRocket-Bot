package hash

import (
    "crypto/sha256"
    "fmt"
)

func SHA256(object interface{}) string {
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf("%v", object)))
    return fmt.Sprintf("%x", h.Sum(nil))
}
