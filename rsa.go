package encutil

import (
  "crypto/rsa"
  "crypto/rand"
)

func GeneratePrivateKey(bits int) (privateKey *rsa.PrivateKey, err error) {
  return rsa.GenerateKey(rand.Reader, bits)
}
