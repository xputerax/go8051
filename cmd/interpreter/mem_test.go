package main

import (
	"testing"
)

func TestMemWriteSFR(t *testing.T) {
	cases := []struct {
		RegisterName    string
		Location        uint8
		Expected        byte
		ReadRegDirectly func(*Machine) byte
		ReadMem         func(*Machine) byte
	}{
		{RegisterName: "ACC", Location: SFR_ACC, Expected: SFR_ACC, ReadRegDirectly: func(m *Machine) byte { return m.registers.ACC }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_ACC); return val }},
		{RegisterName: "B", Location: SFR_B, Expected: SFR_B, ReadRegDirectly: func(m *Machine) byte { return m.registers.B }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_B); return val }},
		{RegisterName: "DPH", Location: SFR_DPH, Expected: SFR_DPH, ReadRegDirectly: func(m *Machine) byte { return m.registers.DPH }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_DPH); return val }},
		{RegisterName: "DPL", Location: SFR_DPL, Expected: SFR_DPL, ReadRegDirectly: func(m *Machine) byte { return m.registers.DPL }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_DPL); return val }},
		{RegisterName: "IE", Location: SFR_IE, Expected: SFR_IE, ReadRegDirectly: func(m *Machine) byte { return m.registers.IE }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_IE); return val }},
		{RegisterName: "IP", Location: SFR_IP, Expected: SFR_IP, ReadRegDirectly: func(m *Machine) byte { return m.registers.IP }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_IP); return val }},
		{RegisterName: "P0", Location: SFR_P0, Expected: SFR_P0, ReadRegDirectly: func(m *Machine) byte { return m.registers.P0 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_P0); return val }},
		{RegisterName: "P1", Location: SFR_P1, Expected: SFR_P1, ReadRegDirectly: func(m *Machine) byte { return m.registers.P1 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_P1); return val }},
		{RegisterName: "P2", Location: SFR_P2, Expected: SFR_P2, ReadRegDirectly: func(m *Machine) byte { return m.registers.P2 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_P2); return val }},
		{RegisterName: "P3", Location: SFR_P3, Expected: SFR_P3, ReadRegDirectly: func(m *Machine) byte { return m.registers.P3 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_P3); return val }},
		{RegisterName: "PCON", Location: SFR_PCON, Expected: SFR_PCON, ReadRegDirectly: func(m *Machine) byte { return m.registers.PCON }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_PCON); return val }},
		{RegisterName: "PSW", Location: SFR_PSW, Expected: SFR_PSW, ReadRegDirectly: func(m *Machine) byte { return m.registers.PSW }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_PSW); return val }},
		{RegisterName: "SCON", Location: SFR_SCON, Expected: SFR_SCON, ReadRegDirectly: func(m *Machine) byte { return m.registers.SCON }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_SCON); return val }},
		{RegisterName: "SBUF", Location: SFR_SBUF, Expected: SFR_SBUF, ReadRegDirectly: func(m *Machine) byte { return m.registers.SBUF }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_SBUF); return val }},
		{RegisterName: "SP", Location: SFR_SP, Expected: SFR_SP, ReadRegDirectly: func(m *Machine) byte { return m.registers.SP }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_SP); return val }},
		{RegisterName: "TMOD", Location: SFR_TMOD, Expected: SFR_TMOD, ReadRegDirectly: func(m *Machine) byte { return m.registers.TMOD }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_TMOD); return val }},
		{RegisterName: "TCON", Location: SFR_TCON, Expected: SFR_TCON, ReadRegDirectly: func(m *Machine) byte { return m.registers.TCON }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_TCON); return val }},
		{RegisterName: "TL0", Location: SFR_TL0, Expected: SFR_TL0, ReadRegDirectly: func(m *Machine) byte { return m.registers.TL0 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_TL0); return val }},
		{RegisterName: "TH0", Location: SFR_TH0, Expected: SFR_TH0, ReadRegDirectly: func(m *Machine) byte { return m.registers.TH0 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_TH0); return val }},
		{RegisterName: "TL1", Location: SFR_TL1, Expected: SFR_TL1, ReadRegDirectly: func(m *Machine) byte { return m.registers.TL1 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_TL1); return val }},
		{RegisterName: "TH1", Location: SFR_TH1, Expected: SFR_TH1, ReadRegDirectly: func(m *Machine) byte { return m.registers.TH1 }, ReadMem: func(m *Machine) byte { val, _ := m.ReadMem(SFR_TH1); return val }},
	}

	for _, tc := range cases {
		vm := NewMachine()
		err := vm.WriteMem(tc.Location, tc.Expected)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		actualReg := tc.ReadRegDirectly(vm)
		if tc.Expected != actualReg {
			t.Errorf("expected %#02x in register %s, got %#02x", tc.Expected, tc.RegisterName, actualReg)
		}

		actualMem := tc.ReadMem(vm)
		if tc.Expected != actualMem {
			t.Errorf("expected %#02x in memory location %#02x, got %#02x", tc.Expected, tc.Location, actualMem)
		}
	}
}

