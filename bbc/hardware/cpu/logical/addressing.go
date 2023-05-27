package logical

import (
	"bbc/hardware"
	"bbc/utils"
)

type ReadFn func(LogicalCPU, *hardware.Bus) (byte, error)
type WriteFn func(byte, LogicalCPU, *hardware.Bus) error
type ReadModifyWriteFn func(OperationRMWFn, LogicalCPU, *hardware.Bus) error

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
)

func readZeroPageOffset(base, offset uint8, cpu LogicalCPU, bus *hardware.Bus) (byte, uint16, error) {
	// 6502 performs a read at base, unused but makes the clocks tick
	if err := bus.Clock.Tick(); err != nil {
		return 0, 0, err
	}
	addr := (uint16(base) + uint16(offset)) & 0xff
	value, err := bus.DirectRead(addr)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

func readIndexedIndirectOffset(ptr, offset uint8, cpu LogicalCPU, bus *hardware.Bus) (byte, uint16, error) {
	// 6502 performs a read at ptr, unused but makes the clocks tick
	if err := bus.Clock.Tick(); err != nil {
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
	addr := utils.AddressFromNibble(high, low)
	value, err := bus.DirectRead(addr)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

func readIndirectIndexedOffset(ptr, offset uint8, cpu LogicalCPU, bus *hardware.Bus) (byte, uint16, error) {
	low, err := bus.DirectRead(uint16(ptr))
	if err != nil {
		return 0, 0, err
	}
	high, err := bus.DirectRead(uint16(ptr + 1))
	if err != nil {
		return 0, 0, err
	}
	value, addr, err := bus.OffsetRead(utils.AddressFromNibble(high, low), offset)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

func writeZeroPageOffset(value byte, base, offset uint8, cpu LogicalCPU, bus *hardware.Bus) error {
	// 6502 performs a read at base, unused but makes the clocks tick
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	return bus.DirectWrite(value, (uint16(base)+uint16(offset))&0xff)
}

func writeIndexedIndirectOffset(value byte, ptr, offset uint8, cpu LogicalCPU, bus *hardware.Bus) error {
	// 6502 performs a read at ptr, unused but makes the clocks tick
	if err := bus.Clock.Tick(); err != nil {
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
	return bus.DirectWrite(value, utils.AddressFromNibble(high, low))
}

func writeIndirectIndexedOffset(value byte, ptr, offset uint8, cpu LogicalCPU, bus *hardware.Bus) error {
	low, err := bus.DirectRead(uint16(ptr))
	if err != nil {
		return err
	}
	high, err := bus.DirectRead(uint16(ptr + 1))
	if err != nil {
		return err
	}
	return bus.OffsetWrite(value, utils.AddressFromNibble(high, low), offset)
}

var immediateRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
	value, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	return value, nil
})

var absoluteRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var absoluteWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, addr)
})

var absoluteRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	value, err := bus.DirectRead(addr)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value)
	if err != nil {
		return err
	}
	return bus.DirectWrite(newValue, addr)
})

var zeroPageRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var zeroPageWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return bus.DirectWrite(value, uint16(addr))
})

var zeroPageRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, err := bus.DirectRead(uint16(addr))
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value)
	if err != nil {
		return err
	}
	return bus.DirectWrite(newValue, uint16(addr))
})

var zeroPageXRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var zeroPageXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeZeroPageOffset(value, addr, cpu.GetRegister(RegisterX), cpu, bus)
})

var zeroPageXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, _, err := readZeroPageOffset(addr, cpu.GetRegister(RegisterX), cpu, bus)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value)
	if err != nil {
		return err
	}
	return writeZeroPageOffset(newValue, addr, cpu.GetRegister(RegisterX), cpu, bus)
})

var zeroPageYRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var zeroPageYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeZeroPageOffset(value, addr, cpu.GetRegister(RegisterY), cpu, bus)
})

var absoluteXRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var absoluteXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.OffsetWrite(value, addr, cpu.GetRegister(RegisterX))
})

var absoluteXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	value, effectiveAddr, err := bus.OffsetRead(addr, cpu.GetRegister(RegisterX))
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value)
	if err != nil {
		return err
	}
	// don't do the boundary check again
	return bus.DirectWrite(newValue, effectiveAddr)
})

var absoluteYRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var absoluteYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	addr, err := cpu.NextWord()
	if err != nil {
		return err
	}
	return bus.OffsetWrite(value, addr, cpu.GetRegister(RegisterY))
})

var relativeRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
	_, err := cpu.NextByte()
	if err != nil {
		return 0, err
	}
	// TODO: finish
	return 0, nil
})

var indirectXRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var indirectXWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeIndexedIndirectOffset(value, ptr, cpu.GetRegister(RegisterX), cpu, bus)
})

var indirectXRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus *hardware.Bus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, addr, err := readIndexedIndirectOffset(ptr, cpu.GetRegister(RegisterX), cpu, bus)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value)
	if err != nil {
		return err
	}
	// don't do the boundary check again
	return bus.DirectWrite(newValue, addr)
})

var indirectYRead = ReadFn(func(cpu LogicalCPU, bus *hardware.Bus) (byte, error) {
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

var indirectYWrite = WriteFn(func(value byte, cpu LogicalCPU, bus *hardware.Bus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	return writeIndirectIndexedOffset(value, ptr, cpu.GetRegister(RegisterY), cpu, bus)
})

var indirectYRMW = ReadModifyWriteFn(func(operation OperationRMWFn, cpu LogicalCPU, bus *hardware.Bus) error {
	ptr, err := cpu.NextByte()
	if err != nil {
		return err
	}
	value, addr, err := readIndirectIndexedOffset(ptr, cpu.GetRegister(RegisterY), cpu, bus)
	if err != nil {
		return err
	}
	// one tick to do the operation
	if err := bus.Clock.Tick(); err != nil {
		return err
	}
	newValue, err := operation(value)
	if err != nil {
		return err
	}
	// don't do the boundary check again
	return bus.DirectWrite(newValue, addr)
})

var AddressModeFetch = map[AccessMode]map[AddressingMode]interface{}{
	Read: {
		Immediate: immediateRead,
		ZeroPage:  zeroPageRead,
		ZeroPageX: zeroPageXRead,
		ZeroPageY: zeroPageYRead,
		Relative:  relativeRead,
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
		ZeroPage:  zeroPageRMW,
		ZeroPageX: zeroPageXRMW,
		Absolute:  absoluteRMW,
		AbsoluteX: absoluteXRMW,
		IndirectX: indirectXRMW,
		IndirectY: indirectYRMW,
	},
	ImpliedAccess: {},
}
