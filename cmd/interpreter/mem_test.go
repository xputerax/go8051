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
		{RegisterName: "ACC", Location: SFR_ACC, Expected: SFR_ACC, ReadRegDirectly: func(m *Machine) byte { return m.Registers.ACC }, ReadMem: func(m *Machine) byte { return m.Data[SFR_ACC] }},
		{RegisterName: "B", Location: SFR_B, Expected: SFR_B, ReadRegDirectly: func(m *Machine) byte { return m.Registers.B }, ReadMem: func(m *Machine) byte { return m.Data[SFR_B] }},
		{RegisterName: "DPH", Location: SFR_DPH, Expected: SFR_DPH, ReadRegDirectly: func(m *Machine) byte { return m.Registers.DPH }, ReadMem: func(m *Machine) byte { return m.Data[SFR_DPH] }},
		{RegisterName: "DPL", Location: SFR_DPL, Expected: SFR_DPL, ReadRegDirectly: func(m *Machine) byte { return m.Registers.DPL }, ReadMem: func(m *Machine) byte { return m.Data[SFR_DPL] }},
		{RegisterName: "IE", Location: SFR_IE, Expected: SFR_IE, ReadRegDirectly: func(m *Machine) byte { return m.Registers.IE }, ReadMem: func(m *Machine) byte { return m.Data[SFR_IE] }},
		{RegisterName: "IP", Location: SFR_IP, Expected: SFR_IP, ReadRegDirectly: func(m *Machine) byte { return m.Registers.IP }, ReadMem: func(m *Machine) byte { return m.Data[SFR_IP] }},
		{RegisterName: "P0", Location: SFR_P0, Expected: SFR_P0, ReadRegDirectly: func(m *Machine) byte { return m.Registers.P0 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_P0] }},
		{RegisterName: "P1", Location: SFR_P1, Expected: SFR_P1, ReadRegDirectly: func(m *Machine) byte { return m.Registers.P1 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_P1] }},
		{RegisterName: "P2", Location: SFR_P2, Expected: SFR_P2, ReadRegDirectly: func(m *Machine) byte { return m.Registers.P2 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_P2] }},
		{RegisterName: "P3", Location: SFR_P3, Expected: SFR_P3, ReadRegDirectly: func(m *Machine) byte { return m.Registers.P3 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_P3] }},
		{RegisterName: "PCON", Location: SFR_PCON, Expected: SFR_PCON, ReadRegDirectly: func(m *Machine) byte { return m.Registers.PCON }, ReadMem: func(m *Machine) byte { return m.Data[SFR_PCON] }},
		{RegisterName: "PSW", Location: SFR_PSW, Expected: SFR_PSW, ReadRegDirectly: func(m *Machine) byte { return m.Registers.PSW }, ReadMem: func(m *Machine) byte { return m.Data[SFR_PSW] }},
		{RegisterName: "SCON", Location: SFR_SCON, Expected: SFR_SCON, ReadRegDirectly: func(m *Machine) byte { return m.Registers.SCON }, ReadMem: func(m *Machine) byte { return m.Data[SFR_SCON] }},
		{RegisterName: "SBUF", Location: SFR_SBUF, Expected: SFR_SBUF, ReadRegDirectly: func(m *Machine) byte { return m.Registers.SBUF }, ReadMem: func(m *Machine) byte { return m.Data[SFR_SBUF] }},
		{RegisterName: "SP", Location: SFR_SP, Expected: SFR_SP, ReadRegDirectly: func(m *Machine) byte { return m.Registers.SP }, ReadMem: func(m *Machine) byte { return m.Data[SFR_SP] }},
		{RegisterName: "TMOD", Location: SFR_TMOD, Expected: SFR_TMOD, ReadRegDirectly: func(m *Machine) byte { return m.Registers.TMOD }, ReadMem: func(m *Machine) byte { return m.Data[SFR_TMOD] }},
		{RegisterName: "TCON", Location: SFR_TCON, Expected: SFR_TCON, ReadRegDirectly: func(m *Machine) byte { return m.Registers.TCON }, ReadMem: func(m *Machine) byte { return m.Data[SFR_TCON] }},
		{RegisterName: "TL0", Location: SFR_TL0, Expected: SFR_TL0, ReadRegDirectly: func(m *Machine) byte { return m.Registers.TL0 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_TL0] }},
		{RegisterName: "TH0", Location: SFR_TH0, Expected: SFR_TH0, ReadRegDirectly: func(m *Machine) byte { return m.Registers.TH0 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_TH0] }},
		{RegisterName: "TL1", Location: SFR_TL1, Expected: SFR_TL1, ReadRegDirectly: func(m *Machine) byte { return m.Registers.TL1 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_TL1] }},
		{RegisterName: "TH1", Location: SFR_TH1, Expected: SFR_TH1, ReadRegDirectly: func(m *Machine) byte { return m.Registers.TH1 }, ReadMem: func(m *Machine) byte { return m.Data[SFR_TH1] }},
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