func TestWriteMemAllOk(t *testing.T) {
	vm := NewMachine()

	var i uint8
	for i = 0; i < 0xFF; i++ {
		err := vm.WriteMem(i, 0xAD)
		if err != nil {
			t.Errorf("unexpected error when writing to address %d (%#02x): %s", i, i, err)
		}
		if i == 0xFF {
			break
		}
	}
}

func TestDerefMem(t *testing.T) {
	cases := []struct {
		Addr     uint8
		Ptr      uint8
		PtrValue byte
	}{
		{Addr: 0x00, Ptr: 0x01, PtrValue: 0xAD},
		{Addr: 0x01, Ptr: 0x02, PtrValue: 0xFF},
		{Addr: 0x7E, Ptr: 0x7F, PtrValue: 0x69},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.WriteMem(tc.Addr, tc.Ptr); err != nil {
			t.Fatalf("failed to write to memory address %#02x: %s", tc.Addr, err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.PtrValue); err != nil {
			t.Fatalf("failed to write to memory address %#02x: %s", tc.Ptr, err)
		}

		if actualValue, err := vm.DerefMem(tc.Addr); err != nil {
			t.Fatalf("deref error: %s", err)
		} else {
			if actualValue != tc.PtrValue {
				t.Fatalf("expected value at address %#02x to be %#02b, got %#02b", tc.Ptr, tc.PtrValue, actualValue)
			}
		}
	}
}

func TestSetrefMem(t *testing.T) {
	cases := []struct {
		Addr          uint8
		Ptr           uint8
		ExpectedValue byte
	}{
		{Addr: 0, Ptr: 1, ExpectedValue: 0xAD},
		{Addr: 1, Ptr: 0, ExpectedValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()
		vm.WriteMem(tc.Addr, tc.Ptr)
		vm.SetrefMem(tc.Addr, tc.ExpectedValue)

		directRead, err := vm.ReadMem(tc.Ptr)
		if err != nil {
			t.Fatal(err)
		}

		indirectRead, err := vm.DerefMem(tc.Addr)
		if err != nil {
			t.Fatal(err)
		}

		if directRead != indirectRead {
			t.Fatalf("directRead and indirectRead returns different value: %#02x, %#02x", directRead, indirectRead)
		}

		if directRead != tc.ExpectedValue {
			t.Fatalf("expected directRead to be %#02x, got %#02x", tc.ExpectedValue, directRead)
		}

		if indirectRead != tc.ExpectedValue {
			t.Fatalf("expected indirectRead to be %#02x, got %#02x", tc.ExpectedValue, indirectRead)
		}
	}
}

func TestBankNo_BankOffset(t *testing.T) {
	cases := []struct {
		RS1            byte
		RS0            byte
		BankNo         byte
		ExpectedOffset byte
	}{
		{RS1: 0, RS0: 0, BankNo: 0, ExpectedOffset: 0 * BANK_SIZE},
		{RS1: 0, RS0: 1, BankNo: 1, ExpectedOffset: 1 * BANK_SIZE},
		{RS1: 1, RS0: 0, BankNo: 2, ExpectedOffset: 2 * BANK_SIZE},
		{RS1: 1, RS0: 1, BankNo: 3, ExpectedOffset: 3 * BANK_SIZE},
	}

	for _, tc := range cases {
		vm := NewMachine()
		psw, _ := vm.ReadMem(SFR_PSW)

		if tc.RS0 == 0 {
			psw = PSW_UNSET(psw, PSW_RS0_MASK)
		} else if tc.RS0 == 1 {
			psw = PSW_SET(psw, PSW_RS0_MASK)
		} else {
			t.Fatalf("invalid value for RS0: %d", tc.RS0)
		}

		if tc.RS1 == 0 {
			psw = PSW_UNSET(psw, PSW_RS1_MASK)
		} else if tc.RS1 == 1 {
			psw = PSW_SET(psw, PSW_RS1_MASK)
		} else {
			t.Fatalf("invalid value for RS1: %d", tc.RS1)
		}

		err := vm.WriteMem(SFR_PSW, psw)
		if err != nil {
			t.Fatal(err)
		}

		actualBankNo := vm.bankNo()
		actualOffset := vm.bankOffset()

		if actualBankNo != tc.BankNo {
			t.Errorf("RS1: %d, RS0: %d: expected bank number to be %d, got %d", tc.RS1, tc.RS0, tc.BankNo, actualBankNo)
		}

		if actualOffset != tc.ExpectedOffset {
			t.Errorf("RS1: %d, RS0: %d: expected bank %d offset to be %d, got %d", tc.RS1, tc.RS0, tc.BankNo, tc.ExpectedOffset, actualOffset)
		}
	}
}

func TestSetBankNo(t *testing.T) {
	cases := []struct {
		expected byte
	}{
		{expected: 0},
		{expected: 1},
		{expected: 2},
		{expected: 3},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.SetBankNo(tc.expected); err != nil {
			t.Fatal(err)
		}

		actual := vm.bankNo()

		if actual != tc.expected {
			t.Errorf("expected bank no to be %d, got %d", tc.expected, actual)
		}
	}
}

func TestSetBankNoErr(t *testing.T) {
	vm := NewMachine()
	err := vm.SetBankNo(4)
	if err == nil {
		t.Errorf("expected SetBankNo to return error, nil given")
	}
}

func TestRead_WriteBankMem(t *testing.T) {
	cases := []struct {
		Name          string
		Reg           uint8
		Bank          byte
		ExpectedValue byte
	}{
		{Name: "R0", Reg: LOC_R0, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R0", Reg: LOC_R0, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R0", Reg: LOC_R0, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R0", Reg: LOC_R0, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R1", Reg: LOC_R1, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R1", Reg: LOC_R1, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R1", Reg: LOC_R1, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R1", Reg: LOC_R1, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R2", Reg: LOC_R2, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R2", Reg: LOC_R2, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R2", Reg: LOC_R2, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R2", Reg: LOC_R2, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R3", Reg: LOC_R3, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R3", Reg: LOC_R3, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R3", Reg: LOC_R3, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R3", Reg: LOC_R3, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R4", Reg: LOC_R4, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R4", Reg: LOC_R4, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R4", Reg: LOC_R4, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R4", Reg: LOC_R4, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R5", Reg: LOC_R5, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R5", Reg: LOC_R5, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R5", Reg: LOC_R5, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R5", Reg: LOC_R5, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R6", Reg: LOC_R6, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R6", Reg: LOC_R6, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R6", Reg: LOC_R6, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R6", Reg: LOC_R6, Bank: 3, ExpectedValue: 0xFF},

		{Name: "R7", Reg: LOC_R7, Bank: 0, ExpectedValue: 0xFF},
		{Name: "R7", Reg: LOC_R7, Bank: 1, ExpectedValue: 0xFF},
		{Name: "R7", Reg: LOC_R7, Bank: 2, ExpectedValue: 0xFF},
		{Name: "R7", Reg: LOC_R7, Bank: 3, ExpectedValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Reg, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.ReadBankMem(tc.Reg)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected value in %s bank %d to be %#02x, got %#02x", tc.Name, tc.Bank, tc.ExpectedValue, actualValue)
		}
	}
}

func TestDerefBank(t *testing.T) {
	cases := []struct {
		Name          string
		Reg           uint8
		Ptr           uint8
		Bank          byte
		ExpectedValue byte
	}{
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Reg, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteMem(tc.Ptr, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefBank(tc.Reg)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected bank %d reg %s value to be %#02x, got %#02x", tc.Bank, tc.Name, tc.ExpectedValue, actualValue)
		}
	}
}

func TestSetrefBank(t *testing.T) {
	cases := []struct {
		Name          string
		Reg           uint8
		Ptr           uint8
		Bank          byte
		ExpectedValue byte
	}{
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R0", Reg: LOC_R0, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R1", Reg: LOC_R1, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R2", Reg: LOC_R2, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R3", Reg: LOC_R3, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R4", Reg: LOC_R4, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R5", Reg: LOC_R5, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R6", Reg: LOC_R6, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},

		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 0, ExpectedValue: 0xFF},
		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 1, ExpectedValue: 0xFF},
		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 2, ExpectedValue: 0xFF},
		{Name: "@R7", Reg: LOC_R7, Ptr: 0x2f, Bank: 3, ExpectedValue: 0xFF},
	}

	for _, tc := range cases {
		vm := NewMachine()
		if err := vm.SetBankNo(tc.Bank); err != nil {
			t.Fatal(err)
		}

		if err := vm.WriteBankMem(tc.Reg, tc.Ptr); err != nil {
			t.Fatal(err)
		}

		if err := vm.SetrefBank(tc.Reg, tc.ExpectedValue); err != nil {
			t.Fatal(err)
		}

		actualValue, err := vm.DerefBank(tc.Reg)
		if err != nil {
			t.Fatal(err)
		}

		if actualValue != tc.ExpectedValue {
			t.Errorf("expected bank %d reg %s value to be %#02x, got %#02x", tc.Bank, tc.Name, tc.ExpectedValue, actualValue)
		}
	}
}
