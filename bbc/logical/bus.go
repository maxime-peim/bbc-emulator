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
	// 1 cycle
	DirectRead(uint16) (byte, error)
	// 1 cycle, +1 if page crossed or forced
	OffsetRead(uint16, uint8, bool) (byte, uint16, error)

	// 1 cycle
	DirectWrite(byte, uint16) error
	// 1 cycle, +1 if page crossed or force
	OffsetWrite(byte, uint16, uint8, bool) (uint16, error)

	Reset()
	Tick() error
}
