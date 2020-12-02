package sunspec_test

import (
	"encoding/binary"
	"fmt"
	"github.com/orlopau/go-energy/pkg/sunspec"
	"github.com/phayes/freeport"
	"github.com/tbrandon/mbserver"
	"testing"
)

func setupModbusServer() (*mbserver.Server, string, error) {
	p, err := freeport.GetFreePort()
	if err != nil {
		return nil, "", err
	}

	server := mbserver.NewServer()
	addr := fmt.Sprintf("127.0.0.1:%v", p)
	err = server.ListenTCP(addr)
	if err != nil {
		return nil, "", err
	}

	return server, addr, nil
}

// setupRegisters sets the registers to test the SunSpec protocol using a real modbus connection.
func setupRegisters(r []uint16) {
	sunID := make([]byte, 4)
	binary.BigEndian.PutUint32(sunID, 0x53756e53)
	r[0] = binary.BigEndian.Uint16(sunID[0:2])
	r[1] = binary.BigEndian.Uint16(sunID[2:4])

	r[2] = 1
	r[3] = 66
	r[70] = 11
	r[71] = 13
	r[72] = 1337

	for i := 71 + 13; i < len(r); i++ {
		r[i] = ^uint16(0)
	}
}

func TestConnect_Scan(t *testing.T) {
	server, addr, err := setupModbusServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	setupRegisters(server.HoldingRegisters)

	ss, err := sunspec.Connect(1, addr)
	if err != nil {
		t.Fatal(err)
	}

	if ss == nil {
		t.Fatal("device is nil")
	}

	hasModel, err := ss.HasModel(11)
	if err != nil {
		t.Fatal(err)
	}

	if !hasModel {
		t.Fatalf("should have model")
	}

	p, err := ss.ReadPointUint16(11, 2)
	if err != nil {
		t.Fatal(err)
	}

	if 1337 != p {
		t.Fatalf("want %v, got %v", 1337, p)
	}
}

func TestConnect_Err(t *testing.T) {
	_, err := sunspec.Connect(1, "127.0.0.2:1234")
	if err == nil {
		t.Fatalf("expected error")
	}
}
