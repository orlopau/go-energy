package meter

import (
	"bytes"
	"encoding/hex"
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	energyMeter, err := Listen()

	if err != nil {
		t.Fatal(err)
	}

	if energyMeter.Conn == nil {
		t.Fatalf("connection is nil")
	}
}

type dummyEnergyMeterConnection struct {
	msg []byte
}

func (d *dummyEnergyMeterConnection) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	reader := bytes.NewReader(d.msg)
	n, err := reader.Read(b)
	if err != nil {
		return 0, nil, err
	}

	return n, nil, nil
}

func (d *dummyEnergyMeterConnection) Close() error {
	return nil
}

func TestEnergyMeter_Read(t *testing.T) {
	t.Parallel()

	hexes := "534d4100000402a000000001024400106069015d71551764e5bdd84c0001040000000bf70001080000000002f8910910000204000" +
		"0000000000208000000000dcdc5c87800030400000000000003080000000001f123bc00000404000000014e00040800000000016a2919e" +
		"80009040000000c09000908000000000397ab5348000a040000000000000a08000000000e84ed5c50000d0400000003e20015040000001" +
		"0a90015080000000005378ef3c800160400000000000016080000000003c80e74480017040000000000001708000000000105d38438001" +
		"80400000001150018080000000000e27d9960001d0400000010b2001d08000000000578325168001e040000000000001e0800000000042" +
		"7f938c0001f0400000008a50020040000038e1200210400000003e600290400000000000029080000000000d60cdf10002a0400000004e" +
		"6002a080000000009ce538888002b040000000005002b080000000000aac925c0002c040000000000002c08000000000031f7dab000310" +
		"400000000000031080000000000ec5dc47800320400000004e60032080000000009dd81cc70003304000000023c0034040000038fc2003" +
		"50400000003e8003d040000000034003d0800000000013ddc2538003e040000000000003e0800000000048a4abc10003f0400000000000" +
		"03f0800000000005def5bc0004004000000003d0040080000000000731be2e80045040000000050004508000000000181b137d00046040" +
		"000000000004608000000000494ea5428004704000000002300480400000391f80049040000000286900000000200105200000000"

	msg, err := hex.DecodeString(hexes)
	if err != nil {
		t.Fatal(err)
	}

	emConn := &dummyEnergyMeterConnection{msg}
	em := &energyMeter{emConn}

	telegram, err := em.ReadTelegram()
	if err != nil {
		t.Fatal(err)
	}

	if telegram == nil {
		t.Fatalf("telegram is nil")
	}
}
