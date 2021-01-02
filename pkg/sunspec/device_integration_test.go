package sunspec_test

import (
	"fmt"
	"github.com/orlopau/go-energy/pkg/modbus"
	"github.com/orlopau/go-energy/pkg/sunspec"
	"github.com/phayes/freeport"
	"github.com/xiegeo/modbusone"
	"golang.org/x/net/context"
	"math"
	"net"
	"runtime"
	"testing"
	"time"
)

func startSunSpecServer(ctx context.Context, t *testing.T, addr string, value int) {
	t.Logf("modbus server starting on %v", addr)

	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		t.Fatal(err)
	}

	server := modbusone.NewTCPServer(listener)

	mockbus := modbus.NewMockbus(1000)
	err = mockbus.AddHoldingRegisterEntries(map[uint16]interface{}{
		0:  uint32(0x53756e53),
		2:  uint16(1),
		3:  uint16(66),
		70: uint16(2),
		71: uint16(1),
		// sunspec point, model 11, point 3, val 100
		72: uint16(value),
		73: uint16(math.MaxUint16),
		74: uint16(0),
	})

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		select {
		case <-ctx.Done():
			err2 := server.Close()
			if err2 != nil {
				t.Fatal(err2)
			}
		}
	}()

	handler := &modbusone.SimpleHandler{
		ReadHoldingRegisters: func(address, quantity uint16) ([]uint16, error) {
			if ctx.Err() != nil {
				// yikes, used to kill the connection
				runtime.Goexit()
			}
			regs, err2 := mockbus.ReadHoldingRegistersUint(address, quantity)
			return regs, err2
		},
	}

	err = server.Serve(handler)
	if err != nil && ctx.Err() == nil {
		t.Fatal(err)
	}

	t.Log("server terminated")
}

func verifySunSpec(t *testing.T, device *sunspec.ModbusDevice, val int) {
	v, err := device.GetAnyPoint(sunspec.Point{
		Point: 2,
		Model: 2,
		T:     uint16(0),
	})
	if err != nil {
		t.Fatal(err)
	}

	if v != float64(val) {
		t.Fatalf("expected %v, got %v", val, v)
	}
}

func TestConnect_Simple(t *testing.T) {
	t.Parallel()

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	addr := fmt.Sprintf(":%v", port)

	started := make(chan bool)

	go func() {
		started <- true
		startSunSpecServer(context.TODO(), t, addr, 100)
	}()

	<-started

	device, err := sunspec.Connect(addr)
	if err != nil {
		t.Fatal(err)
	}

	verifySunSpec(t, device, 100)
}

func TestConnect_InitialRefused(t *testing.T) {
	t.Parallel()

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	addr := fmt.Sprintf(":%v", port)

	go func() {
		<-time.After(time.Second)
		startSunSpecServer(context.TODO(), t, addr, 100)
	}()

	device, err := sunspec.Connect(addr)
	if err != nil {
		t.Fatal(err)
	}

	verifySunSpec(t, device, 100)
}

func TestConnect_InBetweenRefused(t *testing.T) {
	t.Parallel()

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	addr := fmt.Sprintf(":%v", port)

	started := make(chan bool)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		started <- true
		startSunSpecServer(ctx, t, addr, 100)
	}()

	<-started

	device, err := sunspec.Connect(addr)
	if err != nil {
		t.Fatal(err)
	}

	verifySunSpec(t, device, 100)

	// shutdown device
	cancelFunc()

	const secondVal = 200

	// start again
	go func() {
		<-time.After(1 * time.Second)
		startSunSpecServer(context.TODO(), t, addr, secondVal)
	}()

	// verify while server is down to check retry mechanism
	verifySunSpec(t, device, secondVal)
}
