package main

import (
	"testing"
)

func TestPSW_SET(t *testing.T) {
	cases := []struct {
		Mask     byte
		Original byte
		Expected byte
	}{
		{Mask: 0b00000000, Original: 0b00000000, Expected: 0b00000000},
		{Mask: (1 << 0), Original: 0x00, Expected: (1 << 0)},
		{Mask: (1 << 1), Original: 0x00, Expected: (1 << 1)},
		{Mask: (1 << 2), Original: 0x00, Expected: (1 << 2)},
		{Mask: (1 << 3), Original: 0x00, Expected: (1 << 3)},
		{Mask: (1 << 4), Original: 0x00, Expected: (1 << 4)},
		{Mask: (1 << 5), Original: 0x00, Expected: (1 << 5)},
		{Mask: (1 << 6), Original: 0x00, Expected: (1 << 6)},
		{Mask: (1 << 7), Original: 0x00, Expected: (1 << 7)},
	}

	for _, tc := range cases {
		actual := PSW_SET(tc.Original, tc.Mask)

		if actual != tc.Expected {
			t.Errorf("expected %#08b got %#08b", tc.Expected, actual)
		}
	}
}

func TestPSW_UNSET(t *testing.T) {
	cases := []struct {
		Mask     byte
		Original byte
		Expected byte
	}{
		{Mask: 0b00000000, Original: 0b00000000, Expected: 0b00000000},
		{Mask: (1 << 0), Original: 0xFF, Expected: 0b11111110},
		{Mask: (1 << 1), Original: 0xFF, Expected: 0b11111101},
		{Mask: (1 << 2), Original: 0xFF, Expected: 0b11111011},
		{Mask: (1 << 3), Original: 0xFF, Expected: 0b11110111},
		{Mask: (1 << 4), Original: 0xFF, Expected: 0b11101111},
		{Mask: (1 << 5), Original: 0xFF, Expected: 0b11011111},
		{Mask: (1 << 6), Original: 0xFF, Expected: 0b10111111},
		{Mask: (1 << 7), Original: 0xFF, Expected: 0b01111111},
	}

	for _, tc := range cases {
		actual := PSW_UNSET(tc.Original, tc.Mask)

		if actual != tc.Expected {
			t.Errorf("expected %#08b got %#08b", tc.Expected, actual)
		}
	}
}

func TestOp0x03(t *testing.T) {
	cases := []struct {
		Original byte
		Expected byte
	}{
		{Original: 0b00000000, Expected: 0b00000000},
		{Original: 0b11111111, Expected: 0b11111111},
		{Original: 0b10000001, Expected: 0b11000000},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.Original); err != nil {
			t.Fatalf("unexpected error when writing to memory: %s", err)
		}

		err := vm.Feed([]byte{0x03})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if actualAcc, err := vm.ReadMem(SFR_ACC); err != nil {
			t.Fatalf("unexpected error when reading from memory (%#02x): %s", SFR_ACC, err)
		} else {
			if actualAcc != tc.Expected {
				t.Errorf("expected A register to be %#08b, got %#08b", tc.Expected, actualAcc)
			}
		}
	}
}

func TestOp0x04(t *testing.T) {
	cases := []struct {
		Original byte
		Expected byte
	}{
		{Original: 0b00000000, Expected: 0b00000001},
		{Original: 0xFF, Expected: 0x00},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.Original); err != nil {
			t.Fatalf("unexpected error when writing to memory: %s", err)
		}

		vm.Feed([]byte{0x04})

		if actualAcc, err := vm.ReadMem(SFR_ACC); err != nil {
			t.Fatalf("unexpected error when reading from memory (%#02x): %s", SFR_ACC, err)
		} else {
			if actualAcc != tc.Expected {
				t.Errorf("expected A register to be %#08b, got %#08b", tc.Expected, actualAcc)
			}
		}
	}
}

func TestOp0x05(t *testing.T) {
	cases := []struct {
		Addr         uint8
		InitialValue byte
	}{
		{Addr: 0, InitialValue: 10},
		{Addr: 1, InitialValue: 20},
	}

	for _, tc := range cases {
		vm := NewMachine()
		err := vm.WriteMem(tc.Addr, tc.InitialValue)
		if err != nil {
			t.Fatal(err)
		}

		err = vm.Feed([]byte{0x05, tc.Addr})
		if err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue + 1
		if actualValue != expectedValue {
			t.Errorf("expected value at address %#02x to be %#02x, got %#02x", tc.Addr, expectedValue, actualValue)
		}
	}
}

