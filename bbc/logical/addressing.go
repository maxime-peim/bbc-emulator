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

func readIndirectIndexedOffset(ptr, offset uint8, cpu LogicalCPU, bus LogicalBus) (byte, uint16, error) {
	low, err := bus.DirectRead(uint16(ptr))
	if err != nil {
		return 0, 0, err
	}
	high, err := bus.DirectRead(uint16(ptr + 1))
	if err != nil {
		return 0, 0, err
	}
	value, addr, err := bus.OffsetRead(utils.AddressFromNibbles(high, low), offset)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

func writeZeroPageOffset(value byte, base, offset uint8, cpu LogicalCPU, bus LogicalBus) error {
	// 6502 performs a read at base, unused but makes the clocks tick
	if err := bus.Tick(); err != nil {
		return err
	}
	return bus.DirectWrite(value, utils.SamePageOffset(uint16(base), offset))
}

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

func writeIndirectIndexedOffset(value byte, ptr, offset uint8, cpu LogicalCPU, bus LogicalBus) error {
	low, err := bus.DirectRead(uint16(ptr))
	if err != nil {
		return err
	}
	high, err := bus.DirectRead(uint16(ptr + 1))
	if err != nil {
		return err
	}
	return bus.OffsetWrite(value, utils.AddressFromNibbles(high, low), offset)
}

var immediateRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	value, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	return value, nil
})

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

var absoluteWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, addr)
})

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

var zeroPageWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, uint16(addr))
})

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

var zeroPageXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeZeroPageOffset(value, addr, cpu.GetRegister(RegisterX), cpu, bus)
})

var zeroPageXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, _, err := readZeroPageOffset(addr, cpu.GetRegister(RegisterX), cpu, bus)
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
	return writeZeroPageOffset(newValue, addr, cpu.GetRegister(RegisterX), cpu, bus)
})

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

var zeroPageYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeZeroPageOffset(value, addr, cpu.GetRegister(RegisterY), cpu, bus)
})

var absoluteXRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	value, _, err := bus.OffsetRead(addr, cpu.GetRegister(RegisterX))
	if err != nil {
		return 0, err
	}
	return value, nil
})

var absoluteXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.OffsetWrite(value, addr, cpu.GetRegister(RegisterX))
})

var absoluteXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	value, effectiveAddr, err := bus.OffsetRead(addr, cpu.GetRegister(RegisterX))
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

var absoluteYRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	value, _, err := bus.OffsetRead(addr, cpu.GetRegister(RegisterY))
	if err != nil {
		return 0, err
	}
	return value, nil
})

var absoluteYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.OffsetWrite(value, addr, cpu.GetRegister(RegisterY))
})

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
	if err := bus.Tick(); err != nil {
		return err
	}
	pcl := cpu.GetRegister(RegisterPCL)
	pch := cpu.GetRegister(RegisterPCH)
	pc := utils.AddressFromNibbles(pch, pcl)
	if utils.IsPageCrossed(pc, operand) {
		if err := bus.Tick(); err != nil {
			return err
		}
	}
	return nil
})

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

var indirectXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeIndexedIndirectOffset(value, ptr, cpu.GetRegister(RegisterX), cpu, bus)
})

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

var indirectYRead = ReadFn(func(cpu LogicalCPU, bus LogicalBus) (byte, error) {
	ptr, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	value, _, err := readIndirectIndexedOffset(ptr, cpu.GetRegister(RegisterY), cpu, bus)
	if err != nil {
		return 0, err
	}
	return value, nil
})

var indirectYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeIndirectIndexedOffset(value, ptr, cpu.GetRegister(RegisterY), cpu, bus)
})

var indirectYRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, addr, err := readIndirectIndexedOffset(ptr, cpu.GetRegister(RegisterY), cpu, bus)
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

var accumulatorRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus LogicalBus) error {
	// read next instruction byte (and throw it away)
	return bus.Tick()
})

var impliedFn = ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
	// read next instruction byte (and throw it away)
	return bus.Tick()
})

var absoluteJmp = JumpFn(func(cpu LogicalCPU, bus LogicalBus) (uint16, error) {
	addr, err := cpu.NextWord()
	if err != nil {
		return 0, nil
	}
	return addr, nil
})

var indirectJmp = JumpFn(func(cpu LogicalCPU, bus LogicalBus) (uint16, error) {
	ptr, err := cpu.NextWord()
	if err != nil {
		return 0, err
	}
	low, err := bus.DirectRead(ptr)
	if err != nil {
		return 0, err
	}
	high, err := bus.DirectRead(utils.SamePageOffset(ptr, 1))
	if err != nil {
		return 0, err
	}
	addr := utils.AddressFromNibbles(high, low)
	return addr, nil
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
