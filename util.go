package encutil

import (
	"fmt"
	"net"
)

type TooFewDataWrittenError struct {
	desc string
}

func (e TooFewDataWrittenError) Error() string {
	return fmt.Sprintf("failed writing data (too few data written): %s.", e.desc)
}

// DummyConn is a wrapper around the connection so you can treat it like a normal connection and don't have to worry about encryption at all.
type DummyConn struct {
	write func(self *DummyConn, p []byte) (n int, err error)
	read  func(self *DummyConn, p []byte) (n int, err error)

	conn net.Conn
  aesKey AESKey
}

func (c *DummyConn) Write(p []byte) (n int, err error) {
	return c.write(c, c.aesKey.Encrypt(p))
}

func (c *DummyConn) Read(p []byte) (n int, err error) {
  n, err = c.read(c, p)
  if err != nil {
    return n, err
  }
  data, err := c.aesKey.Decrypt(p)
  if err != nil {
    return 0, err
  }
  copy(p, data)
  return len(data), nil
}
