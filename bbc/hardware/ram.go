package hardware

import (
	"bbc/utils"
)

type RAM struct {
	memory []byte
	bus    *Bus
}

func (ram *RAM) GetName() string    { return "RAM" }
func (ram *RAM) PlugToBus(bus *Bus) { ram.bus = bus }
func (ram *RAM) IsWritable() bool   { return true }
func (ram *RAM) IsReadable() bool   { return true }
func (ram *RAM) GetSegment() *utils.Segment {
	return AdressableSegment
}

func (ram *RAM) Start() error {
	return nil
}

func (ram *RAM) Reset() error {
	return nil
}

func (ram *RAM) Stop() error {
	return nil
}

func (ram *RAM) DirectRead(addr uint16) (byte, error) {
	return ram.memory[addr], nil
}

func (ram *RAM) OffsetRead(base uint16, offset uint8) (byte, uint16, error) {
	if uint8(base&0xff)+offset < offset {
		// page crossed
		if err := ram.bus.Clock.Tick(); err != nil {
			return 0, 0, err
		}
	}
	addr := base + uint16(offset)
	value, err := ram.DirectRead(addr)
	if err != nil {
		return 0, 0, err
	}
	return value, addr, nil
}

func (ram *RAM) DirectWrite(value byte, addr uint16) error {
	ram.memory[addr] = value
	return nil
}

func (ram *RAM) OffsetWrite(value byte, base uint16, offset uint8) error {
	if uint8(base&0xff)+offset < offset {
		// page crossed
		if err := ram.bus.Clock.Tick(); err != nil {
			return err
		}
	}
	return ram.DirectWrite(value, base+uint16(offset))
}

func (ram *RAM) Clear() {
	for i := uint16(0); i <= AdressableSegment.End; i++ {
		ram.memory[i] = 0
	}
}

func NewRAM() *RAM {
	return &RAM{
		memory: make([]byte, AdressableSegment.Size()),
	}
}