func TestOp0x06_0x07(t *testing.T) {
	cases := []struct {
		Opcode       byte
		Name         string
		Addr         uint8
		Ptr          uint8
		InitialValue byte
	}{
		{Opcode: 0x06, Name: "@R0", Addr: LOC_R0, Ptr: 1, InitialValue: 0x00}, // value of R0 is 1, value at address 1 is 0x00
		{Opcode: 0x07, Name: "@R1", Addr: LOC_R1, Ptr: 2, InitialValue: 0xA0},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		vm.Feed([]byte{tc.Opcode})

		directRead, err := vm.ReadMem(tc.Ptr)
		if err != nil {
			t.Fatal(err)
		}

		indirectRead, err := vm.DerefMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue + 1
		if directRead != expectedValue {
			t.Errorf("expected %s to be %#02x, got %#02x", tc.Name, expectedValue, directRead)
		}

		if indirectRead != expectedValue {
			t.Errorf("expected %s to be %#02x, got %#02x", tc.Name, expectedValue, indirectRead)
		}
	}
}

// INC R0-R7 (direct addressing)
func TestOp0x08_0x0f(t *testing.T) {
	cases := []struct {
		Name         string
		Opcode       byte
		Addr         uint8
		InitialValue byte
	}{
		{Name: "R0", Opcode: 0x08, Addr: LOC_R0, InitialValue: 0x00},
		{Name: "R0", Opcode: 0x08, Addr: LOC_R0, InitialValue: 0xFF},
		{Name: "R1", Opcode: 0x09, Addr: LOC_R1, InitialValue: 0x01},
		{Name: "R2", Opcode: 0x0a, Addr: LOC_R2, InitialValue: 0x02},
		{Name: "R3", Opcode: 0x0b, Addr: LOC_R3, InitialValue: 0x03},
		{Name: "R4", Opcode: 0x0c, Addr: LOC_R4, InitialValue: 0x04},
		{Name: "R5", Opcode: 0x0d, Addr: LOC_R5, InitialValue: 0x05},
		{Name: "R6", Opcode: 0x0e, Addr: LOC_R6, InitialValue: 0x06},
		{Name: "R7", Opcode: 0x0f, Addr: LOC_R7, InitialValue: 0x07},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue + 1
		if actualValue != expectedValue {
			t.Errorf("expected value in register %s to be %#02x, got %02x (opcode %#02x)",
				tc.Name, expectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x13(t *testing.T) {
	cases := []struct {
		Original      byte
		ExpectedAcc   byte
		ExpectedCarry byte
	}{
		{Original: 0b00000000, ExpectedAcc: 0b00000000, ExpectedCarry: 0},
		{Original: 0b11111111, ExpectedAcc: 0b01111111, ExpectedCarry: 1},
		{Original: 0b10000001, ExpectedAcc: 0b01000000, ExpectedCarry: 1},
		{Original: 0b11000101, ExpectedAcc: 0b01100010, ExpectedCarry: 1},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.Original); err != nil {
			t.Fatalf("unexpected error when writing to memory: %s", err)
		}

		vm.Feed([]byte{0x13})

		actualAcc, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			t.Fatalf("unexpected error when reading from memory (%#02x): %s", SFR_ACC, err)
		}

		psw, err := vm.ReadMem(SFR_PSW)
		if err != nil {
			t.Fatalf("unexpected error when reading PSW register from memory (%#02x): %s", SFR_PSW, err)
		}

		var actualCarry byte
		if psw&PSW_C_MASK > 0 {
			actualCarry = 1
		} else {
			actualCarry = 0
		}

		if actualAcc != tc.ExpectedAcc {
			t.Errorf("expected A register to be %#08b, got %#08b", tc.ExpectedAcc, actualAcc)
		}

		if actualCarry != tc.ExpectedCarry {
			t.Errorf("expected Carry register to be %#08b, got %#08b", tc.ExpectedCarry, actualCarry)
		}
	}
}

func TestOp0x14(t *testing.T) {
	cases := []struct {
		Original byte
		Expected byte
	}{
		{Original: 0x00, Expected: 0xFF},
		{Original: 0x01, Expected: 0x00},
		{Original: 0xFF, Expected: 0xFE},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.Original); err != nil {
			t.Fatalf("unexpected error when writing to memory: %s", err)
		}

		vm.Feed([]byte{0x14})

		if actualAcc, err := vm.ReadMem(SFR_ACC); err != nil {
			t.Fatalf("unexpected error when reading from memory (%#02x): %s", SFR_ACC, err)
		} else {
			if tc.Expected != actualAcc {
				t.Errorf("expected A register to be %#08b, got %08b", tc.Expected, actualAcc)
			}
		}
	}
}

// DEC ramaddr
func TestOp0x15(t *testing.T) {
	cases := []struct {
		Addr         uint8
		InitialValue byte
	}{
		{Addr: 0x00, InitialValue: 0xAB},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0x15, tc.Addr}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue - 1
		if actualValue != expectedValue {
			t.Errorf("expected value in address %#02x to be %#02x, got %#02x", tc.Addr, expectedValue, actualValue)
		}
	}
}

