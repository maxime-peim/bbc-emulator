package logical

func transfer(from, to Register) ExecFn {
	return ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		value := cpu.GetRegister(from)
		cpu.SetStatus(value == 0, ZeroFlagBit)
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetRegister(value, to)
		return nil
	})
}

var txs = InstructionDescription{
	Name:    "TXS",
	SubExec: transfer(RegisterX, RegisterStack),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xBA: Implied,
	},
}

var tsx = InstructionDescription{
	Name:    "TSX",
	SubExec: transfer(RegisterStack, RegisterX),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x9A: Implied,
	},
}

var tax = InstructionDescription{
	Name:    "TAX",
	SubExec: transfer(RegisterA, RegisterX),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xAA: Implied,
	},
}

var tay = InstructionDescription{
	Name:    "TAY",
	SubExec: transfer(RegisterA, RegisterY),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xA8: Implied,
	},
}

var txa = InstructionDescription{
	Name:    "TXA",
	SubExec: transfer(RegisterX, RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x8A: Implied,
	},
}

var tya = InstructionDescription{
	Name:    "TYA",
	SubExec: transfer(RegisterY, RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x98: Implied,
	},
}
