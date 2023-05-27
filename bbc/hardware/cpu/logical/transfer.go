package logical

var TAX = InstructionDescription{
	Name:    "TAX",
	SubExec: transfer(RegisterA, RegisterX),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xAA: Implied,
	},
}

var TAY = InstructionDescription{
	Name:    "TAY",
	SubExec: transfer(RegisterA, RegisterY),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xA8: Implied,
	},
}

var TXA = InstructionDescription{
	Name:    "TXA",
	SubExec: transfer(RegisterX, RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x8A: Implied,
	},
}

var TYA = InstructionDescription{
	Name:    "TYA",
	SubExec: transfer(RegisterY, RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x98: Implied,
	},
}
