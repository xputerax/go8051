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
		BankNo       byte
		Addr         uint8
		Ptr          uint8
		InitialValue byte
	}{
		{Opcode: 0x06, Name: "@R0", BankNo: 0, Addr: LOC_R0, Ptr: 1, InitialValue: 0xA9}, // value of R0 is 1, value at address 1 is 0x00
		{Opcode: 0x06, Name: "@R0", BankNo: 1, Addr: LOC_R0, Ptr: 1, InitialValue: 0xA9},
		{Opcode: 0x06, Name: "@R0", BankNo: 2, Addr: LOC_R0, Ptr: 1, InitialValue: 0xA9},
		{Opcode: 0x06, Name: "@R0", BankNo: 3, Addr: LOC_R0, Ptr: 1, InitialValue: 0xA9},
		{Opcode: 0x07, Name: "@R1", BankNo: 0, Addr: LOC_R1, Ptr: 2, InitialValue: 0xA0},
		{Opcode: 0x07, Name: "@R1", BankNo: 1, Addr: LOC_R1, Ptr: 2, InitialValue: 0xA0},
		{Opcode: 0x07, Name: "@R1", BankNo: 2, Addr: LOC_R1, Ptr: 2, InitialValue: 0xA0},
		{Opcode: 0x07, Name: "@R1", BankNo: 3, Addr: LOC_R1, Ptr: 2, InitialValue: 0xA0},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.SetBankNo(tc.BankNo); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.Ptr); err != nil {
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

		indirectRead, err := vm.DerefMem(tc.Addr + vm.bankOffset())
		if err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue + 1
		if directRead != expectedValue {
			t.Errorf("expected direct read %s to be %#02x, got %#02x", tc.Name, expectedValue, directRead)
		}

		if indirectRead != expectedValue {
			t.Errorf("expected indirect read %s to be %#02x, got %#02x", tc.Name, expectedValue, indirectRead)
		}
	}
}

