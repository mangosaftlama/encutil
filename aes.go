package encutil

import (
	"crypto/aes"
	"crypto/cipher"
  "crypto/rand"
	"errors"
	"io"
)

type AESKey struct {
  key []byte
  cipher cipher.Block
  gcm cipher.AEAD
  nonce []byte
}

func (key *AESKey) Encrypt(data []byte) []byte {
  return key.gcm.Seal(key.nonce, key.nonce, data, nil)
}

func (key *AESKey) Decrypt(data []byte) ([]byte, error) {
  return key.gcm.Open(nil, key.nonce, data, nil)
}

func NewAESKey(key []byte) (*AESKey, error) {
  if len(key) != 16 && len(key) != 24 && len(key) != 32 {
    return nil, errors.New("Invalid key length. It should be 16, 24 or 32 bytes.")
  }
  
  aesKey := &AESKey{
    key: make([]byte, len(key)), 
  }

  copy(aesKey.key, key)

  c, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }
  aesKey.cipher = c
  
  gcm, err := cipher.NewGCM(c)
  if err != nil {
    return nil, err
  }
  aesKey.gcm = gcm
  
  nonce := make([]byte, gcm.NonceSize())

  if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
    return nil, err
  }
  aesKey.nonce = nonce

  return aesKey, nil
}
