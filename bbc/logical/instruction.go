package logical

import (
	"fmt"
)

type Opcode byte
type AfterReadFn func(byte, LogicalCPU) error
type BeforeWriteFn func(LogicalCPU) (byte, error)
type OperationRMWFn func(byte, LogicalCPU) (byte, error)
type TakeBranchFn func(LogicalCPU) (bool, error)
type SetupJumpFn func(uint16, LogicalCPU) error
type ExecFn func(LogicalCPU) error

type Instruction struct {
	Name   string
	Access AccessMode

	subInstructionsByMode   map[AddressingMode]ExecFn
	subInstructionsByOpcode map[Opcode]ExecFn
}

func (instruction *Instruction) Execute(opcode Opcode, cpu LogicalCPU) error {
	exec, ok := instruction.subInstructionsByOpcode[opcode]
	if !ok {
		return fmt.Errorf("Opcode %x does not belong to instruction %s", opcode, instruction.Name)
	}
	return exec(cpu)
}

func (instruction *Instruction) GetOpcodes() []Opcode {
	opcodes := make([]Opcode, len(instruction.subInstructionsByOpcode))
	i := 0
	for opcode := range instruction.subInstructionsByOpcode {
		opcodes[i] = opcode
		i++
	}
	return opcodes
}

type InstructionDescription struct {
	Name          string
	SubExec       interface{}
	Access        AccessMode
	OpcodeMapping map[Opcode]AddressingMode
}

func (ins *InstructionDescription) RegisterTo(cpu LogicalCPU) error {
	if i := cpu.GetInstruction(ins.Name); i != nil {
		return fmt.Errorf("instruction %s already exists", ins.Name)
	}

	instruction := Instruction{
		Name:                    ins.Name,
		subInstructionsByMode:   map[AddressingMode]ExecFn{},
		subInstructionsByOpcode: map[Opcode]ExecFn{},
	}

	addressingFnForAccess, ok := AddressModeFetch[ins.Access]
	if !ok {
		return fmt.Errorf("addressing function does not exist for access %d", ins.Access)
	}

	fmt.Printf("Registration %s\n", ins.Name)

	for opcode, mode := range ins.OpcodeMapping {
		if i := cpu.GetInstructionByOpcode(opcode); i != nil {
			return fmt.Errorf("opcode %x already exists", opcode)
		}

		var execute ExecFn
		switch ins.Access {
		case ImpliedAccess:
			impliedAddressingFn := addressingFnForAccess[mode].(ExecFn)
			impliedInstructionFn, ok := ins.SubExec.(ExecFn)
			if !ok {
				return fmt.Errorf("access mode and sub-execute funtion signature don't match")
			}
			execute = func(cpu LogicalCPU) error {
				if err := impliedAddressingFn(cpu); err != nil {
					return err
				}
				return impliedInstructionFn(cpu)
			}
		case RelativeAccess:
			relativeAddressingFn := addressingFnForAccess[mode].(BranchFn)
			relativeInstructionFn, ok := ins.SubExec.(TakeBranchFn)
			if !ok {
				return fmt.Errorf("access mode and sub-execute funtion signature don't match")
			}
			execute = func(cpu LogicalCPU) error {
				return relativeAddressingFn(relativeInstructionFn, cpu)
			}
		case JumpAccess:
			jumpAddressingFn := addressingFnForAccess[mode].(JumpFn)
			jumpInstructionFn, ok := ins.SubExec.(SetupJumpFn)
			if !ok {
				return fmt.Errorf("access mode and sub-execute funtion signature don't match")
			}
			execute = func(cpu LogicalCPU) error {
				addr, err := jumpAddressingFn(cpu)
				if err != nil {
					return nil
				}
				return jumpInstructionFn(addr, cpu)
			}
		case Read:
			readAddressingFn := addressingFnForAccess[mode].(ReadFn)
			readInstructionFn, ok := ins.SubExec.(AfterReadFn)
			if !ok {
				return fmt.Errorf("access mode and sub-execute funtion signature don't match")
			}
			execute = func(cpu LogicalCPU) error {
				value, err := readAddressingFn(cpu)
				if err != nil {
					return err
				}
				return readInstructionFn(value, cpu)
			}
		case Write:
			writeAddressingFn := addressingFnForAccess[mode].(WriteFn)
			writeInstructionFn, ok := ins.SubExec.(BeforeWriteFn)
			if !ok {
				return fmt.Errorf("access mode and sub-execute funtion signature don't match")
			}
			execute = func(cpu LogicalCPU) error {
				value, err := writeInstructionFn(cpu)
				if err != nil {
					return err
				}
				return writeAddressingFn(value, cpu)
			}
		case ReadModifyWrite:
			rmwAddressingFn := addressingFnForAccess[mode].(ReadModifyWriteFn)
			rmwInstructionFn, ok := ins.SubExec.(OperationRMWFn)
			if !ok {
				return fmt.Errorf("access mode and sub-execute funtion signature don't match")
			}
			execute = func(cpu LogicalCPU) error {
				return rmwAddressingFn(rmwInstructionFn, cpu)
			}
		}

		instruction.subInstructionsByMode[mode] = ExecFn(execute)
		instruction.subInstructionsByOpcode[opcode] = ExecFn(execute)
	}
	cpu.SetInstruction(&instruction)
	return nil
}
