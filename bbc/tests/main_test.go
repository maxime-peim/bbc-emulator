package tests

import (
	"bbc/hardware"
	"fmt"
	"os"
	"testing"
)

type Context struct {
	bus   *hardware.Bus
	cpu   *hardware.CPU
	clock *hardware.Clock
}

func (ctx *Context) Reset() {
	ctx.bus.Reset()
}

var testCtx Context

func TestMain(m *testing.M) {
	clock := hardware.NewClock(2e6)
	cpu := hardware.NewCPU(clock)
	ram := hardware.NewRAM()

	bus, err := hardware.NewBus(clock, cpu, ram)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	testCtx = Context{
		bus:   bus,
		cpu:   cpu,
		clock: clock,
	}

	os.Exit(m.Run())
}
