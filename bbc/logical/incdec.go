package logical

func incdecRegister(register Register, increment bool) ExecFn {
	return ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		value := cpu.GetRegister(register)
		if increment {
			value++
		} else {
			value--
		}
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetStatus(value == 0, ZeroFlagBit)
		cpu.SetRegister(value, register)
		return nil
	})
}

var inc = InstructionDescription{
	Name: "INC",
	SubExec: OperationRMWFn(func(value byte, cpu LogicalCPU, bus LogicalBus) (byte, error) {
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetStatus(value == 0, ZeroFlagBit)
		return value + 1, nil
	}),
	Access: ReadModifyWrite,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xE6: ZeroPage,
		0xF6: ZeroPageX,
		0xEE: Absolute,
		0xFE: AbsoluteX,
	},
}

var inx = InstructionDescription{
	Name:    "INX",
	SubExec: incdecRegister(RegisterX, true),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xE8: Implied,
	},
}

var iny = InstructionDescription{
	Name:    "INY",
	SubExec: incdecRegister(RegisterY, true),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xC8: Implied,
	},
}

var dec = InstructionDescription{
	Name: "DEC",
	SubExec: OperationRMWFn(func(value byte, cpu LogicalCPU, bus LogicalBus) (byte, error) {
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetStatus(value == 0, ZeroFlagBit)
		return value - 1, nil
	}),
	Access: ReadModifyWrite,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xC6: ZeroPage,
		0xD6: ZeroPageX,
		0xCE: Absolute,
		0xDE: AbsoluteX,
	},
}

var dex = InstructionDescription{
	Name:    "DEX",
	SubExec: incdecRegister(RegisterX, false),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xCA: Implied,
	},
}

var dey = InstructionDescription{
	Name:    "DEY",
	SubExec: incdecRegister(RegisterY, false),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x88: Implied,
	},
}
