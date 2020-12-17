package modbus

import (
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

const (
	backoffStart            = 500 * time.Millisecond
	backoffExpoBase         = 1.2
	backoffMax              = 30 * time.Second
	backoffRandomMultiplier = 100 * time.Millisecond
	timeout                 = time.Minute
)

// reconnectingConn is a network connection that tries to reconnect after the connection is closed.
type reconnectingConn struct {
	*net.TCPConn
	addr   string
	cancel chan bool
}

func newReconnectingConn(addr string) (*reconnectingConn, error) {
	return &reconnectingConn{
		addr: addr,
	}, nil
}

func (c *reconnectingConn) verifyConn() error {
	if c.TCPConn == nil {
		return c.reconnect()
	}

	if !isConnAlive(c.TCPConn) {
		err := c.reconnect()
		if err != nil {
			return err
		}
	}

	return nil
}

func isConnAlive(c *net.TCPConn) bool {
	err := c.SetReadDeadline(time.Now())
	if err != nil {
		return false
	}

	defer func() {
		var zero time.Time
		err := c.SetReadDeadline(zero)
		if err != nil {
			panic(err)
		}
	}()

	b := make([]byte, 1)
	if _, err := c.Read(b); err == io.EOF {
		return false
	} else {
		return true
	}
}

func (c *reconnectingConn) reconnect() error {
	if c.TCPConn != nil {
		err := c.TCPConn.Close()
		if err != nil {
			return err
		}
	}

	var i int

	for {
		conn, err := net.DialTimeout("tcp", c.addr, timeout)
		if err == nil {
			tcpConn := conn.(*net.TCPConn)
			err := tcpConn.SetKeepAlive(true)
			if err != nil {
				return err
			}
			err = tcpConn.SetKeepAlivePeriod(30 * time.Second)
			if err != nil {
				return err
			}

			c.TCPConn = tcpConn
			return nil
		}

		randomOffset := rand.Float64() * float64(backoffRandomMultiplier.Milliseconds())
		backoff := math.Pow(backoffExpoBase, float64(i))*float64(backoffStart.Milliseconds()) + randomOffset
		if backoff < float64(backoffMax.Milliseconds()) {
			i++
		} else {
			backoff = float64(backoffMax.Milliseconds())
		}

		log.Printf("backing off %v millis...", backoff)

		select {
		case <-c.cancel:
			return nil
		case <-time.After(time.Duration(backoff) * time.Millisecond):
		}
	}
}

func (c *reconnectingConn) Read(p []byte) (n int, err error) {
	err = c.verifyConn()
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Read(p)
}

func (c *reconnectingConn) Write(p []byte) (n int, err error) {
	err = c.verifyConn()
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Write(p)
}

func (c *reconnectingConn) Close() error {
	if c.TCPConn != nil {
		return c.TCPConn.Close()
	}
	close(c.cancel)

	return nil
}
