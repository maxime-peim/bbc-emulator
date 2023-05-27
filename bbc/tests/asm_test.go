package tests

import (
	"testing"
)

func TestLDA(t *testing.T) {
	testCtx.Reset()

	initialPC := testCtx.cpu.ProgramCounter

	program := []byte{
		0xA9, 0x55, // LDA 0x55
	}
	testCtx.bus.WriteMultiple(program, initialPC)

	if err := testCtx.cpu.ExecuteNext(); err != nil {
		t.Fatalf(err.Error())
	}

	if testCtx.cpu.A != 0x55 {
		t.Fail()
	}
}

func TestSTA(t *testing.T) {
	testCtx.Reset()

	initialPC := testCtx.cpu.ProgramCounter

	program := []byte{
		0xA9, 0x55, // LDA 0x55
		0x85, 0x80, // STA 0x80 zero page
	}
	testCtx.bus.WriteMultiple(program, initialPC)

	for i := 0; i < 2; i++ {
		if err := testCtx.cpu.ExecuteNext(); err != nil {
			t.Fatalf(err.Error())
		}
	}

	value, err := testCtx.bus.DirectRead(0x0080)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value != 0x55 {
		t.Fail()
	}
}

func TestTAX(t *testing.T) {
	testCtx.Reset()

	initialPC := testCtx.cpu.ProgramCounter

	program := []byte{
		0xA9, 0x55, // LDA 0x55
		0xAA, // TAX
	}
	testCtx.bus.WriteMultiple(program, initialPC)

	for i := 0; i < 2; i++ {
		if err := testCtx.cpu.ExecuteNext(); err != nil {
			t.Fatalf(err.Error())
		}
	}

	if testCtx.cpu.X != 0x55 {
		t.Fail()
	}
}
