package encutil

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAESKey(t *testing.T) {
  data := make([]byte, 32)
  n, err := rand.Read(data)
  assert.NoError(t, err, "rand.Read shouldn't cause an error.")
  assert.Equal(t, 32, n, "The buffer should have the length of 32 bytes.")
  _, err = NewAESKey(data[:16])
  assert.NoError(t, err, "The aes key (16 bytes) shouldn't cause an error.")
  _, err = NewAESKey(data[:24])
  assert.NoError(t, err, "The aes key (24 bytes) shouldn't cause an error.")
  _, err = NewAESKey(data[:32])
  assert.NoError(t, err, "The aes key (32 bytes) shouldn't cause an error")
}