// DEC @R0, @R1
func TestOp0x16_0x17(t *testing.T) {
	cases := []struct {
		Opcode       byte
		Name         string
		Addr         uint8
		Ptr          uint8
		InitialValue byte
	}{
		{Opcode: 0x16, Name: "@R0", Addr: LOC_R0, Ptr: 1, InitialValue: 0xFF}, // value of R0 is 1, value at address 1 is 0xFF
		{Opcode: 0x17, Name: "@R1", Addr: LOC_R1, Ptr: 2, InitialValue: 0xAB},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue - 1
		actualValue, err := vm.DerefMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if expectedValue != actualValue {
			t.Errorf("expected value at %s to be %#02x, got %#02x", tc.Name, expectedValue, actualValue)
		}
	}
}

// DEC R0-R7
func TestOp0x18_0x1f(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Addr          uint8
		InitialValue  byte
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x18, Addr: LOC_R0, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R0", Opcode: 0x18, Addr: LOC_R0, InitialValue: 0xFF, ExpectedValue: 0xFF - 1},
		{Name: "R1", Opcode: 0x19, Addr: LOC_R1, InitialValue: 0x01, ExpectedValue: 0x01 - 1},
		{Name: "R2", Opcode: 0x1a, Addr: LOC_R2, InitialValue: 0x02, ExpectedValue: 0x02 - 1},
		{Name: "R3", Opcode: 0x1b, Addr: LOC_R3, InitialValue: 0x03, ExpectedValue: 0x03 - 1},
		{Name: "R4", Opcode: 0x1c, Addr: LOC_R4, InitialValue: 0x04, ExpectedValue: 0x04 - 1},
		{Name: "R5", Opcode: 0x1d, Addr: LOC_R5, InitialValue: 0x05, ExpectedValue: 0x05 - 1},
		{Name: "R6", Opcode: 0x1e, Addr: LOC_R6, InitialValue: 0x06, ExpectedValue: 0x06 - 1},
		{Name: "R7", Opcode: 0x1f, Addr: LOC_R7, InitialValue: 0x07, ExpectedValue: 0x07 - 1},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected value in register %s to be %#02x, got %02x (opcode %#02x)",
				tc.Name, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x23(t *testing.T) {
	cases := []struct {
		Original    byte
		ExpectedAcc byte
	}{
		{Original: 0x00, ExpectedAcc: 0x00},
		{Original: 0xFF, ExpectedAcc: 0xFF},
		{Original: 0b11000101, ExpectedAcc: 0b10001011},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.Original); err != nil {
			t.Fatalf("unexpected error when writing to memory: %s", err)
		}

		vm.Feed([]byte{0x23})

		if actualAcc, err := vm.ReadMem(SFR_ACC); err != nil {
			t.Fatalf("unexpected error when reading from memory (%#02x): %s", SFR_ACC, err)
		} else {
			if tc.ExpectedAcc != actualAcc {
				t.Errorf("expected A register to be %#08b, got %#08b", tc.ExpectedAcc, actualAcc)
			}
		}
	}
}

