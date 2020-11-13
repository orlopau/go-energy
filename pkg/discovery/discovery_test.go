package discovery

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type dummyConn struct {
	net.Conn
	nextPacketAddr net.Addr
}

func (t *dummyConn) ReadFrom(p []byte) (int, net.Addr, error) {
	n, err := t.Read(p)
	if err != nil {
		return 0, nil, err
	}
	if t.nextPacketAddr == nil {
		return 0, nil, fmt.Errorf("no address specified")
	}
	return n, t.nextPacketAddr, nil
}

func (t *dummyConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	panic("Not implemented!")
}

type dummyAddr string

func (d dummyAddr) Network() string { return string(d) }
func (d dummyAddr) String() string  { return string(d) }

type dummyDuplex struct {
	r dummyConn
	w net.Conn
}

func (t *dummyDuplex) WriteToFromAddr(p []byte, addr net.Addr) (int, error) {
	t.r.nextPacketAddr = addr
	return t.w.Write(p)
}

// TestListen_Timeout verifies that the connection timeouts after the specified amount.
func TestListen_Timeout(t *testing.T) {
	t.Parallel()

	conn, _ := net.Pipe()
	testConn := &dummyConn{
		conn,
		dummyAddr("127.0.0.1"),
	}

	timeout := 100 * time.Millisecond
	start := time.Now()

	_, err := listen(testConn, timeout)
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Now().Sub(start)
	expected := timeout + 5*time.Millisecond
	if duration > expected {
		t.Fatalf("should have timed out after %s but did after %s", expected, duration)
	}
}

// TestListen_Addresses verifies that the correct addresses are returned.
func TestListen_Addresses(t *testing.T) {
	t.Parallel()

	r, w := net.Pipe()
	defer func() {
		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	duplex := &dummyDuplex{
		r: dummyConn{
			r,
			nil,
		},
		w: w,
	}

	exAddrs := []net.Addr{
		dummyAddr("127.0.0.1"),
		dummyAddr("127.0.0.2"),
		dummyAddr("127.0.0.3"),
	}
	go func() {
		for _, v := range exAddrs {
			if _, err := duplex.WriteToFromAddr([]byte("msg"), v); err != nil {
				t.Fatal(err)
			}
		}
	}()

	timeout := 100 * time.Millisecond
	addrs, err := listen(&duplex.r, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(exAddrs) != len(addrs) {
		t.Fatalf("expected %d addresses, got %d", len(exAddrs), len(addrs))
	}

	for i := 0; i < len(addrs); i++ {
		if exAddrs[i] != addrs[i] {
			t.Fatalf("addresses not equal")
		}
	}
}
