package discovery

import (
	"encoding/hex"
	"fmt"
	"net"
	"testing"
	"time"
)

type dummyAddr string
type dummyConn struct {
	net.Conn
	nextPacketAddr net.Addr
}
type dummyDuplex struct {
	r dummyConn
	w net.Conn
}
type testData struct {
	addr    net.Addr
	payload string
}

func (d dummyAddr) Network() string { return string(d) }
func (d dummyAddr) String() string  { return string(d) }

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

func (t *dummyDuplex) WriteToFromAddr(p []byte, addr net.Addr) (int, error) {
	t.r.nextPacketAddr = addr
	return t.w.Write(p)
}

func (t *dummyDuplex) Close() error {
	if err := t.r.Close(); err != nil {
		return err
	}
	if err := t.w.Close(); err != nil {
		return err
	}

	return nil
}

func setupDiscoveryTest(data []testData) (duplex *dummyDuplex, err error) {
	r, w := net.Pipe()

	duplex = &dummyDuplex{
		r: dummyConn{
			r,
			nil,
		},
		w: w,
	}

	go func() {
		for _, v := range data {
			bytes, err := hex.DecodeString(v.payload)
			if err != nil {
				panic(err)
			}
			if _, err := duplex.WriteToFromAddr(bytes, v.addr); err != nil {
				panic(err)
			}
		}
	}()

	return duplex, nil
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

	testData := []testData{
		{
			dummyAddr("127.0.0.1"),
			discoveryVerifyPayload,
		},
		{
			dummyAddr("127.0.0.2"),
			discoveryVerifyPayload,
		},
	}

	duplex, err := setupDiscoveryTest(testData)
	if err != nil {
		t.Fatal(err)
	}

	timeout := 100 * time.Millisecond
	addrs, err := listen(&duplex.r, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(testData) != len(addrs) {
		t.Fatalf("expected %d addresses, got %d", len(testData), len(addrs))
	}

	for i := 0; i < len(addrs); i++ {
		if testData[i].addr != addrs[i] {
			t.Fatalf("addresses not equal")
		}
	}

	if err := duplex.Close(); err != nil {
		t.Fatal(err)
	}
}

// TestListen_Addresses verifies that returned packages with invalid payloads are not added to the list.
func TestListen_Invalid_Verify(t *testing.T) {
	t.Parallel()

	testData := []testData{
		{
			dummyAddr("127.0.0.1"),
			"feedbeaf",
		},
		{
			dummyAddr("127.0.0.2"),
			discoveryVerifyPayload,
		},
	}

	duplex, err := setupDiscoveryTest(testData)
	if err != nil {
		t.Fatal(err)
	}

	timeout := 100 * time.Millisecond
	addrs, err := listen(&duplex.r, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("expected 1 addresses, got %d", len(addrs))
	}

	if addrs[0] != testData[1].addr {
		t.Fatalf("incorrect address, expected %s got %s", testData[1].addr, addrs[0])
	}
}

func TestListen_No_Duplicates(t *testing.T) {
	t.Parallel()

	testData := []testData{
		{
			dummyAddr("127.0.0.1"),
			discoveryVerifyPayload,
		},
		{
			dummyAddr("127.0.0.1"),
			discoveryVerifyPayload,
		},
	}

	duplex, err := setupDiscoveryTest(testData)
	if err != nil {
		t.Fatal(err)
	}

	timeout := 100 * time.Millisecond
	addrs, err := listen(&duplex.r, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("expected 1 addresses, got %d", len(addrs))
	}

	if addrs[0] != testData[0].addr {
		t.Fatalf("incorrect address, expected %s got %s", testData[1].addr, addrs[0])
	}
}

func TestDiscoverInverters(t *testing.T) {
	go func() {
		// send on loopback
		verifyPayload, err := hex.DecodeString(discoveryVerifyPayload)
		if err != nil {
			t.Fatal(err)
		}

		conn, err := net.Dial("udp4", "239.12.255.254:9522")
		if err != nil {
			t.Fatal(err)
		}

		_, err = conn.Write(verifyPayload)
		if err != nil {
			t.Fatal(err)
		}
	}()

	addrs, err := DiscoverInverters(nil, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(addrs)

	if len(addrs) != 1 {
		t.Fatal("no addresses found")
	}
}
