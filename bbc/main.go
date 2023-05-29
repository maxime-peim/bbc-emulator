package main

import (
	"bbc/hardware"
	"fmt"
	"os"
)

func main() {
	// BBC micro run at 2MHz
	clock := hardware.NewClock(2e6)
	cpu := hardware.NewCPU(clock)
	ram := hardware.NewRAM()

	bus, err := hardware.NewBus(clock, cpu, ram)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := clock.Start(); err != nil {
		fmt.Printf("Error while starting clock: %v", err)
	}

	initialPC := cpu.ProgramCounter

	program := []byte{
		0xA9, 0x55, // LDA #$55
		0x4C, 0x00, 0x00, // JMP 0000
	}
	bus.WriteMultiple(program, initialPC)

	if err := cpu.Start(); err != nil {
		fmt.Printf("Error while executing: %v", err)
	}
}
