package utils

func AddressFromNibbles(high, low uint8) uint16 {
	return (uint16(high) << 8) | uint16(low)
}

func AddressToNibbles(addr uint16) (uint8, uint8) {
	return uint8(addr >> 8), uint8(addr & 0xFF)
}

func IsPageCrossed(addr uint16, offset uint8) bool {
	return (addr+uint16(offset))&0xFF00 != addr&0xFF00
}

func SamePageOffset(addr uint16, offset uint8) uint16 {
	page := addr & 0xFF00
	low := uint8(addr&0xff) + offset
	return page | uint16(low)
}
