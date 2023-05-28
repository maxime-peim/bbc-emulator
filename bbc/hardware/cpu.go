package hardware

import (
	"bbc/logical"
	"bbc/utils"
	"fmt"

	"github.com/kelindar/bitmap"
)

type CPU struct {
	A      uint8
	X      uint8
	Y      uint8
	Status bitmap.Bitmap

	StackPointer   uint8
	ProgramCounter uint16

	instructionSet      map[string]*logical.Instruction
	instructionByOpcode map[logical.Opcode]*logical.Instruction

	bus *Bus
}

func (cpu *CPU) checkBus() {
	if cpu.bus == nil {
		panic(fmt.Errorf("cpu not plugged to bus"))
	}
}

func (cpu *CPU) executeOpcode(opcode logical.Opcode) error {
	instruction, ok := cpu.instructionByOpcode[opcode]
	if !ok {
		return fmt.Errorf("no instruction registered for opcode %x", opcode)
	}
	fmt.Printf("Executing %s instruction (opcode %x)\n", instruction.Name, opcode)
	return instruction.Execute(opcode, cpu, cpu.bus)
}

func (cpu *CPU) GetInstruction(name string) *logical.Instruction {
	instruction, ok := cpu.instructionSet[name]
	if !ok {
		return nil
	}
	return instruction
}

func (cpu *CPU) GetInstructionByOpcode(opcode logical.Opcode) *logical.Instruction {
	instruction, ok := cpu.instructionByOpcode[opcode]
	if !ok {
		return nil
	}
	return instruction
}

func (cpu *CPU) SetInstruction(instruction *logical.Instruction) error {
	if _, ok := cpu.instructionSet[instruction.Name]; ok {
		return fmt.Errorf("instruction %s already registered", instruction.Name)
	}
	cpu.instructionSet[instruction.Name] = instruction
	for _, opcode := range instruction.GetOpcodes() {
		cpu.instructionByOpcode[opcode] = instruction
	}
	return nil
}

func (cpu *CPU) SetRegister(value byte, register logical.Register) {
	switch register {
	case logical.RegisterA:
		cpu.A = value
	case logical.RegisterX:
		cpu.X = value
	case logical.RegisterY:
		cpu.Y = value
	case logical.RegisterPCL:
		cpu.ProgramCounter = cpu.ProgramCounter&0xFF00 | uint16(value)
	case logical.RegisterPCH:
		cpu.ProgramCounter = cpu.ProgramCounter&0x00FF | uint16(value)<<8
	case logical.RegisterStack:
		cpu.StackPointer = value
	case logical.RegisterStatus:
		cpu.Status = bitmap.FromBytes([]byte{value})
	}
}

func (cpu *CPU) GetRegister(readable logical.Register) byte {
	switch readable {
	case logical.RegisterA:
		return cpu.A
	case logical.RegisterX:
		return cpu.X
	case logical.RegisterY:
		return cpu.Y
	case logical.RegisterPCL:
		return byte(cpu.ProgramCounter & 0xFF)
	case logical.RegisterPCH:
		return byte((cpu.ProgramCounter >> 8) & 0xFF)
	case logical.RegisterStack:
		return cpu.StackPointer
	case logical.RegisterStatus:
		return cpu.Status.ToBytes()[0]
	}
	return 0
}

func (cpu *CPU) SetStatus(set bool, flag logical.StatusFlag) {
	if set {
		cpu.Status.Set(uint32(flag))
	} else {
		cpu.Status.Remove(uint32(flag))
	}
}

func (cpu *CPU) GetStatus(flag logical.StatusFlag) bool {
	return cpu.Status.Contains(uint32(flag))
}

func (cpu *CPU) ExecuteNext() error {
	preCycles := cpu.bus.Clock.GetCycles()
	opcode, err := cpu.NextByte()
	if err != nil {
		return err
	}
	if err := cpu.executeOpcode(logical.Opcode(opcode)); err != nil {
		return err
	}
	postCycles := cpu.bus.Clock.GetCycles()
	fmt.Printf("Took %d cycles\n", postCycles-preCycles)
	return nil
}

func (cpu *CPU) GetPC() uint16 { return cpu.ProgramCounter }
func (cpu *CPU) SetPC(pc uint16) {
	cpu.ProgramCounter = pc
}

func (cpu *CPU) Push(value byte) error {
	stackTop := logical.StackSegment.OffsetIn(uint16(cpu.StackPointer))
	cpu.bus.DirectWrite(value, stackTop)
	cpu.StackPointer--
	return nil
}

func (cpu *CPU) Pop() (byte, error) {
	cpu.StackPointer++
	stackTop := logical.StackSegment.OffsetIn(uint16(cpu.StackPointer))
	value, err := cpu.bus.DirectRead(stackTop)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (cpu *CPU) NextByte() (byte, error) {
	cpu.checkBus()
	value, err := cpu.bus.DirectRead(cpu.ProgramCounter)
	if err != nil {
		return 0, err
	}
	cpu.ProgramCounter++
	return value, nil
}

func (cpu *CPU) NextWord() (uint16, error) {
	cpu.checkBus()
	low, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	high, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	return utils.AddressFromNibbles(high, low), nil
}

func (cpu *CPU) GetName() string { return "CPU" }

func (cpu *CPU) Start() error {
	cpu.checkBus()
	return nil
}

func (cpu *CPU) Reset() error {
	return nil
}

func (cpu *CPU) Stop() error {
	return nil
}

func (cpu *CPU) PlugToBus(bus *Bus) {
	cpu.bus = bus
}

func NewCPU() *CPU {
	status := bitmap.Bitmap{}
	status.Grow(8)
	cpu := CPU{
		StackPointer:        uint8(logical.StackSegment.Start & 0xff),
		Status:              status,
		instructionSet:      map[string]*logical.Instruction{},
		instructionByOpcode: map[logical.Opcode]*logical.Instruction{},
	}

	for _, ins := range logical.BaseInstructionSet {
		if err := ins.RegisterTo(&cpu); err != nil {
			fmt.Printf("error while asm registration: %s", err.Error())
		}
	}

	return &cpu
}
