package logical

import (
	"bbc/utils"
)

type ReadFn func(LogicalCPU, LogicalBus) (byte, error)
type WriteFn func(byte, LogicalCPU, LogicalBus) error
type ReadModifyWriteFn func(OperationRMWFn, LogicalCPU, LogicalBus) error
type BranchFn func(TakeBranchFn, LogicalCPU, LogicalBus) error
type JumpFn func(LogicalCPU, LogicalBus) (uint16, error)

type AddressingMode uint8
type AccessMode uint8

const (
	Implied AddressingMode = iota
	Accumulator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Relative
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
	NbAddressingMode
)

const (
	Read AccessMode = iota
	Write
	ReadModifyWrite
	ImpliedAccess
	RelativeAccess
	JumpAccess
)

// 2 cycles
func readZeroPageOffset(base, offset uint8, cpu LogicalCPU, bus LogicalBus) (byte, uint16, error) {
	// 6502 performs a read at base, unused but makes the clocks tick
	if err := bus.Tick(); err != nil {
		return 0, 0, err
	}
	addr := utils.SamePageOffset(uint16(base), offset)
	value, err := bus.DirectRead(addr)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

// 4 cycles
func readIndexedIndirectOffset(ptr, offset uint8, cpu LogicalCPU, bus LogicalBus) (byte, uint16, error) {
	// 6502 performs a read at ptr, unused but makes the clocks tick
	if err := bus.Tick(); err != nil {
		return 0, 0, err
	}
	low, err := bus.DirectRead(uint16(ptr + offset))
	if err != nil {
		return 0, 0, err
	}
	high, err := bus.DirectRead(uint16(ptr + offset + 1))
	if err != nil {
		return 0, 0, err
	}
	addr := utils.AddressFromNibbles(high, low)
	value, err := bus.DirectRead(addr)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

// 3 cycles, +1 if page crossed or forced
func readIndirectIndexedOffset(ptr, offset uint8, forceFixing bool, cpu LogicalCPU, bus LogicalBus) (byte, uint16, error) {
	low, err := bus.DirectRead(uint16(ptr))
	if err != nil {
		return 0, 0, err
	}
	high, err := bus.DirectRead(uint16(ptr + 1))
	if err != nil {
		return 0, 0, err
	}
	value, addr, err := bus.OffsetRead(utils.AddressFromNibbles(high, low), offset, forceFixing)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

// 2 cycles
func writeZeroPageOffset(value byte, base, offset uint8, cpu LogicalCPU, bus LogicalBus) error {
	// 6502 performs a read at base, unused but makes the clocks tick
	if err := bus.Tick(); err != nil {
		return err
	}
	return bus.DirectWrite(value, utils.SamePageOffset(uint16(base), offset))
}

// 4 cycles
func writeIndexedIndirectOffset(value byte, ptr, offset uint8, cpu LogicalCPU, bus LogicalBus) error {
	// 6502 performs a read at ptr, unused but makes the clocks tick
	if err := bus.Tick(); err != nil {
		return err
	}
	low, err := bus.DirectRead(uint16(ptr + offset))
	if err != nil {
		return err
	}
	high, err := bus.DirectRead(uint16(ptr + offset + 1))
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, utils.AddressFromNibbles(high, low))
}

// 4 cycles
func writeIndirectIndexedOffset(value byte, ptr, offset uint8, cpu LogicalCPU, bus LogicalBus) (uint16, error) {
	low, err := bus.DirectRead(uint16(ptr))
	if err != nil {
		return 0, err
	}
	high, err := bus.DirectRead(uint16(ptr + 1))
	if err != nil {
		return 0, err
	}
	// one tick to fix the high byte of the address
	if err := bus.Tick(); err != nil {
		return 0, err
	}
	addr := utils.AddressFromNibbles(high, low) + uint16(offset)
	return addr, bus.DirectWrite(value, addr)
}

// 1 cycle
var immediateRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	value, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 3 cycles
var absoluteRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	value, err := bus.DirectRead(addr)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 3 cycles
var absoluteWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, addr)
})