// INC R0-R7 (direct addressing)
func TestOp0x08_0x0f(t *testing.T) {
	cases := []struct {
		Name         string
		Opcode       byte
		Bank         byte
		Addr         uint8
		InitialValue byte
	}{
		{Name: "R0", Opcode: 0x08, Bank: 0, Addr: LOC_R0, InitialValue: 0xFF},
		{Name: "R0", Opcode: 0x08, Bank: 1, Addr: LOC_R0, InitialValue: 0xFF},
		{Name: "R1", Opcode: 0x09, Bank: 0, Addr: LOC_R1, InitialValue: 0xFF},
		{Name: "R1", Opcode: 0x09, Bank: 1, Addr: LOC_R1, InitialValue: 0xFF},
		{Name: "R2", Opcode: 0x0a, Bank: 0, Addr: LOC_R2, InitialValue: 0xFF},
		{Name: "R2", Opcode: 0x0a, Bank: 1, Addr: LOC_R2, InitialValue: 0xFF},
		{Name: "R3", Opcode: 0x0b, Bank: 0, Addr: LOC_R3, InitialValue: 0xFF},
		{Name: "R3", Opcode: 0x0b, Bank: 1, Addr: LOC_R3, InitialValue: 0xFF},
		{Name: "R4", Opcode: 0x0c, Bank: 0, Addr: LOC_R4, InitialValue: 0xFF},
		{Name: "R4", Opcode: 0x0c, Bank: 1, Addr: LOC_R4, InitialValue: 0xFF},
		{Name: "R5", Opcode: 0x0d, Bank: 0, Addr: LOC_R5, InitialValue: 0xFF},
		{Name: "R5", Opcode: 0x0d, Bank: 1, Addr: LOC_R5, InitialValue: 0xFF},
		{Name: "R6", Opcode: 0x0e, Bank: 0, Addr: LOC_R6, InitialValue: 0xFF},
		{Name: "R6", Opcode: 0x0e, Bank: 1, Addr: LOC_R6, InitialValue: 0xFF},
		{Name: "R7", Opcode: 0x0f, Bank: 0, Addr: LOC_R7, InitialValue: 0xFF},
		{Name: "R7", Opcode: 0x0f, Bank: 1, Addr: LOC_R7, InitialValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadBankMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue + 1
		if actualValue != expectedValue {
			t.Errorf("expected value in bank %d register %s to be %#02x, got %02x (opcode %#02x)",
				tc.Bank, tc.Name, expectedValue, actualValue, tc.Opcode)
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

// DEC @R0/@R1
func TestOp0x16_0x17(t *testing.T) {
	cases := []struct {
		Opcode       byte
		Name         string
		Addr         uint8
		Ptr          uint8
		InitialValue byte
		Bank         byte
	}{
		{Opcode: 0x16, Name: "@R0", Addr: LOC_R0, Ptr: 1, InitialValue: 0xFF, Bank: 1}, // value of R0 is 1, value at address 1 is 0xFF
		{Opcode: 0x17, Name: "@R1", Addr: LOC_R1, Ptr: 2, InitialValue: 0xAB, Bank: 1},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		expectedValue := tc.InitialValue - 1
		actualValue, err := vm.DerefBank(tc.Addr)
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
		Bank          byte
		Addr          uint8
		InitialValue  byte
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x18, Bank: 0, Addr: LOC_R0, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R0", Opcode: 0x18, Bank: 1, Addr: LOC_R0, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R0", Opcode: 0x18, Bank: 2, Addr: LOC_R0, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R0", Opcode: 0x18, Bank: 3, Addr: LOC_R0, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R1", Opcode: 0x19, Bank: 0, Addr: LOC_R1, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R1", Opcode: 0x19, Bank: 1, Addr: LOC_R1, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R1", Opcode: 0x19, Bank: 2, Addr: LOC_R1, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R1", Opcode: 0x19, Bank: 3, Addr: LOC_R1, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R2", Opcode: 0x1A, Bank: 0, Addr: LOC_R2, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R2", Opcode: 0x1A, Bank: 1, Addr: LOC_R2, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R2", Opcode: 0x1A, Bank: 2, Addr: LOC_R2, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R2", Opcode: 0x1A, Bank: 3, Addr: LOC_R2, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R3", Opcode: 0x1B, Bank: 0, Addr: LOC_R3, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R3", Opcode: 0x1B, Bank: 1, Addr: LOC_R3, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R3", Opcode: 0x1B, Bank: 2, Addr: LOC_R3, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R3", Opcode: 0x1B, Bank: 3, Addr: LOC_R3, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R4", Opcode: 0x1C, Bank: 0, Addr: LOC_R4, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R4", Opcode: 0x1C, Bank: 1, Addr: LOC_R4, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R4", Opcode: 0x1C, Bank: 2, Addr: LOC_R4, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R4", Opcode: 0x1C, Bank: 3, Addr: LOC_R4, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R5", Opcode: 0x1D, Bank: 0, Addr: LOC_R5, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R5", Opcode: 0x1D, Bank: 1, Addr: LOC_R5, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R5", Opcode: 0x1D, Bank: 2, Addr: LOC_R5, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R5", Opcode: 0x1D, Bank: 3, Addr: LOC_R5, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R6", Opcode: 0x1E, Bank: 0, Addr: LOC_R6, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R6", Opcode: 0x1E, Bank: 1, Addr: LOC_R6, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R6", Opcode: 0x1E, Bank: 2, Addr: LOC_R6, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R6", Opcode: 0x1E, Bank: 3, Addr: LOC_R6, InitialValue: 0x00, ExpectedValue: 0xFF},

		{Name: "R7", Opcode: 0x1F, Bank: 0, Addr: LOC_R7, InitialValue: 0x00 + 1, ExpectedValue: 0x00},
		{Name: "R7", Opcode: 0x1F, Bank: 1, Addr: LOC_R7, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R7", Opcode: 0x1F, Bank: 2, Addr: LOC_R7, InitialValue: 0x00, ExpectedValue: 0xFF},
		{Name: "R7", Opcode: 0x1F, Bank: 3, Addr: LOC_R7, InitialValue: 0x00, ExpectedValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadBankMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected value in bank %d register %s to be %#02x, got %#02x (opcode %#02x)",
				tc.Bank, tc.Name, tc.ExpectedValue, actualValue, tc.Opcode)
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
		Bank          byte
		Addr          uint8
		Ptr           uint8
		InitialValue  byte
		AddAmount     byte
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x26, Bank: 0, Addr: LOC_R0, Ptr: 0x30, InitialValue: 0xAA, AddAmount: 0x01, ExpectedValue: 0xAB},
		{Name: "@R0", Opcode: 0x26, Bank: 1, Addr: LOC_R0, Ptr: 0x30, InitialValue: 0xAA, AddAmount: 0x01, ExpectedValue: 0xAB},
		{Name: "@R0", Opcode: 0x26, Bank: 2, Addr: LOC_R0, Ptr: 0x30, InitialValue: 0xAA, AddAmount: 0x01, ExpectedValue: 0xAB},
		{Name: "@R0", Opcode: 0x26, Bank: 3, Addr: LOC_R0, Ptr: 0x30, InitialValue: 0xAA, AddAmount: 0x01, ExpectedValue: 0xAB},

		{Name: "@R1", Opcode: 0x27, Bank: 0, Addr: LOC_R1, Ptr: 0x30, InitialValue: 0xFF, AddAmount: 0x01, ExpectedValue: 0x00},
		{Name: "@R1", Opcode: 0x27, Bank: 1, Addr: LOC_R1, Ptr: 0x30, InitialValue: 0xFF, AddAmount: 0x01, ExpectedValue: 0x00},
		{Name: "@R1", Opcode: 0x27, Bank: 2, Addr: LOC_R1, Ptr: 0x30, InitialValue: 0xFF, AddAmount: 0x01, ExpectedValue: 0x00},
		{Name: "@R1", Opcode: 0x27, Bank: 3, Addr: LOC_R1, Ptr: 0x30, InitialValue: 0xFF, AddAmount: 0x01, ExpectedValue: 0x00},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.AddAmount); err != nil {
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
			t.Fatalf("expected bank %d reg %s to be %#02x, got %#02x (opcode %#02x)", tc.Bank, tc.Name, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x28_0x2f(t *testing.T) {
	cases := []struct {
		Name            string
		Opcode          byte
		Bank            byte
		Addr            byte
		AccInitialValue byte
		AddAmount       byte
		ExpectedValue   byte
	}{
		{Name: "R0", Opcode: 0x28, Bank: 0, Addr: LOC_R0, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R0", Opcode: 0x28, Bank: 1, Addr: LOC_R0, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R0", Opcode: 0x28, Bank: 2, Addr: LOC_R0, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R0", Opcode: 0x28, Bank: 3, Addr: LOC_R0, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R1", Opcode: 0x29, Bank: 0, Addr: LOC_R1, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R1", Opcode: 0x29, Bank: 1, Addr: LOC_R1, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R1", Opcode: 0x29, Bank: 2, Addr: LOC_R1, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R1", Opcode: 0x29, Bank: 3, Addr: LOC_R1, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R2", Opcode: 0x2A, Bank: 0, Addr: LOC_R2, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R2", Opcode: 0x2A, Bank: 1, Addr: LOC_R2, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R2", Opcode: 0x2A, Bank: 2, Addr: LOC_R2, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R2", Opcode: 0x2A, Bank: 3, Addr: LOC_R2, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R3", Opcode: 0x2B, Bank: 0, Addr: LOC_R3, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R3", Opcode: 0x2B, Bank: 1, Addr: LOC_R3, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R3", Opcode: 0x2B, Bank: 2, Addr: LOC_R3, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R3", Opcode: 0x2B, Bank: 3, Addr: LOC_R3, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R4", Opcode: 0x2C, Bank: 0, Addr: LOC_R4, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R4", Opcode: 0x2C, Bank: 1, Addr: LOC_R4, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R4", Opcode: 0x2C, Bank: 2, Addr: LOC_R4, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R4", Opcode: 0x2C, Bank: 3, Addr: LOC_R4, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R5", Opcode: 0x2D, Bank: 0, Addr: LOC_R5, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R5", Opcode: 0x2D, Bank: 1, Addr: LOC_R5, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R5", Opcode: 0x2D, Bank: 2, Addr: LOC_R5, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R5", Opcode: 0x2D, Bank: 3, Addr: LOC_R5, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R6", Opcode: 0x2E, Bank: 0, Addr: LOC_R6, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R6", Opcode: 0x2E, Bank: 1, Addr: LOC_R6, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R6", Opcode: 0x2E, Bank: 2, Addr: LOC_R6, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R6", Opcode: 0x2E, Bank: 3, Addr: LOC_R6, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},

		{Name: "R7", Opcode: 0x2F, Bank: 0, Addr: LOC_R7, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R7", Opcode: 0x2F, Bank: 1, Addr: LOC_R7, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R7", Opcode: 0x2F, Bank: 2, Addr: LOC_R7, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
		{Name: "R7", Opcode: 0x2F, Bank: 3, Addr: LOC_R7, AccInitialValue: 0x01 + 0, AddAmount: 0x02, ExpectedValue: 0x03 + 0},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(SFR_ACC, tc.AccInitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.AddAmount); err != nil {
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
		Bank          byte
		Addr          uint8
		Ptr           uint8
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x76, Bank: 0, Addr: LOC_R0, Ptr: 0x2F, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0x76, Bank: 1, Addr: LOC_R0, Ptr: 0x2F, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0x76, Bank: 2, Addr: LOC_R0, Ptr: 0x2F, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0x76, Bank: 3, Addr: LOC_R0, Ptr: 0x2F, ExpectedValue: 0xAA},

		{Name: "@R1", Opcode: 0x77, Bank: 0, Addr: LOC_R1, Ptr: 0x2F, ExpectedValue: 0xDD},
		{Name: "@R1", Opcode: 0x77, Bank: 1, Addr: LOC_R1, Ptr: 0x2F, ExpectedValue: 0xDD},
		{Name: "@R1", Opcode: 0x77, Bank: 2, Addr: LOC_R1, Ptr: 0x2F, ExpectedValue: 0xDD},
		{Name: "@R1", Opcode: 0x77, Bank: 3, Addr: LOC_R1, Ptr: 0x2F, ExpectedValue: 0xDD},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0x76, tc.ExpectedValue}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefBank(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected bank %d reg %s to be %#02x, got %#02x", tc.Bank, tc.Name, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x78_0x7f(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		Addr          uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x78, Bank: 0, Addr: LOC_R0, ExpectedValue: 0xAA + 0},
		{Name: "R0", Opcode: 0x78, Bank: 1, Addr: LOC_R0, ExpectedValue: 0xAA + 0},

		{Name: "R1", Opcode: 0x79, Bank: 0, Addr: LOC_R1, ExpectedValue: 0xAA + 1},
		{Name: "R1", Opcode: 0x79, Bank: 1, Addr: LOC_R1, ExpectedValue: 0xAA + 1},

		{Name: "R2", Opcode: 0x7a, Bank: 0, Addr: LOC_R2, ExpectedValue: 0xAA + 2},
		{Name: "R2", Opcode: 0x7a, Bank: 1, Addr: LOC_R2, ExpectedValue: 0xAA + 2},

		{Name: "R3", Opcode: 0x7b, Bank: 0, Addr: LOC_R3, ExpectedValue: 0xAA + 3},
		{Name: "R3", Opcode: 0x7b, Bank: 1, Addr: LOC_R3, ExpectedValue: 0xAA + 3},

		{Name: "R4", Opcode: 0x7c, Bank: 0, Addr: LOC_R4, ExpectedValue: 0xAA + 4},
		{Name: "R4", Opcode: 0x7c, Bank: 1, Addr: LOC_R4, ExpectedValue: 0xAA + 4},

		{Name: "R5", Opcode: 0x7d, Bank: 0, Addr: LOC_R5, ExpectedValue: 0xAA + 5},
		{Name: "R5", Opcode: 0x7d, Bank: 1, Addr: LOC_R5, ExpectedValue: 0xAA + 5},

		{Name: "R6", Opcode: 0x7e, Bank: 0, Addr: LOC_R6, ExpectedValue: 0xAA + 6},
		{Name: "R6", Opcode: 0x7e, Bank: 1, Addr: LOC_R6, ExpectedValue: 0xAA + 6},

		{Name: "R7", Opcode: 0x7f, Bank: 0, Addr: LOC_R7, ExpectedValue: 0xAA + 7},
		{Name: "R7", Opcode: 0x7f, Bank: 1, Addr: LOC_R7, ExpectedValue: 0xAA + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.ExpectedValue}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadBankMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected value in bank %d reg %s to be %#02x, got %#02x (opcode %#02x)",
				tc.Bank, tc.Name, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x86_0x87(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		MoveFrom      uint8
		Ptr           uint8
		MoveInto      uint8
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x86, Bank: 0, MoveFrom: LOC_R0, Ptr: 0x01, MoveInto: 0x02, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0x86, Bank: 1, MoveFrom: LOC_R0, Ptr: 0x01, MoveInto: 0x02, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0x86, Bank: 2, MoveFrom: LOC_R0, Ptr: 0x01, MoveInto: 0x02, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0x86, Bank: 3, MoveFrom: LOC_R0, Ptr: 0x01, MoveInto: 0x02, ExpectedValue: 0xAA},

		{Name: "@R1", Opcode: 0x87, Bank: 0, MoveFrom: LOC_R1, Ptr: 0x03, MoveInto: 0x04, ExpectedValue: 0xBB},
		{Name: "@R1", Opcode: 0x87, Bank: 1, MoveFrom: LOC_R1, Ptr: 0x03, MoveInto: 0x04, ExpectedValue: 0xBB},
		{Name: "@R1", Opcode: 0x87, Bank: 2, MoveFrom: LOC_R1, Ptr: 0x03, MoveInto: 0x04, ExpectedValue: 0xBB},
		{Name: "@R1", Opcode: 0x87, Bank: 3, MoveFrom: LOC_R1, Ptr: 0x03, MoveInto: 0x04, ExpectedValue: 0xBB},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.MoveFrom, tc.Ptr); err != nil {
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
			t.Errorf("expected bank %d address %#02x to be %#02x, got %#02x (opcode %#02x)",
				tc.Bank, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0x88_0x8f(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x88, Bank: 0, MoveFrom: LOC_R0, MoveInto: 0x00 + 0, ExpectedValue: 0xAA + 1},
		{Name: "R0", Opcode: 0x88, Bank: 1, MoveFrom: LOC_R0, MoveInto: 0x00 + 0, ExpectedValue: 0xAA + 1},
		{Name: "R0", Opcode: 0x88, Bank: 2, MoveFrom: LOC_R0, MoveInto: 0x00 + 0, ExpectedValue: 0xAA + 1},
		{Name: "R0", Opcode: 0x88, Bank: 3, MoveFrom: LOC_R0, MoveInto: 0x00 + 0, ExpectedValue: 0xAA + 1},

		{Name: "R1", Opcode: 0x89, Bank: 0, MoveFrom: LOC_R1, MoveInto: 0x00 + 1, ExpectedValue: 0xAA + 1},
		{Name: "R1", Opcode: 0x89, Bank: 1, MoveFrom: LOC_R1, MoveInto: 0x00 + 1, ExpectedValue: 0xAA + 1},
		{Name: "R1", Opcode: 0x89, Bank: 2, MoveFrom: LOC_R1, MoveInto: 0x00 + 1, ExpectedValue: 0xAA + 1},
		{Name: "R1", Opcode: 0x89, Bank: 3, MoveFrom: LOC_R1, MoveInto: 0x00 + 1, ExpectedValue: 0xAA + 1},

		{Name: "R2", Opcode: 0x8a, Bank: 0, MoveFrom: LOC_R2, MoveInto: 0x00 + 2, ExpectedValue: 0xAA + 1},
		{Name: "R2", Opcode: 0x8a, Bank: 1, MoveFrom: LOC_R2, MoveInto: 0x00 + 2, ExpectedValue: 0xAA + 1},
		{Name: "R2", Opcode: 0x8a, Bank: 2, MoveFrom: LOC_R2, MoveInto: 0x00 + 2, ExpectedValue: 0xAA + 1},
		{Name: "R2", Opcode: 0x8a, Bank: 3, MoveFrom: LOC_R2, MoveInto: 0x00 + 2, ExpectedValue: 0xAA + 1},

		{Name: "R3", Opcode: 0x8b, Bank: 0, MoveFrom: LOC_R3, MoveInto: 0x00 + 3, ExpectedValue: 0xAA + 1},
		{Name: "R3", Opcode: 0x8b, Bank: 1, MoveFrom: LOC_R3, MoveInto: 0x00 + 3, ExpectedValue: 0xAA + 1},
		{Name: "R3", Opcode: 0x8b, Bank: 2, MoveFrom: LOC_R3, MoveInto: 0x00 + 3, ExpectedValue: 0xAA + 1},
		{Name: "R3", Opcode: 0x8b, Bank: 3, MoveFrom: LOC_R3, MoveInto: 0x00 + 3, ExpectedValue: 0xAA + 1},

		{Name: "R4", Opcode: 0x8c, Bank: 0, MoveFrom: LOC_R4, MoveInto: 0x00 + 4, ExpectedValue: 0xAA + 1},
		{Name: "R4", Opcode: 0x8c, Bank: 1, MoveFrom: LOC_R4, MoveInto: 0x00 + 4, ExpectedValue: 0xAA + 1},
		{Name: "R4", Opcode: 0x8c, Bank: 2, MoveFrom: LOC_R4, MoveInto: 0x00 + 4, ExpectedValue: 0xAA + 1},
		{Name: "R4", Opcode: 0x8c, Bank: 3, MoveFrom: LOC_R4, MoveInto: 0x00 + 4, ExpectedValue: 0xAA + 1},

		{Name: "R5", Opcode: 0x8d, Bank: 0, MoveFrom: LOC_R5, MoveInto: 0x00 + 5, ExpectedValue: 0xAA + 1},
		{Name: "R5", Opcode: 0x8d, Bank: 1, MoveFrom: LOC_R5, MoveInto: 0x00 + 5, ExpectedValue: 0xAA + 1},
		{Name: "R5", Opcode: 0x8d, Bank: 2, MoveFrom: LOC_R5, MoveInto: 0x00 + 5, ExpectedValue: 0xAA + 1},
		{Name: "R5", Opcode: 0x8d, Bank: 3, MoveFrom: LOC_R5, MoveInto: 0x00 + 5, ExpectedValue: 0xAA + 1},

		{Name: "R6", Opcode: 0x8e, Bank: 0, MoveFrom: LOC_R6, MoveInto: 0x00 + 6, ExpectedValue: 0xAA + 1},
		{Name: "R6", Opcode: 0x8e, Bank: 1, MoveFrom: LOC_R6, MoveInto: 0x00 + 6, ExpectedValue: 0xAA + 1},
		{Name: "R6", Opcode: 0x8e, Bank: 2, MoveFrom: LOC_R6, MoveInto: 0x00 + 6, ExpectedValue: 0xAA + 1},
		{Name: "R6", Opcode: 0x8e, Bank: 3, MoveFrom: LOC_R6, MoveInto: 0x00 + 6, ExpectedValue: 0xAA + 1},

		{Name: "R7", Opcode: 0x8f, Bank: 0, MoveFrom: LOC_R7, MoveInto: 0x00 + 7, ExpectedValue: 0xAA + 1},
		{Name: "R7", Opcode: 0x8f, Bank: 1, MoveFrom: LOC_R7, MoveInto: 0x00 + 7, ExpectedValue: 0xAA + 1},
		{Name: "R7", Opcode: 0x8f, Bank: 2, MoveFrom: LOC_R7, MoveInto: 0x00 + 7, ExpectedValue: 0xAA + 1},
		{Name: "R7", Opcode: 0x8f, Bank: 3, MoveFrom: LOC_R7, MoveInto: 0x00 + 7, ExpectedValue: 0xAA + 1},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
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
			t.Errorf("%s: expected value at bank %d address %#02x to be %#02x, got %#02x (opcode %#02x)",
				tc.Name, tc.Bank, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xA6_0xA7(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		Dest          uint8
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		// R0 points to 0x01, addr 0x02 contains the data that we want (0xAA). deref'ing R0 should return 0xAA
		{Name: "@R0", Opcode: 0xA6, Bank: 0, Dest: LOC_R0, MoveInto: 0x20, MoveFrom: 0x10, ExpectedValue: 0xAA},
		{Name: "@R0", Opcode: 0xA6, Bank: 1, Dest: LOC_R0, MoveInto: 0x20, MoveFrom: 0x10, ExpectedValue: 0xAA},

		{Name: "@R1", Opcode: 0xA7, Bank: 0, Dest: LOC_R1, MoveInto: 0x20, MoveFrom: 0x10, ExpectedValue: 0xAB},
		{Name: "@R1", Opcode: 0xA7, Bank: 1, Dest: LOC_R1, MoveInto: 0x20, MoveFrom: 0x10, ExpectedValue: 0xAB},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Dest, tc.MoveInto); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveFrom}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefBank(tc.Dest)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected value at bank %d address %#02x to be %#02x, got %#02x (opcode %#02x)",
				tc.Name, tc.Bank, tc.Dest, tc.ExpectedValue, actualValue, tc.Opcode,
			)
		}
	}
}

func TestOp0xA8_0xAF(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0xA8, Bank: 0, MoveInto: LOC_R0, MoveFrom: 0xAA + 0, ExpectedValue: 0xBB + 0},
		{Name: "R0", Opcode: 0xA8, Bank: 1, MoveInto: LOC_R0, MoveFrom: 0xAA + 0, ExpectedValue: 0xBB + 0},
		{Name: "R1", Opcode: 0xA9, Bank: 0, MoveInto: LOC_R1, MoveFrom: 0xAA + 1, ExpectedValue: 0xBB + 1},
		{Name: "R1", Opcode: 0xA9, Bank: 1, MoveInto: LOC_R1, MoveFrom: 0xAA + 1, ExpectedValue: 0xBB + 1},
		{Name: "R2", Opcode: 0xAA, Bank: 0, MoveInto: LOC_R2, MoveFrom: 0xAA + 2, ExpectedValue: 0xBB + 2},
		{Name: "R2", Opcode: 0xAA, Bank: 1, MoveInto: LOC_R2, MoveFrom: 0xAA + 2, ExpectedValue: 0xBB + 2},
		{Name: "R3", Opcode: 0xAB, Bank: 0, MoveInto: LOC_R3, MoveFrom: 0xAA + 3, ExpectedValue: 0xBB + 3},
		{Name: "R3", Opcode: 0xAB, Bank: 1, MoveInto: LOC_R3, MoveFrom: 0xAA + 3, ExpectedValue: 0xBB + 3},
		{Name: "R4", Opcode: 0xAC, Bank: 0, MoveInto: LOC_R4, MoveFrom: 0xAA + 4, ExpectedValue: 0xBB + 4},
		{Name: "R4", Opcode: 0xAC, Bank: 1, MoveInto: LOC_R4, MoveFrom: 0xAA + 4, ExpectedValue: 0xBB + 4},
		{Name: "R5", Opcode: 0xAD, Bank: 0, MoveInto: LOC_R5, MoveFrom: 0xAA + 5, ExpectedValue: 0xBB + 5},
		{Name: "R5", Opcode: 0xAD, Bank: 1, MoveInto: LOC_R5, MoveFrom: 0xAA + 5, ExpectedValue: 0xBB + 5},
		{Name: "R6", Opcode: 0xAE, Bank: 0, MoveInto: LOC_R6, MoveFrom: 0xAA + 6, ExpectedValue: 0xBB + 6},
		{Name: "R6", Opcode: 0xAE, Bank: 1, MoveInto: LOC_R6, MoveFrom: 0xAA + 6, ExpectedValue: 0xBB + 6},
		{Name: "R7", Opcode: 0xAF, Bank: 0, MoveInto: LOC_R7, MoveFrom: 0xAA + 7, ExpectedValue: 0xBB + 7},
		{Name: "R7", Opcode: 0xAF, Bank: 1, MoveInto: LOC_R7, MoveFrom: 0xAA + 7, ExpectedValue: 0xBB + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode, tc.MoveFrom}); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadBankMem(tc.MoveInto)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("%s: expected bank %d address %#02x to be %#02x, got %#02x (opcode %#02x)",
				tc.Name, tc.Bank, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xE6_0xE7(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		MoveFrom      uint8
		Ptr           uint8
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0xE6, Bank: 0, MoveFrom: LOC_R0, Ptr: 0xAA, ExpectedValue: 0xDD},
		{Name: "@R0", Opcode: 0xE6, Bank: 1, MoveFrom: LOC_R0, Ptr: 0xAA, ExpectedValue: 0xDD},

		{Name: "@R1", Opcode: 0xE7, Bank: 0, MoveFrom: LOC_R1, Ptr: 0xAB, ExpectedValue: 0xDD},
		{Name: "@R1", Opcode: 0xE7, Bank: 1, MoveFrom: LOC_R1, Ptr: 0xAB, ExpectedValue: 0xDD},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.MoveFrom, tc.Ptr); err != nil {
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
			t.Errorf("%s: expected bank %d A register to be %#02x, got %#02x (opcode %#02x)",
				tc.Name, tc.Bank, tc.ExpectedValue, actualValue, tc.Opcode)
		}
	}
}

func TestOp0xE5(t *testing.T) {
	// {Name: "ramaddr", Opcode: 0xE5, MoveInto: SFR_ACC, MoveFrom: 0xAD, ExpectedValue: 0xBB + 0},
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xE8_0xEF(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		MoveInto      uint8
		MoveFrom      uint8
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0xE8, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R0, ExpectedValue: 0xBB + 0},
		{Name: "R0", Opcode: 0xE8, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R0, ExpectedValue: 0xBB + 0},
		{Name: "R1", Opcode: 0xE9, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R1, ExpectedValue: 0xBB + 1},
		{Name: "R1", Opcode: 0xE9, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R1, ExpectedValue: 0xBB + 1},
		{Name: "R2", Opcode: 0xEA, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R2, ExpectedValue: 0xBB + 2},
		{Name: "R2", Opcode: 0xEA, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R2, ExpectedValue: 0xBB + 2},
		{Name: "R3", Opcode: 0xEB, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R3, ExpectedValue: 0xBB + 3},
		{Name: "R3", Opcode: 0xEB, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R3, ExpectedValue: 0xBB + 3},
		{Name: "R4", Opcode: 0xEC, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R4, ExpectedValue: 0xBB + 4},
		{Name: "R4", Opcode: 0xEC, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R4, ExpectedValue: 0xBB + 4},
		{Name: "R5", Opcode: 0xED, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R5, ExpectedValue: 0xBB + 5},
		{Name: "R5", Opcode: 0xED, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R5, ExpectedValue: 0xBB + 5},
		{Name: "R6", Opcode: 0xEE, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R6, ExpectedValue: 0xBB + 6},
		{Name: "R6", Opcode: 0xEE, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R6, ExpectedValue: 0xBB + 6},
		{Name: "R7", Opcode: 0xEF, Bank: 0, MoveInto: SFR_ACC, MoveFrom: LOC_R7, ExpectedValue: 0xBB + 7},
		{Name: "R7", Opcode: 0xEF, Bank: 1, MoveInto: SFR_ACC, MoveFrom: LOC_R7, ExpectedValue: 0xBB + 7},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.MoveFrom, tc.ExpectedValue); err != nil {
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
			t.Errorf("%s: expected bank %d address %#02x to be %#02x, got %#02x (opcode %#02x)",
				tc.Name, tc.Bank, tc.MoveInto, tc.ExpectedValue, actualValue, tc.Opcode)
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

func TestOp0x42(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x43(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x44(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x45(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x46_0x47(t *testing.T) {
	cases := []struct {
		Name      string
		Opcode    byte
		Addr      uint8
		InitialA  byte
		Bank      byte
		Ptr       uint8
		RegValue  byte
		ExpectedA byte
	}{
		{Name: "@R0", Opcode: 0x46, Addr: LOC_R0, InitialA: 0b11110000, Bank: 0, Ptr: 0x2F, RegValue: 0b00001111, ExpectedA: 0b11111111},
		{Name: "@R0", Opcode: 0x46, Addr: LOC_R0, InitialA: 0b11110000, Bank: 1, Ptr: 0x2F, RegValue: 0b00001111, ExpectedA: 0b11111111},
		{Name: "@R1", Opcode: 0x47, Addr: LOC_R1, InitialA: 0b11110000, Bank: 0, Ptr: 0x2F, RegValue: 0b00001111, ExpectedA: 0b11111111},
		{Name: "@R1", Opcode: 0x47, Addr: LOC_R1, InitialA: 0b11110000, Bank: 1, Ptr: 0x2F, RegValue: 0b00001111, ExpectedA: 0b11111111},
	}

	for _, tc := range cases {
		var initialA byte = tc.InitialA
		var ptr byte = tc.Ptr
		var val byte = tc.RegValue
		var expectedA byte = tc.ExpectedA
		var bankNo byte = tc.Bank

		vm := NewMachine()

		if err := vm.SetBankNo(bankNo); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(SFR_ACC, initialA); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.SetrefBank(tc.Addr, val); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{tc.Opcode}); err != nil {
			t.Fatal(err)
		}

		actualA, _ := vm.ReadMem(SFR_ACC)

		if actualA != expectedA {
			t.Errorf("%s: expected bank %d A register to be %#08b, got %#08b", tc.Name, tc.Bank, expectedA, actualA)
		}
	}
}

func TestOp0x48(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		Addr          byte
		InitialValue  byte
		AddrValue     byte
		ExpectedValue byte
	}{
		{Name: "R0", Opcode: 0x48, Bank: 0, Addr: LOC_R0, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R0", Opcode: 0x48, Bank: 1, Addr: LOC_R0, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R1", Opcode: 0x49, Bank: 0, Addr: LOC_R1, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R1", Opcode: 0x49, Bank: 1, Addr: LOC_R1, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R2", Opcode: 0x4a, Bank: 0, Addr: LOC_R2, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R2", Opcode: 0x4a, Bank: 1, Addr: LOC_R2, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R3", Opcode: 0x4b, Bank: 0, Addr: LOC_R3, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R3", Opcode: 0x4b, Bank: 1, Addr: LOC_R3, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R4", Opcode: 0x4c, Bank: 0, Addr: LOC_R4, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R4", Opcode: 0x4c, Bank: 1, Addr: LOC_R4, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R5", Opcode: 0x4d, Bank: 0, Addr: LOC_R5, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R5", Opcode: 0x4d, Bank: 1, Addr: LOC_R5, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R6", Opcode: 0x4e, Bank: 0, Addr: LOC_R6, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R6", Opcode: 0x4e, Bank: 1, Addr: LOC_R6, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R7", Opcode: 0x4f, Bank: 0, Addr: LOC_R7, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
		{Name: "R7", Opcode: 0x4f, Bank: 1, Addr: LOC_R7, InitialValue: 0b11110000, AddrValue: 0b00001111, ExpectedValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.AddrValue); err != nil {
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
			t.Errorf("opcode %#02x: expected bank %d register %s value to be %#08b, got %#08b",
				tc.Opcode, tc.Bank, tc.Name, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x49(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x4a(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x4b(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x4c(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x4d(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x4e(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x4f(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x52(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x53(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x54(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x55(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x56_0x57(t *testing.T) {
	cases := []struct {
		Name          string
		Opcode        byte
		Bank          byte
		Addr          uint8
		Ptr           uint8
		InitialValue  byte
		PtrValue      byte
		ExpectedValue byte
	}{
		{Name: "@R0", Opcode: 0x56, Bank: 0, Addr: LOC_R0, Ptr: 0x2F, InitialValue: 0xFF, PtrValue: 0x00, ExpectedValue: 0x00},
		{Name: "@R0", Opcode: 0x56, Bank: 1, Addr: LOC_R0, Ptr: 0x2F, InitialValue: 0xFF, PtrValue: 0x00, ExpectedValue: 0x00},
		{Name: "@R1", Opcode: 0x57, Bank: 0, Addr: LOC_R1, Ptr: 0x2F, InitialValue: 0xFF, PtrValue: 0x00, ExpectedValue: 0x00},
		{Name: "@R1", Opcode: 0x57, Bank: 1, Addr: LOC_R1, Ptr: 0x2F, InitialValue: 0xFF, PtrValue: 0x00, ExpectedValue: 0x00},
	}

	for _, tc := range cases {
		vm := NewMachine()

		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.PtrValue); err != nil {
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
			t.Errorf("%s: expected bank %d register %#02x to be %#08b, got %#08b",
				tc.Name, tc.Bank, tc.Addr, tc.ExpectedValue, actualValue)
		}
	}
}

func TestOp0x58(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x59(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x5a(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x5b(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x5c(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x5d(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x5e(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x5f(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x62(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x63(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x64(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x65(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x66(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x67(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x68(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x69(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x6a(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x6b(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x6c(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x6d(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x6e(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x6f(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0x85(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xC4(t *testing.T) {
	cases := []struct {
		InitialValue  byte
		ExpectedValue byte
	}{
		{InitialValue: 0x00, ExpectedValue: 0x00},
		{InitialValue: 0xFF, ExpectedValue: 0xFF},
		{InitialValue: 0b11110000, ExpectedValue: 0b00001111},
		{InitialValue: 0b00001111, ExpectedValue: 0b11110000},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(SFR_ACC, tc.InitialValue); err != nil {
			t.Fatal(err)
		}

		if err := vm.Feed([]byte{0xC4}); err != nil {
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

func TestOp0xC5(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xC6(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xC7(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xC8(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xC9(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xCA(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xCB(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xCC(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xCD(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xCE(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xCF(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xD6(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xD7(t *testing.T) {
	// TODO: implement
	t.Skipf("TODO: implement")
}

func TestOp0xE4(t *testing.T) {
	var initialValue byte = 0xFF
	var expectedValue byte = 0x00

	vm := NewMachine()
	if err := vm.WriteMem(SFR_ACC, initialValue); err != nil {
		t.Fatal(err)
	}

	if err := vm.Feed([]byte{0xe4}); err != nil {
		t.Fatal(err)
	}

	actualValue, err := vm.ReadMem(SFR_ACC)
	if err != nil {
		t.Fatal(err)
	}

	if actualValue != expectedValue {
		t.Errorf("expected A register to be %#08b, got %#08b", expectedValue, actualValue)
	}
}

func TestOp0xC0_0xD0(t *testing.T) {
	var srcAddr byte = 0x2F
	var expectedValue byte = 0xFF

	vm := NewMachine()
	err := vm.WriteMem(srcAddr, expectedValue)
	if err != nil {
		t.Fatal(err)
	}
	currentSp := vm.SP
	expectedNewSp := currentSp + 1

	if err := vm.Feed([]byte{0xc0, srcAddr}); err != nil {
		t.Fatal(err)
	}

	newSp := vm.SP

	if newSp != expectedNewSp {
		t.Errorf("expected SP to be %d, got %d", expectedNewSp, newSp)
	}

	spVal, err := vm.ReadMem(newSp)
	if err != nil {
		t.Fatal(err)
	}

	if spVal != expectedValue {
		t.Errorf("direct read: expected value at SP (%d) to be %#02x, got %#02x", newSp, expectedValue, spVal)
	}

	var popInto byte = 0x00

	if err := vm.Feed([]byte{0xd0, popInto}); err != nil {
		t.Fatal(err)
	}

	poppedValue, err := vm.ReadMem(popInto)
	if err != nil {
		t.Fatal(err)
	}

	if poppedValue != expectedValue {
		t.Errorf("pop: expected POPped value to be %#02x, got %#02x", expectedValue, poppedValue)
	}
}
