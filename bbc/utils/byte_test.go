package utils

import (
	"testing"
)

func TestFromNibbles(t *testing.T) {
	base1 := uint16(0x2144)
	base2 := uint16(0x45F1)
	offset := uint8(0x20)

	if GetAddressPage(base1) != 0x2100 || GetAddressPage(base2) != 0x4500 {
		t.Fatal("wrong page")
	}

	if IsPageCrossed(base1, offset) || !IsPageCrossed(base2, offset) {
		t.Fatal("wrong crossing")
	}

	if GetAddressPage(SamePageOffset(base2, offset)) != GetAddressPage(base2) {
		t.Fatal("wrong page after same page offset")
	}
}
