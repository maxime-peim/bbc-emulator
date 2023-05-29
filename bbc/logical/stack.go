package logical

func pushRegister(register Register) ExecFn {
	return ExecFn(func(cpu LogicalCPU) error {
		value := cpu.GetRegister(register)
		cpu.Push(value)
		return nil
	})
}

func pullRegister(register Register) ExecFn {
	return ExecFn(func(cpu LogicalCPU) error {
		value, err := cpu.Pop()
		if err != nil {
			return err
		}
		cpu.SetStatus(value == 0, ZeroFlagBit)
		cpu.SetStatus(value&0x80 != 0, NegativeFlagBit)
		cpu.SetRegister(value, register)
		return nil
	})
}

var pha = InstructionDescription{
	Name:    "PHA",
	SubExec: pushRegister(RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x48: Implied,
	},
}

var php = InstructionDescription{
	Name:    "PHP",
	SubExec: pushRegister(RegisterStatus),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x08: Implied,
	},
}

var pla = InstructionDescription{
	Name:    "PLA",
	SubExec: pullRegister(RegisterA),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x68: Implied,
	},
}

var plp = InstructionDescription{
	Name:    "PLP",
	SubExec: pullRegister(RegisterStatus),
	Access:  ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x28: Implied,
	},
}
