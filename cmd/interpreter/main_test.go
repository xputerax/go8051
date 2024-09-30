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
