package logical

const (
	andOp uint8 = iota
	xorOp
	orOp
)

func logicalOperation(operation uint8) AfterReadFn {
	return AfterReadFn(func(value byte, cpu LogicalCPU) error {
		A := cpu.GetRegister(RegisterA)
		switch operation {
		case andOp:
			value &= A
		case xorOp:
			value ^= A
		case orOp:
			value |= A
		}
		cpu.SetStatus(value == 0, ZeroFlagBit)
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetRegister(value, RegisterA)
		return nil
	})
}

var and = InstructionDescription{
	Name:    "AND",
	SubExec: logicalOperation(andOp),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x29: Immediate,
		0x25: ZeroPage,
		0x35: ZeroPageX,
		0x2D: Absolute,
		0x3D: AbsoluteX,
		0x39: AbsoluteY,
		0x21: IndirectX,
		0x31: IndirectY,
	},
}

var eor = InstructionDescription{
	Name:    "EOR",
	SubExec: logicalOperation(xorOp),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x49: Immediate,
		0x45: ZeroPage,
		0x55: ZeroPageX,
		0x4D: Absolute,
		0x5D: AbsoluteX,
		0x59: AbsoluteY,
		0x41: IndirectX,
		0x51: IndirectY,
	},
}

var ora = InstructionDescription{
	Name:    "ORA",
	SubExec: logicalOperation(orOp),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x09: Immediate,
		0x05: ZeroPage,
		0x15: ZeroPageX,
		0x0D: Absolute,
		0x1D: AbsoluteX,
		0x19: AbsoluteY,
		0x01: IndirectX,
		0x11: IndirectY,
	},
}

var bit = InstructionDescription{
	Name: "BIT",
	SubExec: AfterReadFn(func(value byte, cpu LogicalCPU) error {
		anded := value & cpu.GetRegister(RegisterA)
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetStatus(anded == 0, ZeroFlagBit)
		cpu.SetStatus(value&0x40 != 0, OverflowFlagBit)
		return nil
	}),
	Access: Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x24: ZeroPage,
		0x2C: Absolute,
	},
}
