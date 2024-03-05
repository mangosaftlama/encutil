package encutil

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGenerateRSAPrivateKey(t *testing.T) {
  bits := 2048
  privateKey, err := GeneratePrivateKey(bits)
  assert.NoError(t, err, "GeneratePrivateKey should not return an error")
  
  // validate private key
  assert.NotNil(t, privateKey, "Private key should not be nil")
  assert.Equal(t, bits, privateKey.N.BitLen(),"Private key length should match")
}
