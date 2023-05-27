package tests

import (
	"bbc/hardware"
	"bbc/hardware/cpu"
	"fmt"
	"os"
	"testing"
)

type Context struct {
	bus   *hardware.Bus
	cpu   *cpu.CPU
	clock *hardware.Clock
}

func (ctx *Context) Reset() {
	ctx.bus.Reset()
}

var testCtx Context

func TestMain(m *testing.M) {
	cpu := cpu.NewCPU()
	clock := hardware.NewClock(2e6)
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
