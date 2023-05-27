package logical

import "bbc/hardware"

type StatusFlag uint32

const (
	CarryFlagBit StatusFlag = iota
	ZeroFlagBit
	InterruptDisableFlagBit
	DecimalModeFlagBit
	BreakFlagBit
	UnusedFlagBit
	OverflowFlagBit
	NegativeFlagBit
)

type Register uint8

const (
	RegisterA Register = iota
	RegisterX
	RegisterY
	RegisterStack
	RegisterStatus
)

type LogicalCPU interface {
	GetInstruction(string) *Instruction
	GetInstructionByOpcode(Opcode) *Instruction
	SetInstruction(*Instruction) error
	SetRegister(byte, Register)
	GetRegister(Register) byte
	UpdateStatus(byte, ...StatusFlag)

	Push(byte)
	Pop() byte

	NextByte() (byte, error)
	NextWord() (uint16, error)
}

func loadTo(register Register) AfterReadFn {
	return AfterReadFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
		cpu.SetRegister(value, register)
		cpu.UpdateStatus(value, ZeroFlagBit, NegativeFlagBit)
		return nil
	})
}

func storeFrom(register Register) BeforeWriteFn {
	return BeforeWriteFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
		return cpu.GetRegister(register), nil
	})
}

func transfer(from, to Register) ExecFn {
	return ExecFn(func(cpu LogicalCPU, bus *hardware.Bus) error {
		value := cpu.GetRegister(from)
		cpu.SetRegister(value, to)
		cpu.UpdateStatus(value, ZeroFlagBit, NegativeFlagBit)
		return nil
	})
}

func pushRegister(register Register) ExecFn {
	return ExecFn(func(cpu LogicalCPU, bus *hardware.Bus) error {
		value := cpu.GetRegister(register)
		cpu.Push(value)
		return nil
	})
}

func pullRegister(register Register) ExecFn {
	return ExecFn(func(cpu LogicalCPU, bus *hardware.Bus) error {
		value := cpu.Pop()
		cpu.SetRegister(value, register)
		cpu.UpdateStatus(value, ZeroFlagBit, NegativeFlagBit)
		return nil
	})
}
