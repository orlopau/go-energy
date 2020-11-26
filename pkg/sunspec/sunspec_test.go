package sunspec

import (
	"github.com/goburrow/modbus"
	"os"
	"testing"
)

func TestSunSpecReader_Scan_SMAIntegration(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping this integration test using a real inverter in the CI...")
	}

	handler := modbus.NewTCPClientHandler("192.168.188.34:502")
	handler.SlaveId = 126
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := handler.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	client := modbus.NewClient(handler)

	s := &ModbusModelReader{
		ModbusReader: ModbusReader{
			client,
		},
	}

	err = s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if len(s.Models) == 0 {
		t.Fatalf("no models found")
	}

	t.Log(s.Models)
}
