package summary_api

import (
	"golang.org/x/sync/errgroup"
	"sync"
)

type PowerReader interface {
	ReadPower() (uint, error)
}

type BatteryReader interface {
	PowerReader
	ReadSoC() (uint, error)
}

type MeterReader interface {
	ReadGrid() (int, error)
}

type Plant struct {
	PV    []PowerReader
	Bat   BatteryReader
	Meter MeterReader
}

type PlantSummary struct {
	Grid            int
	PV, Bat         uint
	SelfConsumption int
	BatPercentage   uint
}

func fetchSum(readers ...PowerReader) (uint, error) {
	powc := make(chan uint)
	errc := make(chan error)

	quitc := make(chan bool)
	defer func() {
		close(quitc)
	}()

	for _, v := range readers {
		go func(reader PowerReader) {
			power, err := reader.ReadPower()
			if err != nil {
				select {
				case errc <- err:
				case <-quitc:
				}
			} else {
				powc <- power
			}
		}(v)
	}

	var sum uint
	var i int

	for {
		select {
		case err := <-errc:
			return 0, err
		case <-quitc:
			return sum, nil
		case pow := <-powc:
			sum += pow
			i++
			if i == len(readers) {
				return sum, nil
			}
		}
	}
}

func (r *Plant) FetchSummary() (PlantSummary, error) {
	var summary PlantSummary
	var m sync.Mutex

	// wait for em message
	grid, err := r.Meter.ReadGrid()
	if err != nil {
		return PlantSummary{}, err
	}
	summary.Grid = grid

	// fetch SunSpec data
	var g errgroup.Group

	// fetch PV wattage
	g.Go(func() error {
		pv, err := fetchSum(r.PV...)
		if err != nil {
			return err
		}

		m.Lock()
		summary.PV = pv
		m.Unlock()
		return nil
	})

	// fetch battery wattage
	g.Go(func() error {
		power, err := r.Bat.ReadPower()
		if err != nil {
			return err
		}

		m.Lock()
		summary.Bat = power
		m.Unlock()
		return nil
	})

	// fetch battery soc
	g.Go(func() error {
		soc, err := r.Bat.ReadSoC()
		if err != nil {
			return err
		}

		m.Lock()
		summary.BatPercentage = soc
		m.Unlock()
		return nil
	})

	err = g.Wait()
	if err != nil {
		return PlantSummary{}, err
	}

	summary.SelfConsumption = int(summary.PV) + int(summary.Bat) + summary.Grid

	return summary, nil
}