func TestOp0x24(t *testing.T) {
	cases := []struct {
		InitialValue  byte
		AddAmount     byte
		ExpectedValue byte
	}{
		{InitialValue: 0x00, AddAmount: 0xFF, ExpectedValue: 0xFF},
		{InitialValue: 0xFF, AddAmount: 0x1, ExpectedValue: 0},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0x24, tc.AddAmount}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected A register to be %#02x, got %#02x", tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x25(t *testing.T) {
	cases := []struct {
		Addr          uint8
		InitialValue  byte
		AddAmount     byte
		ExpectedValue byte
	}{
		{Addr: 0x00, InitialValue: 0x01, AddAmount: 0xAA, ExpectedValue: 0xAB},
		{Addr: 0x01, InitialValue: 0xFF, AddAmount: 0x1, ExpectedValue: 0x0},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.AddAmount, tc.AddAmount); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0x25, tc.AddAmount}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.AddAmount)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected A register to be %#02x, got %#02x", tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x26_0x27(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Addr          uint8
		Ptr           uint8
		InitialValue  byte
		AddAmount     byte
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x26, Addr: LOC_R0, Ptr: 0x02, InitialValue: 0xAA, AddAmount: 0x01, ExpectedValue: 0xAB},
		{Name: "@R1", Opcode: 0x27, Addr: LOC_R1, Ptr: 0xFF, InitialValue: 0xFF, AddAmount: 0x01, ExpectedValue: 0x00},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.AddAmount); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Fatalf("expected %s to be %#02x, got %#02x (opcode %#02x)", tc.Name, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x28_0x2f(t *testing.T) {
	cases := []struct {
		Name            string
		Opcode          byte
		Addr            byte
		AccInitialValue byte
		AddAmount       byte
		ExpectedValue   byte
	}{
		{Name: "R0", Opcode: 0x28, Addr: LOC_R0, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R1", Opcode: 0x29, Addr: LOC_R1, AccInitialValue: 0x01 + 1, AddAmount: 0x02, ExpectedValue: 0x03 + 1},
		{Name: "R2", Opcode: 0x2a, Addr: LOC_R2, AccInitialValue: 0x01 + 2, AddAmount: 0x02, ExpectedValue: 0x03 + 2},
		{Name: "R3", Opcode: 0x2b, Addr: LOC_R3, AccInitialValue: 0x01 + 3, AddAmount: 0x02, ExpectedValue: 0x03 + 3},
		{Name: "R4", Opcode: 0x2c, Addr: LOC_R4, AccInitialValue: 0x01 + 4, AddAmount: 0x02, ExpectedValue: 0x03 + 4},
		{Name: "R5", Opcode: 0x2d, Addr: LOC_R5, AccInitialValue: 0x01 + 5, AddAmount: 0x02, ExpectedValue: 0x03 + 5},
		{Name: "R6", Opcode: 0x2e, Addr: LOC_R6, AccInitialValue: 0x01 + 6, AddAmount: 0x02, ExpectedValue: 0x03 + 6},
		{Name: "R7", Opcode: 0x2f, Addr: LOC_R7, AccInitialValue: 0x01 + 7, AddAmount: 0x02, ExpectedValue: 0x03 + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.AccInitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Addr, tc.AddAmount); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected A+%s (%#02x+%#02x) to be %#02x, got %#02x (opcode %#02x)", tc.Name, tc.AccInitialValue, tc.AddAmount, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x33(t *testing.T) {
	cases := []struct {
		A             byte
		ExpectedA     byte
		ExpectedCarry byte
	}{
		{A: 0b11000101, ExpectedA: 0b10001010, ExpectedCarry: 1},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(SFR_ACC, tc.A); err != nil {
			t.Fatalf("unexpected error when writing to memory: %s", err)
		}

		vm.Feed([]byte{0x33})

		if actualAcc, err := vm.ReadMem(SFR_ACC); err != nil {
			t.Fatalf("unexpected error when reading from memory (%#02x): %s", SFR_ACC, err)
		} else {
			if actualAcc != tc.ExpectedA {
				t.Errorf("expected A register to be %#08b, got %#08b", tc.ExpectedA, actualAcc)
			}
		}
	}
}

func TestOp0x74(t *testing.T) {
	cases := []struct {
		ExpectedValue byte
	}{
		{ExpectedValue: 0xAA},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.Feed([]byte{0x74, tc.ExpectedValue}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected register A to be %#02x, got %#02x", tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x75(t *testing.T) {
	cases := []struct {
		Addr          uint8
		ExpectedValue byte
	}{
		{Addr: 0x0, ExpectedValue: 0xAD},
		{Addr: 0xFF, ExpectedValue: 0xBB},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.Feed([]byte{0x75, tc.Addr, tc.ExpectedValue}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected value at address %#02x to be %#02x, got %#02x", tc.Addr, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x76_0x77(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Addr          uint8
		Ptr           uint8
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x76, Addr: LOC_R0, Ptr: 0x02, ExpectedValue: 0xAA},
		{Name: "@R1", Opcode: 0x77, Addr: LOC_R0, Ptr: 0x02, ExpectedValue: 0xDD},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0x76, tc.ExpectedValue}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected %s to be %#02x, got %#02x", tc.Name, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x78_0x7f(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Addr          uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x78, Addr: LOC_R0, ExpectedValue: 0xAA + 0},
		{Name: "R1", Opcode: 0x79, Addr: LOC_R1, ExpectedValue: 0xAA + 1},
		{Name: "R2", Opcode: 0x7a, Addr: LOC_R2, ExpectedValue: 0xAA + 2},
		{Name: "R3", Opcode: 0x7b, Addr: LOC_R3, ExpectedValue: 0xAA + 3},
		{Name: "R4", Opcode: 0x7c, Addr: LOC_R4, ExpectedValue: 0xAA + 4},
		{Name: "R5", Opcode: 0x7d, Addr: LOC_R5, ExpectedValue: 0xAA + 5},
		{Name: "R6", Opcode: 0x7e, Addr: LOC_R6, ExpectedValue: 0xAA + 6},
		{Name: "R7", Opcode: 0x7f, Addr: LOC_R7, ExpectedValue: 0xAA + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.Feed([]byte{tc.Opcode, tc.ExpectedValue}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected value in %s to be %#02x, got %#02x", tc.Name, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x86_0x87(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		MoveFrom      uint8
		Ptr           uint8
		MoveInto      uint8
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x86, MoveFrom: LOC_R0, Ptr: 0x01, MoveInto: 0x02, ExpectedValue: 0xAA},
		{Name: "@R1", Opcode: 0x87, MoveFrom: LOC_R1, Ptr: 0x03, MoveInto: 0x04, ExpectedValue: 0xBB},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.WriteMem(tc.MoveFrom, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveInto}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.MoveInto)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected memory address %#02x to be %#02x, got %#02x", tc.MoveInto, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x88_0x8f(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x88, MoveFrom: LOC_R0, MoveInto: 0x00 + 0, ExpectedValue: 0xAA + 1},
		{Name: "R1", Opcode: 0x89, MoveFrom: LOC_R1, MoveInto: 0x00 + 1, ExpectedValue: 0xAA + 1},
		{Name: "R2", Opcode: 0x8a, MoveFrom: LOC_R2, MoveInto: 0x00 + 2, ExpectedValue: 0xAA + 1},
		{Name: "R3", Opcode: 0x8b, MoveFrom: LOC_R3, MoveInto: 0x00 + 3, ExpectedValue: 0xAA + 1},
		{Name: "R4", Opcode: 0x8c, MoveFrom: LOC_R4, MoveInto: 0x00 + 4, ExpectedValue: 0xAA + 1},
		{Name: "R5", Opcode: 0x8d, MoveFrom: LOC_R5, MoveInto: 0x00 + 5, ExpectedValue: 0xAA + 1},
		{Name: "R6", Opcode: 0x8e, MoveFrom: LOC_R6, MoveInto: 0x00 + 6, ExpectedValue: 0xAA + 1},
		{Name: "R7", Opcode: 0x8f, MoveFrom: LOC_R7, MoveInto: 0x00 + 7, ExpectedValue: 0xAA + 1},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveInto}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.MoveInto)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected value at address %#02x to be %#02x, got %#02x (opcode %#02x)", tc.Name, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xA6_0xA7(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Dest          uint8
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		// R0 points to 0x01, addr 0x02 contains the data that we want (0xAA). deref'ing R0 should return 0xAA
		{Name: "@R0", Opcode: 0xA6, Dest: LOC_R0, MoveInto: 0x20, MoveFrom: 0x10, ExpectedValue: 0xAA},
		{Name: "@R1", Opcode: 0xA7, Dest: LOC_R1, MoveInto: 0x20, MoveFrom: 0x10, ExpectedValue: 0xAB},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Dest, tc.MoveInto); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveFrom}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefMem(tc.Dest)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected value at address %#02x to be %#02x, got %#02x (opcode %#02x)",
				tc.Name, tc.Dest, tc.ExpectedValue, actualValue, tc.Opcode,
			)
		}
	}
}

