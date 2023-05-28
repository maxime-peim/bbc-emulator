package logical

func loadTo(register Register) AfterReadFn {
	return AfterReadFn(func(value byte, cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(value == 0, ZeroFlagBit)
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetRegister(value, register)
		return nil
	})
}

var lda = InstructionDescription{
	Name:    "LDA",
	SubExec: loadTo(RegisterA),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xA9: Immediate,
		0xB9: AbsoluteY,
		0xA5: ZeroPage,
		0xB5: ZeroPageX,
		0xAD: Absolute,
		0xBD: AbsoluteX,
		0xA1: IndirectX,
		0xB1: IndirectY,
	},
}

var ldx = InstructionDescription{
	Name:    "LDX",
	SubExec: loadTo(RegisterX),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xA2: Immediate,
		0xBE: AbsoluteY,
		0xA6: ZeroPage,
		0xAE: Absolute,
		0xB6: ZeroPageY,
	},
}

var ldy = InstructionDescription{
	Name:    "LDY",
	SubExec: loadTo(RegisterY),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xA0: Immediate,
		0xA4: ZeroPage,
		0xB4: ZeroPageX,
		0xAC: Absolute,
		0xBC: AbsoluteX,
	},
}
