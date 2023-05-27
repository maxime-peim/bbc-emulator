package main

import (
	"bbc/hardware"
	"bbc/hardware/cpu"
	"fmt"
	"os"
)

func main() {
	// BBC micro run at 2MHz
	cpu := cpu.NewCPU()
	clock := hardware.NewClock(2e6)
	ram := hardware.NewRAM()

	_, err := hardware.NewBus(clock, cpu, ram)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := cpu.Start(); err != nil {
		fmt.Printf("Error while executing: %v", err)
	}
}