func TestOp0xA8_0xAF(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0xA8, MoveInto: LOC_R0, MoveFrom: 0xAA + 0, ExpectedValue: 0xBB + 0},
		{Name: "R1", Opcode: 0xA9, MoveInto: LOC_R1, MoveFrom: 0xAA + 1, ExpectedValue: 0xBB + 1},
		{Name: "R2", Opcode: 0xAA, MoveInto: LOC_R2, MoveFrom: 0xAA + 2, ExpectedValue: 0xBB + 2},
		{Name: "R3", Opcode: 0xAB, MoveInto: LOC_R3, MoveFrom: 0xAA + 3, ExpectedValue: 0xBB + 3},
		{Name: "R4", Opcode: 0xAC, MoveInto: LOC_R4, MoveFrom: 0xAA + 4, ExpectedValue: 0xBB + 4},
		{Name: "R5", Opcode: 0xAD, MoveInto: LOC_R5, MoveFrom: 0xAA + 5, ExpectedValue: 0xBB + 5},
		{Name: "R6", Opcode: 0xAE, MoveInto: LOC_R6, MoveFrom: 0xAA + 6, ExpectedValue: 0xBB + 6},
		{Name: "R7", Opcode: 0xAF, MoveInto: LOC_R7, MoveFrom: 0xAA + 7, ExpectedValue: 0xBB + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveFrom}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.MoveInto)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected address %#02x to be %#02x, got %#02x (opcode %#02x)", tc.Name, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xE6_0xE7(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		MoveFrom      uint8
		Ptr           uint8
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0xE6, MoveFrom: LOC_R0, Ptr: 0xAA, ExpectedValue: 0xDD},
		{Name: "@R1", Opcode: 0xE7, MoveFrom: LOC_R1, Ptr: 0xAB, ExpectedValue: 0xDD},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.MoveFrom, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected A register to be %#02x, got %#02x (opcode %#02x)", tc.Name, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xE5_0xE8_0xEF(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		{Name: "ramaddr", Opcode: 0xE5, MoveInto: SFR_ACC, MoveFrom: 0xAD, ExpectedValue: 0xBB + 0},
		{Name: "R0", Opcode: 0xE8, MoveInto: SFR_ACC, MoveFrom: LOC_R0, ExpectedValue: 0xBB + 0},
		{Name: "R1", Opcode: 0xE9, MoveInto: SFR_ACC, MoveFrom: LOC_R1, ExpectedValue: 0xBB + 1},
		{Name: "R2", Opcode: 0xEA, MoveInto: SFR_ACC, MoveFrom: LOC_R2, ExpectedValue: 0xBB + 2},
		{Name: "R3", Opcode: 0xEB, MoveInto: SFR_ACC, MoveFrom: LOC_R3, ExpectedValue: 0xBB + 3},
		{Name: "R4", Opcode: 0xEC, MoveInto: SFR_ACC, MoveFrom: LOC_R4, ExpectedValue: 0xBB + 4},
		{Name: "R5", Opcode: 0xED, MoveInto: SFR_ACC, MoveFrom: LOC_R5, ExpectedValue: 0xBB + 5},
		{Name: "R6", Opcode: 0xEE, MoveInto: SFR_ACC, MoveFrom: LOC_R6, ExpectedValue: 0xBB + 6},
		{Name: "R7", Opcode: 0xEF, MoveInto: SFR_ACC, MoveFrom: LOC_R7, ExpectedValue: 0xBB + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveFrom}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(tc.MoveInto)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected address %#02x to be %#02x, got %#02x (opcode %#02x)", tc.Name, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xF5(t *testing.T) {
	cases := []struct {
		CopyInto      uint8
		ExpectedValue byte
	}{
		{CopyInto: 0xFF, ExpectedValue: 0xAA},
	}

	for _, tc := range cases {
		vm := NewMachine()

		var expectedValue byte = tc.ExpectedValue
		var copyInto uint8 = tc.CopyInto

		if err := vm.WriteMem(SFR_ACC, expectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0xf5, copyInto}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadMem(copyInto)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != expectedValue {
			t.Errorf("expected address %#02x to be %#02x, got %#02x", copyInto, expectedValue, actualValue)
		}
	}
}

func TestOp0xF6(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xF7(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xF8(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xF9(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xFA(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xFB(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xFC(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xFD(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xFE(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xFF(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}
