package hardware

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Clock struct {
	Frequency uint64

	cycles           atomic.Uint64
	stopChannel      chan bool
	freqTimer        *time.Ticker
	lastTicks        uint64
	status           bool
	lastTickTime     time.Time
	waitBetweenTicks time.Duration
}

type ClockHandler struct {
	*Clock
}

func (clock *Clock) Tick() error {
	time.Sleep(clock.waitBetweenTicks - time.Since(clock.lastTickTime))
	if clock.cycles.Add(1) == 0 {
		return fmt.Errorf("clock cycles count wrap")
	}
	return nil
}

func (clock Clock) GetCycles() uint64 {
	return clock.cycles.Load()
}

func (clock *Clock) Reset() {
	clock.cycles.Store(0)
}

func (clock *Clock) Start() error {
	if clock.status {
		return fmt.Errorf("clock already started")
	}
	clock.status = true

	go func() {
		for {
			select {
			case <-clock.stopChannel:
				clock.status = false
				return
			case currentTime := <-clock.freqTimer.C:
				currentTicks := clock.GetCycles()
				simulatedFrequency := uint64(float64(currentTicks-clock.lastTicks) / currentTime.Sub(clock.lastTickTime).Seconds())
				fmt.Printf("Simulated frequency: %d Hz\n", simulatedFrequency)
				clock.lastTicks = currentTicks
				clock.lastTickTime = currentTime
			}
		}
	}()

	return nil
}

func (clock *Clock) Pause() {

}

func NewClock(frequency uint64) *Clock {
	return &Clock{
		Frequency: frequency,

		cycles:           atomic.Uint64{},
		stopChannel:      make(chan bool),
		freqTimer:        time.NewTicker(time.Second),
		lastTicks:        0,
		status:           false,
		lastTickTime:     time.Now(),
		waitBetweenTicks: time.Duration(1. / float64(frequency) * float64(time.Second)),
	}
}
