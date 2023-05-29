package logical

func shiftUpdate(left bool) OperationRMWFn {
	return OperationRMWFn(func(value byte, cpu LogicalCPU) (byte, error) {
		newValue := byte(0)
		if left {
			newValue = value << 1
			cpu.SetStatus(value&0x80 != 0, CarryFlagBit)
		} else {
			newValue = value >> 1
			cpu.SetStatus(value&0x1 != 0, CarryFlagBit)
		}
		cpu.SetStatus(newValue == 0, ZeroFlagBit)
		cpu.SetStatus(newValue&0x80 != 0, NegativeFlagBit)
		cpu.SetRegister(newValue, RegisterA)
		return newValue, nil
	})
}

var (
	shiftUpdateLeft  = shiftUpdate(true)
	shiftUpdateRight = shiftUpdate(false)
)

var asl = InstructionDescription{
	Name:    "ASL",
	SubExec: shiftUpdateLeft,
	Access:  ReadModifyWrite,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x0A: Accumulator,
		0x06: ZeroPage,
		0x16: ZeroPageX,
		0x0E: Absolute,
		0x1E: AbsoluteX,
	},
}

var lsr = InstructionDescription{
	Name:    "LSR",
	SubExec: shiftUpdateRight,
	Access:  ReadModifyWrite,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x4A: Accumulator,
		0x46: ZeroPage,
		0x56: ZeroPageX,
		0x4E: Absolute,
		0x5E: AbsoluteX,
	},
}

var rol = InstructionDescription{
	Name:    "ROL",
	SubExec: shiftUpdateLeft,
	Access:  ReadModifyWrite,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x2A: Accumulator,
		0x26: ZeroPage,
		0x36: ZeroPageX,
		0x2E: Absolute,
		0x3E: AbsoluteX,
	},
}

var ror = InstructionDescription{
	Name:    "ROR",
	SubExec: shiftUpdateRight,
	Access:  ReadModifyWrite,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x6A: Accumulator,
		0x66: ZeroPage,
		0x76: ZeroPageX,
		0x6E: Absolute,
		0x7E: AbsoluteX,
	},
}
