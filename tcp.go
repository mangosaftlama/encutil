package encutil

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"net"
	"sync"
)

type TCPServer struct {
	ln         net.Listener
	conns      []net.Conn // required for the broadcast function
	handleConn func(conn *DummyConn)

	connsMu sync.Mutex // mutex for the conns list
}

func NewTCPServer(addr string, handleConn func(conn *DummyConn)) (*TCPServer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		ln: ln,
    handleConn: handleConn,
	}, nil
}

func (s *TCPServer) Accept() (conn net.Conn, err error) {
	return s.ln.Accept()
}

func (s *TCPServer) Addr() net.Addr {
	return s.ln.Addr()
}

func (s *TCPServer) Close() error {
	return s.ln.Close()
}

func (s *TCPServer) AcceptClient() error {
	for {
		conn, err := s.Accept()
		if err != nil {
			return err
		}

		s.connsMu.Lock()
		s.conns = append(s.conns, conn)
		s.connsMu.Unlock()

    dummyConn := &DummyConn{
      write: func(self *DummyConn, p []byte) (n int, err error) {return conn.Write(p)},
      read: func(self *DummyConn, p []byte) (n int, err error) {return conn.Read(p)},
      conn: conn,
    }

		go s.handleConn(dummyConn)
	}
}

// server side key exchange
func (server *TCPServer) KeyExchange(conn net.Conn, publicKey *rsa.PublicKey) (clientPublicKey *rsa.PublicKey, err error) {
	// send rsa key
	rawKey := x509.MarshalPKCS1PublicKey(publicKey)
	written, err := conn.Write(rawKey)
	if err != nil {
		return nil, err
	}
	if written != len(rawKey) {
		return nil, &TooFewDataWrittenError{desc: "ServerKeyExchange() failed!"}
	}
	// read public key from client
	reader := bufio.NewReader(conn)

	clientKeyBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	clientKey, err := x509.ParsePKCS1PublicKey(clientKeyBytes)
	if err != nil {
		return nil, err
	}
	return clientKey, nil
}
