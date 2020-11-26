package summary_api

import (
	"fmt"
	"go.uber.org/goleak"
	"testing"
)

type dummyBatteryPowerReader struct {
	isErr bool
	power uint
	soc   uint
}

func (r *dummyBatteryPowerReader) ReadSoC() (uint, error) {
	if r.isErr {
		return 0, fmt.Errorf("dummy error soc")
	}
	return r.soc, nil
}

func (r *dummyBatteryPowerReader) ReadPower() (uint, error) {
	if r.isErr {
		return 0, fmt.Errorf("dummy error power")
	}
	return r.power, nil
}

type dummyEnergyMeter struct {
	isErr bool
	grid  int
}

func (d *dummyEnergyMeter) ReadGrid() (int, error) {
	if d.isErr {
		return 0, fmt.Errorf("dummy error meter")
	}
	return d.grid, nil
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func Test_fetchSum(t *testing.T) {
	defer goleak.VerifyNone(t)

	powerReaders := []PowerReader{
		&dummyBatteryPowerReader{
			power: 100,
		},
		&dummyBatteryPowerReader{
			power: 200,
		},
		&dummyBatteryPowerReader{
			power: 300,
		},
	}

	sum, err := fetchSum(powerReaders...)
	if err != nil {
		t.Fatal(err)
	}

	if sum != 600 {
		t.Fatalf("expected sum of 600, got %d", sum)
	}
}

func Test_fetchSum_err(t *testing.T) {
	defer goleak.VerifyNone(t)

	powerReaders := []PowerReader{
		&dummyBatteryPowerReader{
			power: 100,
			isErr: true,
		},
		&dummyBatteryPowerReader{
			power: 200,
			isErr: true,
		},
		&dummyBatteryPowerReader{
			power: 300,
		},
	}

	_, err := fetchSum(powerReaders...)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPlant_FetchSummary(t *testing.T) {
	defer goleak.VerifyNone(t)

	powerReaders := []PowerReader{
		&dummyBatteryPowerReader{
			power: 100,
		},
		&dummyBatteryPowerReader{
			power: 200,
		},
		&dummyBatteryPowerReader{
			power: 300,
		},
	}

	battReader := dummyBatteryPowerReader{
		power: 200,
		soc:   60,
	}

	energyMeter := dummyEnergyMeter{
		grid: -500,
	}

	plant := Plant{
		PV:    powerReaders,
		Bat:   &battReader,
		Meter: &energyMeter,
	}

	summary, err := plant.FetchSummary()
	if err != nil {
		t.Fatal(err)
	}

	expected := PlantSummary{
		Grid:            -500,
		PV:              600,
		Bat:             200,
		SelfConsumption: 300,
		BatPercentage:   60,
	}

	if summary != expected {
		t.Fatalf("expected %v got %v", expected, summary)
	}
}
