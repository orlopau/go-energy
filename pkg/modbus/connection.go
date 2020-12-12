package modbus

import (
	"io"
	"math"
	"math/rand"
	"net"
	"time"
)

const (
	backoffBase             = time.Second
	backoffMax              = time.Minute * 2
	backoffRandomMultiplier = time.Second
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

	var backoff int64
	var i int

	for {
		conn, err := net.DialTimeout("tcp", c.addr, timeout)
		tcpConn := conn.(*net.TCPConn)
		if err == nil {
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

		if backoff <= backoffMax.Milliseconds() {
			randomOffset := rand.Float64() * float64(backoffRandomMultiplier.Milliseconds())
			newBackoff := math.Pow(float64(backoffBase.Milliseconds()), float64(i)) + randomOffset
			backoff = int64(math.Max(newBackoff, float64(backoffMax.Milliseconds())))
			i++
		}

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