// 5 cycles
var absoluteRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	value, err := bus.DirectRead(addr)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	return bus.DirectWrite(newValue, addr)
})

// 2 cycles
var zeroPageRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	value, err := bus.DirectRead(uint16(addr))
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 2 cycles
var zeroPageWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, uint16(addr))
})

// 4 cycles
var zeroPageRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, err := bus.DirectRead(uint16(addr))
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	return bus.DirectWrite(newValue, uint16(addr))
})

// 3 cycles
var zeroPageXRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	value, _, err := readZeroPageOffset(addr, cpu.GetRegister(RegisterX), cpu, bus)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 3 cycles
var zeroPageXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeZeroPageOffset(value, addr, cpu.GetRegister(RegisterX), cpu, bus)
})

// 5 cycles
var zeroPageXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, _, err := readZeroPageOffset(addr, cpu.GetRegister(RegisterX), cpu, bus)
	if err != nil {
		return err
	}
	// one tick to do the operation, counted in writeZeroPageOffset
	newValue, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	return writeZeroPageOffset(newValue, addr, cpu.GetRegister(RegisterX), cpu, bus)
})

// 3 cycles
var zeroPageYRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	value, _, err := readZeroPageOffset(addr, cpu.GetRegister(RegisterY), cpu, bus)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 3 cycles
var zeroPageYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeZeroPageOffset(value, addr, cpu.GetRegister(RegisterY), cpu, bus)
})

// 3 cycles, +1 if page crossed
var absoluteXRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	value, _, err := bus.OffsetRead(addr, cpu.GetRegister(RegisterX), false)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 4 cycles
var absoluteXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	// one tick to fix the high byte of the address
	if err := bus.Tick(); err != nil {
		return err
	}
	effectiveAddr := addr + uint16(cpu.GetRegister(RegisterX))
	return bus.DirectWrite(value, effectiveAddr)
})

// 6 cycles
var absoluteXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	// one tick to fix the high byte of the address
	if err := bus.Tick(); err != nil {
		return err
	}
	effectiveAddr := addr + uint16(cpu.GetRegister(RegisterX))
	value, err := bus.DirectRead(addr)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	// don't do the boundary check again
	return bus.DirectWrite(newValue, effectiveAddr)
})

// 3 cycles, +1 if page crossed
var absoluteYRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	value, _, err := bus.OffsetRead(addr, cpu.GetRegister(RegisterY), false)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 4 cycles
var absoluteYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	// one tick to fix the high byte of the address
	if err := bus.Tick(); err != nil {
		return err
	}
	effectiveAddr := addr + uint16(cpu.GetRegister(RegisterY))
	return bus.DirectWrite(value, effectiveAddr)
})

// 1 cycle if branch not taken
// 3 cycle if taken, +1 if page crossed
var relativeFn = BranchFn(func(take TakeBranchFn, cpu LogicalCPU, bus LogicalBus) error {
	operand, err := cpu.NextByte()
	if err != nil {
		return err
	}
	takeBranch, err := take(cpu, bus)
	if err != nil {
		return err
	}
	if !takeBranch {
		return nil
	}

	pcl := cpu.GetRegister(RegisterPCL)
	pch := cpu.GetRegister(RegisterPCH)
	pc := utils.AddressFromNibbles(pch, pcl)

	// 1 cycle to add operance
	if err := bus.Tick(); err != nil {
		return err
	}

	// 1 cycle to fix PCH
	if err := bus.Tick(); err != nil {
		return err
	}

	if utils.IsPageCrossed(pc, operand) {
		if err := bus.Tick(); err != nil {
			return err
		}
	}

	pch, pcl = utils.AddressToNibbles(pc + uint16(operand))
	cpu.SetRegister(pcl, RegisterPCL)
	cpu.SetRegister(pch, RegisterPCH)
	return nil
})

