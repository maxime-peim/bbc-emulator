package logical

var clc = InstructionDescription{
	Name: "CLC",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(false, CarryFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x18: Implied,
	},
}

var cld = InstructionDescription{
	Name: "CLD",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(false, DecimalModeFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xD8: Implied,
	},
}

var cli = InstructionDescription{
	Name: "CLI",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(false, InterruptDisableFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x58: Implied,
	},
}

var clv = InstructionDescription{
	Name: "CLV",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(false, OverflowFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xB8: Implied,
	},
}

var sec = InstructionDescription{
	Name: "SEC",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(true, CarryFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x38: Implied,
	},
}

var sed = InstructionDescription{
	Name: "SED",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(true, CarryFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xF8: Implied,
	},
}

var sei = InstructionDescription{
	Name: "SEI",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		cpu.SetStatus(true, CarryFlagBit)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x78: Implied,
	},
}
