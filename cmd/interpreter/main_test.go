package main

import (
	"testing"
)

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