// 5 cycles
var indirectXRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	ptr, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	value, _, err := readIndexedIndirectOffset(ptr, cpu.GetRegister(RegisterX), cpu, bus)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 5 cycles
var indirectXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeIndexedIndirectOffset(value, ptr, cpu.GetRegister(RegisterX), cpu, bus)
})

// 7 cycles
var indirectXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, addr, err := readIndexedIndirectOffset(ptr, cpu.GetRegister(RegisterX), cpu, bus)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	// don't do the boundary check again
	return bus.DirectWrite(newValue, addr)
})

// 4 cycles, +1 if page crossed
var indirectYRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	ptr, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	value, _, err := readIndirectIndexedOffset(ptr, cpu.GetRegister(RegisterY), false, cpu, bus)
	if err != nil {
		return 0, err
	}
	return value, nil
})

// 5 cycles
var indirectYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	_, err = writeIndirectIndexedOffset(value, ptr, cpu.GetRegister(RegisterY), cpu, bus)
	return err
})

// 7 cycles
var indirectYRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, addr, err := readIndirectIndexedOffset(ptr, cpu.GetRegister(RegisterY), true, cpu, bus)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	// don't do the boundary check again
	return bus.DirectWrite(newValue, addr)
})

// 1 cycle
var accumulatorRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	value := cpu.GetRegister(RegisterA)
	newA, err := operation(value, cpu, bus)
	if err != nil {
		return err
	}
	cpu.SetRegister(newA, RegisterA)
	// read next instruction byte (and throw it away)
	return bus.Tick()
})

// 1 cycle
var impliedFn = ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
	// read next instruction byte (and throw it away)
	return bus.Tick()
})

// 2 cycles
var absoluteJmp = JumpFn(func(cpu LogicalCPU, bus LogicalBus) (uint16, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, nil
	}
	return addr, nil
})

// 4 cycles
var indirectJmp = JumpFn(func(cpu LogicalCPU, bus LogicalBus) (uint16, error) {
	ptr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	pcl, err := bus.DirectRead(ptr)
	if err != nil {
		return 0, err
	}
	pch, err := bus.DirectRead(utils.SamePageOffset(ptr, 1))
	if err != nil {
		return 0, err
	}
	pc := utils.AddressFromNibbles(pch, pcl)
	return pc, nil
})

var AddressModeFetch = map[AccessMode]map[AddressingMode]interface{}{
	Read: {
		Immediate: immediateRead,
		ZeroPage:  zeroPageRead,
		ZeroPageX: zeroPageXRead,
		ZeroPageY: zeroPageYRead,
		Absolute:  absoluteRead,
		AbsoluteX: absoluteXRead,
		AbsoluteY: absoluteYRead,
		IndirectX: indirectXRead,
		IndirectY: indirectYRead,
	},
	Write: {
		ZeroPage:  zeroPageWrite,
		ZeroPageX: zeroPageXWrite,
		ZeroPageY: zeroPageYWrite,
		Absolute:  absoluteWrite,
		AbsoluteX: absoluteXWrite,
		AbsoluteY: absoluteYWrite,
		IndirectX: indirectXWrite,
		IndirectY: indirectYWrite,
	},
	ReadModifyWrite: {
		Accumulator: accumulatorRMW,
		ZeroPage:    zeroPageRMW,
		ZeroPageX:   zeroPageXRMW,
		Absolute:    absoluteRMW,
		AbsoluteX:   absoluteXRMW,
		IndirectX:   indirectXRMW,
		IndirectY:   indirectYRMW,
	},
	ImpliedAccess: {
		Implied: impliedFn,
	},
	RelativeAccess: {
		Relative: relativeFn,
	},
	JumpAccess: {
		Absolute: absoluteJmp,
		Indirect: indirectJmp,
	},
}
