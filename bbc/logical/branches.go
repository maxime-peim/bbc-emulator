package logical

var bcc = InstructionDescription{
	Name: "BCC",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return !cpu.GetStatus(CarryFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var bcs = InstructionDescription{
	Name: "BCS",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return cpu.GetStatus(CarryFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var beq = InstructionDescription{
	Name: "BEQ",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return cpu.GetStatus(ZeroFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var bmi = InstructionDescription{
	Name: "BMI",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return cpu.GetStatus(NegativeFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var bne = InstructionDescription{
	Name: "BNE",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return !cpu.GetStatus(ZeroFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var bpl = InstructionDescription{
	Name: "BPL",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return !cpu.GetStatus(NegativeFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var bvc = InstructionDescription{
	Name: "BVC",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return !cpu.GetStatus(OverflowFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}

var bvs = InstructionDescription{
	Name: "BVS",
	SubExec: TakeBranchFn(func(cpu LogicalCPU, bus LogicalBus) (bool, error) {
		return cpu.GetStatus(OverflowFlagBit), nil
	}),
	Access: RelativeAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x90: Relative,
	},
}
