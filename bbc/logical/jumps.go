package logical

import "bbc/utils"

var jmp = InstructionDescription{
	Name: "JMP",
	SubExec: SetupJumpFn(func(addr uint16, cpu LogicalCPU, bus LogicalBus) error {
		nextPCH, nextPCL := utils.AddressToNibbles(addr)
		cpu.SetRegister(nextPCL, RegisterPCL)
		cpu.SetRegister(nextPCH, RegisterPCH)
		return nil
	}),
	Access: JumpAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x4C: Absolute,
		0x6C: Indirect,
	},
}

var jsr = InstructionDescription{
	Name: "JSR",
	SubExec: SetupJumpFn(func(addr uint16, cpu LogicalCPU, bus LogicalBus) error {
		pch := cpu.GetRegister(RegisterPCH)
		pcl := cpu.GetRegister(RegisterPCL)
		pc := utils.AddressFromNibbles(pch, pcl)
		nextPCH, nextPCL := utils.AddressToNibbles(pc + 2)
		// https://www.nesdev.org/6502_cpu.txt
		// internal operation (predecrement S?)
		if err := bus.Tick(); err != nil {
			return err
		}
		cpu.Push(nextPCL)
		cpu.Push(nextPCH)
		nextPCH, nextPCL = utils.AddressToNibbles(addr)
		cpu.SetRegister(nextPCL, RegisterPCL)
		cpu.SetRegister(nextPCH, RegisterPCH)
		return nil
	}),
	Access: JumpAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x20: Absolute,
	},
}

var rts = InstructionDescription{
	Name: "RTS",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		pcl, err := cpu.Pop()
		if err != nil {
			return err
		}
		pch, err := cpu.Pop()
		if err != nil {
			return err
		}
		nextPC := utils.AddressFromNibbles(pch, pcl) + 1
		nextPCH, nextPCL := utils.AddressToNibbles(nextPC)
		cpu.SetRegister(nextPCL, RegisterPCL)
		cpu.SetRegister(nextPCH, RegisterPCH)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x60: Implied,
	},
}
