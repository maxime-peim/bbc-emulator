package logical

type StatusFlag uint32

const (
	CarryFlagBit StatusFlag = iota
	ZeroFlagBit
	InterruptDisableFlagBit
	DecimalModeFlagBit
	BreakFlagBit
	UnusedFlagBit
	OverflowFlagBit
	NegativeFlagBit
)

type Register uint8

const (
	RegisterA Register = iota
	RegisterX
	RegisterY
	RegisterStack
	RegisterStatus
	RegisterPCL
	RegisterPCH
)

type LogicalCPU interface {
	GetInstruction(string) *Instruction
	GetInstructionByOpcode(Opcode) *Instruction
	SetInstruction(*Instruction) error

	SetRegister(byte, Register)
	GetRegister(Register) byte

	SetStatus(bool, StatusFlag)
	GetStatus(StatusFlag) bool

	Push(byte) error
	Pop() (byte, error)

	NextByte() (byte, error)
	NextWord() (uint16, error)
}

var BaseInstructionSet = []InstructionDescription{
	lda, ldx, ldy, // load
	sta, stx, sty, // store
	tax, txa, tay, tya, txs, tsx, // transfert
	pha, php, pla, plp, // stack
	and, eor, ora, bit, // logical
	inc, inx, iny, dec, dex, dey, // increment / decrement
	adc, sbc, cmp, cpx, cpy, // arithmetic
	asl, lsr, rol, ror, // shifts
	clc, cld, cli, clv, sec, sed, sei, // status
	bcc, bcs, beq, bmi, bne, bpl, bvc, bvs, // branches
	jmp, jsr, rts, // jumps
}
