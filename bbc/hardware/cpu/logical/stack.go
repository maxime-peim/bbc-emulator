package logical

var TXS = InstructionDescription{
	Name:    "TXS",
	SubExec: transfer(RegisterX, RegisterStack),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xBA: Implied,
	},
}

var TSX = InstructionDescription{
	Name:    "TSX",
	SubExec: transfer(RegisterStack, RegisterX),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x9A: Implied,
	},
}

var PHA = InstructionDescription{
	Name:    "PHA",
	SubExec: pushRegister(RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x48: Implied,
	},
}

var PHP = InstructionDescription{
	Name:    "PHP",
	SubExec: pushRegister(RegisterStatus),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x08: Implied,
	},
}

var PLA = InstructionDescription{
	Name:    "PLA",
	SubExec: pullRegister(RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x68: Implied,
	},
}

var PLP = InstructionDescription{
	Name:    "PLP",
	SubExec: pullRegister(RegisterStatus),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x28: Implied,
	},
}
