package logical

import "bbc/utils"

var brk = InstructionDescription{
	Name: "BRK",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		pcl := cpu.GetRegister(RegisterPCL)
		if err := cpu.Push(pcl); err != nil {
			return err
		}
		pch := cpu.GetRegister(RegisterPCH)
		if err := cpu.Push(pch); err != nil {
			return err
		}
		status := cpu.GetRegister(RegisterStatus)
		if err := cpu.Push(status); err != nil {
			return err
		}
		cpu.SetStatus(true, BreakFlagBit)
		pcl, err := bus.DirectRead(IRQVectorAddr0)
		if err != nil {
			return err
		}
		pch, err = bus.DirectRead(IRQVectorAddr1)
		if err != nil {
			return err
		}
		nextPC := utils.AddressFromNibbles(pch, pcl)
		nextPCH, nextPCL := utils.AddressToNibbles(nextPC)
		cpu.SetRegister(nextPCL, RegisterPCL)
		cpu.SetRegister(nextPCH, RegisterPCH)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x00: Implied,
	},
}

var nop = InstructionDescription{
	Name: "NOP",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0xEA: Implied,
	},
}

var rti = InstructionDescription{
	Name: "RTI",
	SubExec: ExecFn(func(cpu LogicalCPU, bus LogicalBus) error {
		status, err := cpu.Pop()
		if err != nil {
			return err
		}
		cpu.SetRegister(status, RegisterStatus)
		pcl, err := cpu.Pop()
		if err != nil {
			return err
		}
		pch, err := cpu.Pop()
		if err != nil {
			return err
		}
		cpu.SetRegister(pcl, RegisterPCL)
		cpu.SetRegister(pch, RegisterPCH)
		return nil
	}),
	Access: ImpliedAccess,
	OpcodeMapping: map[Opcode]AddressingMode{
		0x40: Implied,
	},
}
