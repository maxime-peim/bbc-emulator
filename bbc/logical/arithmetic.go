package logical

func addToAccumulator(invert bool) AfterReadFn {
	return AfterReadFn(func(value byte, cpu LogicalCPU) error {
		// invert value bit to perform substraction
		if invert {
			value = ^value
		}
		carry := uint8(0)
		if cpu.GetStatus(CarryFlagBit) {
			carry = 1
		}
		A := cpu.GetRegister(RegisterA)
		result7Bits := value&0x7F + A&0x7F + carry
		carry6 := result7Bits >> 7
		bits67 := value>>7 + A>>7 + carry6
		carry7 := bits67 >> 1
		result := result7Bits&0x7F + bits67<<7
		cpu.SetStatus(carry7 != 0, CarryFlagBit)
		cpu.SetStatus(carry7^carry6 != 0, OverflowFlagBit)
		cpu.SetRegister(result, RegisterA)
		return nil
	})
}

func cmpRegister(register Register) AfterReadFn {
	return AfterReadFn(func(value byte, cpu LogicalCPU) error {
		reg := cpu.GetRegister(register)
		cpu.SetStatus(reg >= value, CarryFlagBit)
		diff := reg - value
		cpu.SetStatus(diff&0x80 != 0, NegativeFlagBit)
		cpu.SetStatus(diff == 0, ZeroFlagBit)
		return nil
	})
}

var adc = InstructionDescription{
	Name:    "ADC",
	SubExec: addToAccumulator(false),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x69: Immediate,
		0x65: ZeroPage,
		0x75: ZeroPageX,
		0x6D: Absolute,
		0x7D: AbsoluteX,
		0x79: AbsoluteY,
		0x61: IndirectX,
		0x71: IndirectY,
	},
}

var sbc = InstructionDescription{
	Name:    "SBC",
	SubExec: addToAccumulator(true),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xE9: Immediate,
		0xE5: ZeroPage,
		0xF5: ZeroPageX,
		0xED: Absolute,
		0xFD: AbsoluteX,
		0xF9: AbsoluteY,
		0xE1: IndirectX,
		0xF1: IndirectY,
	},
}

var cmp = InstructionDescription{
	Name:    "CMP",
	SubExec: cmpRegister(RegisterA),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xC9: Immediate,
		0xC5: ZeroPage,
		0xD5: ZeroPageX,
		0xCD: Absolute,
		0xDD: AbsoluteX,
		0xD9: AbsoluteY,
		0xC1: IndirectX,
		0xD1: IndirectY,
	},
}

var cpx = InstructionDescription{
	Name:    "CPX",
	SubExec: cmpRegister(RegisterX),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xE0: Immediate,
		0xE4: ZeroPage,
		0xEC: Absolute,
	},
}

var cpy = InstructionDescription{
	Name:    "CPY",
	SubExec: cmpRegister(RegisterY),
	Access:  Read,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xC0: Immediate,
		0xC4: ZeroPage,
		0xCC: Absolute,
	},
}
