package modbus

import (
	"bytes"
	"fmt"
	"github.com/phayes/freeport"
	"go.uber.org/goleak"
	"io"
	"net"
	"sync"
	"testing"
)

func getFreeAddr() string {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(":%v", port)
}

func newTcpServer(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		conn, err := tcpListener.Accept()
		if err != nil {
			panic(err)
		}

		for {
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			_, err = conn.Write([]byte("OK"))
			if err != nil {
				panic(err)
			}
		}
	}()

	return tcpListener, nil
}

func verifyPing(c net.Conn) error {
	_, err := c.Write([]byte("TEST"))
	if err != nil {
		return err
	}

	b := make([]byte, 1024)
	l, err := c.Read(b)
	if err != nil {
		return err
	}

	if bytes.Compare([]byte("OK"), b[:l]) != 0 {
		return fmt.Errorf("ping not equal")
	}

	return nil
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestReconnectingConn_Ping(t *testing.T) {
	addr := getFreeAddr()

	server, err := newTcpServer(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			panic(err)
		}
	}()

	conn, err := newReconnectingConn(server.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	err = verifyPing(conn)
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReconnectingConn_ReconnectInitial(t *testing.T) {
	addr := getFreeAddr()

	conn, err := newReconnectingConn(addr)
	if err != nil {
		t.Fatal(err)
	}

	var wgPing sync.WaitGroup
	var wgServerStart sync.WaitGroup
	wgPing.Add(1)
	wgServerStart.Add(1)
	go func() {
		wgServerStart.Done()
		err = verifyPing(conn)
		if err != nil {
			t.Fatal(err)
		}
		wgPing.Done()
	}()

	wgServerStart.Wait()
	server, err := newTcpServer(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			panic(err)
		}
	}()

	wgPing.Wait()

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}
