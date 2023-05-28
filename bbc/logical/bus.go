package logical

import "bbc/utils"

const (
	IRQVectorAddr0 uint16 = 0xFFFE
	IRQVectorAddr1 uint16 = 0xFFFF
)

var (
	AdressableSegment = utils.NewSegment(0x0000, 0xFFFF)
	StackSegment      = utils.NewSegment(0x0100, 0x1FF)
	ZeroPageSegment   = utils.NewSegment(0x0000, 0x00FF)
)

type LogicalBus interface {
	DirectRead(uint16) (byte, error)
	OffsetRead(uint16, uint8) (byte, uint16, error)

	DirectWrite(byte, uint16) error
	OffsetWrite(byte, uint16, uint8) error

	Reset()
	Tick() error
}
