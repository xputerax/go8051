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
		vm := Machine{Registers: Register{ACC: tc.Original}}
		err := vm.Feed([]byte{0x03})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if vm.Registers.ACC != tc.Expected {
			t.Errorf("expected A register to be %#08b, got %#08b", tc.Expected, vm.Registers.ACC)
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
		vm := Machine{Registers: Register{ACC: tc.Original}}

		vm.Feed([]byte{0x04})

		if vm.Registers.ACC != tc.Expected {
			t.Errorf("expected A register to be %#08b, got %#08b", tc.Expected, vm.Registers.ACC)
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
		{Original: 0b11111111, ExpectedAcc: 0b11111111, ExpectedCarry: 1},
		{Original: 0b10000001, ExpectedAcc: 0b11000000, ExpectedCarry: 1},
	}

	for _, tc := range cases {
		vm := Machine{Registers: Register{ACC: tc.Original}}
		vm.Feed([]byte{0x13})

		actualAcc := vm.Registers.ACC
		var actualCarry byte
		if vm.Registers.PSW&PSW_C_MASK > 0 {
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
		vm := Machine{Registers: Register{ACC: tc.Original}}
		vm.Feed([]byte{0x14})

		actualAcc := vm.Registers.ACC

		if tc.Expected != actualAcc {
			t.Errorf("expected A register to be %#08b, got %08b", tc.Expected, actualAcc)
		}
	}
}
