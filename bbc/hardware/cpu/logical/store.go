package logical

var STA = InstructionDescription{
	Name:    "STA",
	SubExec: storeFrom(RegisterA),
	Access:  Write,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x99: AbsoluteY,
		0x85: ZeroPage,
		0x95: ZeroPageX,
		0x8D: Absolute,
		0x9D: AbsoluteX,
		0x81: IndirectX,
		0x91: IndirectY,
	},
}

var STX = InstructionDescription{
	Name:    "STX",
	SubExec: storeFrom(RegisterX),
	Access:  Write,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x86: ZeroPage,
		0x96: ZeroPageY,
		0x8E: Absolute,
	},
}

var STY = InstructionDescription{
	Name:    "STY",
	SubExec: storeFrom(RegisterX),
	Access:  Write,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x84: ZeroPage,
		0x94: ZeroPageX,
		0x8C: Absolute,
	},
}
