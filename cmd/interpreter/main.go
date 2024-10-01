package main

import (
	"fmt"
	"log"
)

/** Special function registers - 80h - FFh */
type Register struct {
	ACC  byte // A Register (for arithmetic) E0H
	B    byte // B Register (for arithmetic) F0H
	DPH  byte // Addressing External Memory 83H
	DPL  byte // Addressing External Memory 82H
	IE   byte // Interrupt Enable Control A8H
	IP   byte // Interrupt Priority B8H
	P0   byte // Port 0 Latch 80H
	P1   byte // Port 1 Latch 90H
	P2   byte // Port 2 Latch A0H
	P3   byte // Port 3 Latch B0H
	PCON byte // Power Control 87H
	PSW  byte // Program Status Word D0H
	SCON byte // Serial Port Control 98H
	SBUF byte // Serial Port Data Buffer 99H
	SP   byte // Stack Pointer 81H
	TMOD byte // Timer / Counter Mode Control 89H
	TCON byte // Timer / Counter Control 88H
	TL0  byte // Timer 0 LOW Byte 8AH
	TH0  byte // Timer 0 HIGH Byte 8CH
	TL1  byte // Timer 1 LOW Byte 8BH
	TH1  byte // Timer 1 HIGH Byte 8DH
}

const PSW_C_MASK byte = (1 << 7)
const PSW_AC_MASK byte = (1 << 6)
const PSW_F0_MASK byte = (1 << 5)
const PSW_RS1_MASK byte = (1 << 4)
const PSW_RS0_MASK byte = (1 << 3)
const PSW_OV_MASK byte = (1 << 2)
const PSW_RESERVED_MASK byte = (1 << 1)
const PSW_P_MASK byte = (1 << 0)

// Carry bit
func PSW_C(psw byte) bool {
	return (psw & PSW_C_MASK) == PSW_C_MASK
}

// Aux carry bit
func PSW_AC(psw byte) bool {
	return (psw & PSW_AC_MASK) == PSW_AC_MASK
}

// Flag 0
func PSW_F0(psw byte) bool {
	return (psw & PSW_F0_MASK) == PSW_F0_MASK
}

// Register bank selection 1
func PSW_RS1(psw byte) bool {
	return (psw & PSW_RS1_MASK) == PSW_RS1_MASK
}

// Register bank selection 0
func PSW_RS0(psw byte) bool {
	return (psw & PSW_RS0_MASK) == PSW_RS0_MASK
}

// Overflow
func PSW_OV(psw byte) bool {
	return (psw & PSW_OV_MASK) == PSW_OV_MASK
}

// Parity
func PSW_P(psw byte) bool {
	return (psw & PSW_P_MASK) == PSW_P_MASK
}

func PSW_SET(psw byte, mask byte) byte {
	return (psw | mask)
}

func PSW_UNSET(psw byte, mask byte) byte {
	return (psw & ^mask)
}

const SFR_ACC uint8 = 0xE0
const SFR_B uint8 = 0xF0
const SFR_DPH uint8 = 0x83
const SFR_DPL uint8 = 0x82
const SFR_IE uint8 = 0xA8
const SFR_IP uint8 = 0xB8
const SFR_P0 uint8 = 0x80
const SFR_P1 uint8 = 0x90
const SFR_P2 uint8 = 0xA0
const SFR_P3 uint8 = 0xB0
const SFR_PCON uint8 = 0x87
const SFR_PSW uint8 = 0xD0
const SFR_SCON uint8 = 0x98
const SFR_SBUF uint8 = 0x99
const SFR_SP uint8 = 0x81
const SFR_TMOD uint8 = 0x89
const SFR_TCON uint8 = 0x88
const SFR_TL0 uint8 = 0x8A
const SFR_TH0 uint8 = 0x8C
const SFR_TL1 uint8 = 0x8B
const SFR_TH1 uint8 = 0x8D

const LOC_R0 uint8 = 0
const LOC_R1 uint8 = 1
const LOC_R2 uint8 = 2
const LOC_R3 uint8 = 3
const LOC_R4 uint8 = 4
const LOC_R5 uint8 = 5
const LOC_R6 uint8 = 6
const LOC_R7 uint8 = 7

// Why uint16: https://stackoverflow.com/questions/57535586/why-does-the-program-counter-in-8051-is-16-bit-and-stack-pointer-is-8-bit-in-805
type Machine struct {
	registers Register
	Program   []byte
	Data      []byte
	PC        uint16 // Program counter / instruction pointer
}

func NewMachine() *Machine {
	vm := Machine{
		registers: Register{},
		Program:   make([]byte, 4*1024, 4*1024), // pre-allocate 4KB ROM
		Data:      make([]byte, 256, 256),       // pre-allocate 256B RAM
		PC:        0,
	}

	return &vm
}

type EvalOperation func(vm *Machine, operands []byte) error

type Opcode struct {
	Name string
	Eval EvalOperation
}

var OPCODES map[byte]Opcode = operationTable()

// []byte{op, oper1, oper2}
func (m *Machine) Feed(instructions []byte) error {
	if len(instructions) < 1 {
		return fmt.Errorf("instruction must be at least 1 byte")
	}

	opcode := instructions[0]
	operands := instructions[1:]

	op, ok := OPCODES[opcode]
	if !ok {
		return fmt.Errorf("opcode '%02x' does not exist in OPCODES", opcode)
	}

	// we know that the instruction will be 3 bytes at most
	// so we can cast the length to uint16
	if m.PC > uint16(0xFFFF)-uint16(len(instructions)) {
		return fmt.Errorf("cannot execute instruction because program counter exceeds 0xFFFF (65535)")
	}

	log.Printf("executing instruction '%02X' (%s) with operand '%v'", opcode, op.Name, operands)

	log.Printf("BEFORE: %+v\n", m.registers)

	evalErr := op.Eval(m, operands)
	if evalErr != nil {
		return fmt.Errorf("VM eval error: %s", evalErr)
	}

	log.Printf("AFTER: %+v\n", m.registers)

	m.PC += uint16(len(instructions))

	return nil
}

// TODO: 8051 has 256B of memory
// but technically it can be extended, so should the location be byte or int?
func (m *Machine) WriteMem(loc uint8, value byte) error {
	if int(loc) >= len(m.Data) {
		return fmt.Errorf("location %#02x exceeds memory capacity of %#02x (%dB)", loc, cap(m.Data), cap(m.Data))
	}

	m.Data[loc] = value

	// these registers are accessible by memory, so we also have to write to it
	// TODO: ada problem kalau write to register directly, nanti value dia tak tally dengan memory
	// TODO: maybe kena validate value for certain registers??
	switch loc {
	case SFR_ACC:
		m.registers.ACC = value
	case SFR_B:
		m.registers.B = value
	case SFR_DPH:
		m.registers.DPH = value
	case SFR_DPL:
		m.registers.DPL = value
	case SFR_IE:
		m.registers.IE = value
	case SFR_IP:
		m.registers.IP = value
	case SFR_P0:
		m.registers.P0 = value
	case SFR_P1:
		m.registers.P1 = value
	case SFR_P2:
		m.registers.P2 = value
	case SFR_P3:
		m.registers.P3 = value
	case SFR_PCON:
		m.registers.PCON = value
	case SFR_PSW:
		m.registers.PSW = value
	case SFR_SCON:
		m.registers.SCON = value
	case SFR_SBUF:
		m.registers.SBUF = value
	case SFR_SP:
		m.registers.SP = value
	case SFR_TMOD:
		m.registers.TMOD = value
	case SFR_TCON:
		m.registers.TCON = value
	case SFR_TL0:
		m.registers.TL0 = value
	case SFR_TH0:
		m.registers.TH0 = value
	case SFR_TL1:
		m.registers.TL1 = value
	case SFR_TH1:
		m.registers.TH1 = value
	}

	return nil
}

func (m *Machine) ReadMem(loc uint8) (byte, error) {
	var value byte = m.Data[loc]

	var registerValue byte

	switch loc {
	case SFR_ACC:
		registerValue = m.registers.ACC
		if registerValue != value {
			fmt.Printf("ACC register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_ACC)
		}
	case SFR_B:
		registerValue = m.registers.B
		if registerValue != value {
			fmt.Printf("B register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_B)
		}
	case SFR_DPH:
		registerValue = m.registers.DPH
		if registerValue != value {
			fmt.Printf("DPH register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_DPH)
		}
	case SFR_DPL:
		registerValue = m.registers.DPL
		if registerValue != value {
			fmt.Printf("DPL register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_DPL)
		}
	case SFR_IE:
		registerValue = m.registers.IE
		if registerValue != value {
			fmt.Printf("IE register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_IE)
		}
	case SFR_IP:
		registerValue = m.registers.IP
		if registerValue != value {
			fmt.Printf("IP register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_IP)
		}
	case SFR_P0:
		registerValue = m.registers.P0
		if registerValue != value {
			fmt.Printf("P0 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_P0)
		}
	case SFR_P1:
		registerValue = m.registers.P1
		if registerValue != value {
			fmt.Printf("P1 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_P1)
		}
	case SFR_P2:
		registerValue = m.registers.P2
		if registerValue != value {
			fmt.Printf("P2 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_P2)
		}
	case SFR_P3:
		registerValue = m.registers.P3
		if registerValue != value {
			fmt.Printf("P3 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_P3)
		}
	case SFR_PCON:
		registerValue = m.registers.PCON
		if registerValue != value {
			fmt.Printf("PCON register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_PCON)
		}
	case SFR_PSW:
		registerValue = m.registers.PSW
		if registerValue != value {
			fmt.Printf("PSW register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_PSW)
		}
	case SFR_SCON:
		registerValue = m.registers.SCON
		if registerValue != value {
			fmt.Printf("SCON register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_SCON)
		}
	case SFR_SBUF:
		registerValue = m.registers.SBUF
		if registerValue != value {
			fmt.Printf("SBUF register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_SBUF)
		}
	case SFR_SP:
		registerValue = m.registers.SP
		if registerValue != value {
			fmt.Printf("SP register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_SP)
		}
	case SFR_TMOD:
		registerValue = m.registers.TMOD
		if registerValue != value {
			fmt.Printf("TMOD register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_TMOD)
		}
	case SFR_TCON:
		registerValue = m.registers.TCON
		if registerValue != value {
			fmt.Printf("TCON register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_TCON)
		}
	case SFR_TL0:
		registerValue = m.registers.TL0
		if registerValue != value {
			fmt.Printf("TL0 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_TL0)
		}
	case SFR_TH0:
		registerValue = m.registers.TH0
		if registerValue != value {
			fmt.Printf("TH0 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_TH0)
		}
	case SFR_TL1:
		registerValue = m.registers.TL1
		if registerValue != value {
			fmt.Printf("TL1 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_TL1)
		}
	case SFR_TH1:
		registerValue = m.registers.TH1
		if registerValue != value {
			fmt.Printf("TH1 register value '%#02x' does not match memory at location %#02x\n", registerValue, SFR_TH1)
		}
	}

	return value, nil
}

func (m *Machine) DerefMem(loc uint8) (byte, error) {
	// TODO: take into account memory banks
	ptr, err := m.ReadMem(loc)
	if err != nil {
		return 0, fmt.Errorf("failed to read memory location %#02x: %s", loc, err)
	}

	val, err := m.ReadMem(ptr)
	if err != nil {
		return 0, fmt.Errorf("dereference error at address %02x: %s", ptr, err)
	}

	return val, nil
}

// Indirect assignment to specified address
// TLDR: int a = 0; int *b = &a; *b = 0xAA;
// Now, a will be 0xAA
func (m *Machine) SetrefMem(loc uint8, value byte) error {
	ptr, err := m.ReadMem(loc)
	if err != nil {
		return err
	}

	err = m.WriteMem(ptr, value)

	return err
}

func main() {
	m := Machine{
		registers: Register{},
	}

	// ni kalau receive raw instruction/byte code
	// TODO: pass assembly, software akan translate jadi byte code, feed masuk VM
	err := m.Feed([]byte{0x24, 0xFF})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
}

func operationTable() map[byte]Opcode {
	tbl := make(map[byte]Opcode)
	tbl[0x00] = Opcode{Name: "NOP", Eval: func(vm *Machine, operands []byte) error {
		fmt.Println("PERFORMING NOP")
		return nil
	}}

	tbl[0x01] = Opcode{Name: "AJMP codeaddr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x02] = Opcode{Name: "LJMP codeaddr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x03] = Opcode{Name: "RR A", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		lsb := A & 0b00000001
		lsbToMsb := (lsb << 7)
		ARightShifted := A >> 1
		retval := lsbToMsb | ARightShifted

		err = vm.WriteMem(SFR_ACC, retval)
		return err
	}}

	tbl[0x04] = Opcode{Name: "INC A", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		A += 1

		err = vm.WriteMem(SFR_ACC, A)
		return err
	}}

	tbl[0x05] = Opcode{Name: "INC ramaddr", Eval: func(vm *Machine, operands []byte) error {
		addr := operands[0]
		val, err := vm.ReadMem(addr)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(addr, val)
		return err
	}}

	tbl[0x06] = Opcode{Name: "INC @R0", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R0 // TODO: take into account memory bank
		val, err := vm.DerefMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.SetrefMem(loc, val)
		return err
	}}

	tbl[0x07] = Opcode{Name: "INC @R1", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R1 // TODO: take into account memory bank
		val, err := vm.DerefMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.SetrefMem(loc, val)
		return err
	}}

	tbl[0x08] = Opcode{Name: "INC R0", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R0 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x09] = Opcode{Name: "INC R1", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R1 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x0a] = Opcode{Name: "INC R2", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R2 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x0b] = Opcode{Name: "INC R3", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R3 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x0c] = Opcode{Name: "INC R4", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R4 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x0d] = Opcode{Name: "INC R5", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R5 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x0e] = Opcode{Name: "INC R6", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R6 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x0f] = Opcode{Name: "INC R7", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R7 // TODO: take into account memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val += 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x10] = Opcode{Name: "JBC bit,rel", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x11] = Opcode{Name: "ACALL page0", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x12] = Opcode{Name: "LCALL codeaddr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x13] = Opcode{Name: "RRC A", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		psw, err := vm.ReadMem(SFR_PSW)
		if err != nil {
			return err
		}

		lsb := A & 0b00000001
		newCarry := lsb

		if newCarry == 0 {
			if err := vm.WriteMem(SFR_PSW, PSW_UNSET(psw, PSW_C_MASK)); err != nil {
				return err
			}
		} else if newCarry == 1 {
			if err := vm.WriteMem(SFR_PSW, PSW_SET(psw, PSW_C_MASK)); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unexpected value %#08b for new carry value", newCarry)
		}

		currentAcc, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		currentAcc >>= 1
		err = vm.WriteMem(SFR_ACC, currentAcc)
		return err
	}}

	tbl[0x14] = Opcode{Name: "DEC A", Eval: func(vm *Machine, operands []byte) error {
		acc, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		acc -= 1
		err = vm.WriteMem(SFR_ACC, acc)
		return err
	}}

	tbl[0x15] = Opcode{Name: "DEC ramaddr", Eval: func(vm *Machine, operands []byte) error {
		addr := operands[0]
		val, err := vm.ReadMem(addr)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(addr, val)
		return err
	}}

	tbl[0x16] = Opcode{Name: "DEC @R0", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R0 // TODO: memory bank
		val, err := vm.DerefMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.SetrefMem(loc, val)
		return err
	}}

	tbl[0x17] = Opcode{Name: "DEC @R1", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R1 // TODO: memory bank
		val, err := vm.DerefMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.SetrefMem(loc, val)
		return err
	}}

	tbl[0x18] = Opcode{Name: "DEC R0", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R0 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x19] = Opcode{Name: "DEC R1", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R1 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x1a] = Opcode{Name: "DEC R2", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R2 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x1b] = Opcode{Name: "DEC R3", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R3 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x1c] = Opcode{Name: "DEC R4", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R4 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x1d] = Opcode{Name: "DEC R5", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R5 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x1e] = Opcode{Name: "DEC R6", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R6 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x1f] = Opcode{Name: "DEC R7", Eval: func(vm *Machine, operands []byte) error {
		loc := LOC_R7 // TODO: memory bank
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}

		val -= 1
		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x20] = Opcode{Name: "JB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x21] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x22] = Opcode{Name: "RET", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x23] = Opcode{Name: "RL A", Eval: func(vm *Machine, operands []byte) error {
		acc, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		msb := acc & 0b10000000
		acc <<= 1
		acc |= (msb >> 7)

		err = vm.WriteMem(SFR_ACC, acc)
		return err
	}}

	// TODO: flags
	tbl[0x24] = Opcode{Name: "ADD A,#data", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		A += operands[0]
		err = vm.WriteMem(SFR_ACC, A)
		return err
	}}

	// TODO: flags
	tbl[0x25] = Opcode{Name: "ADD A,dataaddr", Eval: func(vm *Machine, operands []byte) error {
		addr := operands[0]

		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(addr)
		if err != nil {
			return err
		}

		err = vm.WriteMem(addr, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x26] = Opcode{Name: "ADD A,@R0", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.DerefMem(LOC_R0) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.SetrefMem(LOC_R0, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x27] = Opcode{Name: "ADD A,@R1", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.DerefMem(LOC_R1) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.SetrefMem(LOC_R1, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x28] = Opcode{Name: "ADD A, R0", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R0) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x29] = Opcode{Name: "ADD A, R1", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R1) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x2a] = Opcode{Name: "ADD A, R2", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R2) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err

	}}

	// TODO: flags
	tbl[0x2b] = Opcode{Name: "ADD A, R3", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R3) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x2c] = Opcode{Name: "ADD A, R4", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R4) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x2d] = Opcode{Name: "ADD A, R5", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R5) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x2e] = Opcode{Name: "ADD A, R6", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R6) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	// TODO: flags
	tbl[0x2f] = Opcode{Name: "ADD A, R7", Eval: func(vm *Machine, operands []byte) error {
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		val, err := vm.ReadMem(LOC_R7) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(SFR_ACC, A+val)
		return err
	}}

	tbl[0x30] = Opcode{Name: "JNB bit addr,code addr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x31] = Opcode{Name: "ACALL code addr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x32] = Opcode{Name: "RETI", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x33] = Opcode{Name: "RLC A", Eval: func(vm *Machine, operands []byte) error {
		acc, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}

		psw, err := vm.ReadMem(SFR_PSW)
		if err != nil {
			return err
		}

		msb := acc & 0b10000000
		if msb == 0b10000000 {
			err := vm.WriteMem(SFR_PSW, PSW_SET(psw, PSW_C_MASK))
			if err != nil {
				return err
			}
		} else {
			err := vm.WriteMem(SFR_PSW, PSW_UNSET(psw, PSW_C_MASK))
			if err != nil {
				return err
			}
		}
		acc <<= 1

		err = vm.WriteMem(SFR_ACC, acc)
		return err
	}}

	tbl[0x34] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x35] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x36] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x37] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x38] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x39] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x3a] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x3b] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x3c] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x3d] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x3e] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x3f] = Opcode{Name: "ADDC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x40] = Opcode{Name: "JC reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x41] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x42] = Opcode{Name: "ORL data addr,A", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x43] = Opcode{Name: "ORL data addr,#data", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x44] = Opcode{Name: "ORL A,#data", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x45] = Opcode{Name: "ORL A,data add", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x46] = Opcode{Name: "ORL A,@R0", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x47] = Opcode{Name: "ORL A,@R1", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x48] = Opcode{Name: "ORL A,R0", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x49] = Opcode{Name: "ORL A,R1", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4a] = Opcode{Name: "ORL A,R2", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4b] = Opcode{Name: "ORL A,R3", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4c] = Opcode{Name: "ORL A,R4", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4d] = Opcode{Name: "ORL A,R5", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4e] = Opcode{Name: "ORL A,R6", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4f] = Opcode{Name: "ORL A,R7", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x50] = Opcode{Name: "JNC reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x51] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x52] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x53] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x54] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x55] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x56] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x57] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x58] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x59] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x5a] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x5b] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x5c] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x5d] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x5e] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x5f] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x60] = Opcode{Name: "JZ reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x61] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x62] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x63] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x64] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x65] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x66] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x67] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x68] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x69] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x6a] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x6b] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x6c] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x6d] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x6e] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x6f] = Opcode{Name: "XRL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x70] = Opcode{Name: "JNZ reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x71] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x72] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x73] = Opcode{Name: "JMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x74] = Opcode{Name: "MOV A,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		A, err := vm.ReadMem(SFR_ACC)
		if err != nil {
			return err
		}
		err = vm.WriteMem(SFR_ACC, A+data)
		return err
	}}

	tbl[0x75] = Opcode{Name: "MOV ramaddr,#data", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		data := operands[1]
		err := vm.WriteMem(loc, data)
		return err
	}}

	tbl[0x76] = Opcode{Name: "MOV @R0,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.SetrefMem(LOC_R0, data) // TODO: memory bank
		return err
	}}

	tbl[0x77] = Opcode{Name: "MOV @R1,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.SetrefMem(LOC_R1, data) // TODO: memory bank
		return err
	}}

	tbl[0x78] = Opcode{Name: "MOV R0,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R0, data) // TODO: memory bank
		return err
	}}

	tbl[0x79] = Opcode{Name: "MOV R1,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R1, data) // TODO: memory bank
		return err
	}}

	tbl[0x7a] = Opcode{Name: "MOV R2,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R2, data) // TODO: memory bank
		return err
	}}

	tbl[0x7b] = Opcode{Name: "MOV R3,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R3, data) // TODO: memory bank
		return err
	}}

	tbl[0x7c] = Opcode{Name: "MOV R4,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R4, data) // TODO: memory bank
		return err
	}}

	tbl[0x7d] = Opcode{Name: "MOV R5,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R5, data) // TODO: memory bank
		return err
	}}

	tbl[0x7e] = Opcode{Name: "MOV R6,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R6, data) // TODO: memory bank
		return err
	}}

	tbl[0x7f] = Opcode{Name: "MOV R7,#data", Eval: func(vm *Machine, operands []byte) error {
		data := operands[0]
		err := vm.WriteMem(LOC_R7, data) // TODO: memory bank
		return err
	}}

	tbl[0x80] = Opcode{Name: "SJMP reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x81] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x82] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x83] = Opcode{Name: "MOVC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x84] = Opcode{Name: "DIV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x85] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x86] = Opcode{Name: "MOV ramaddr,@R0", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.DerefMem(LOC_R0) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x87] = Opcode{Name: "MOV ramaddr,@R1", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.DerefMem(LOC_R1) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x88] = Opcode{Name: "MOV ramaddr,R0", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R0) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x89] = Opcode{Name: "MOV ramaddr,R1", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R1) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x8a] = Opcode{Name: "MOV ramaddr,R2", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R2) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x8b] = Opcode{Name: "MOV ramaddr,R3", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R3) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x8c] = Opcode{Name: "MOV ramaddr,R4", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R4) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x8d] = Opcode{Name: "MOV ramaddr,R5", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R5) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x8e] = Opcode{Name: "MOV ramaddr,R6", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R6) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x8f] = Opcode{Name: "MOV ramaddr,R7", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]

		val, err := vm.ReadMem(LOC_R7) // TODO: memory bank
		if err != nil {
			return err
		}

		err = vm.WriteMem(loc, val)
		return err
	}}

	tbl[0x90] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x91] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x92] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x93] = Opcode{Name: "MOVC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x94] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x95] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x96] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x97] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x98] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x99] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x9a] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x9b] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x9c] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x9d] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x9e] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x9f] = Opcode{Name: "SUBB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa0] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa1] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa2] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa3] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa4] = Opcode{Name: "MUL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa5] = Opcode{Name: "?", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa6] = Opcode{Name: "MOV @R0,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.SetrefMem(LOC_R0, val) // TODO: memory bank
		return err
	}}

	tbl[0xa7] = Opcode{Name: "MOV @R1,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.SetrefMem(LOC_R1, val) // TODO: memory bank
		return err
	}}

	tbl[0xa8] = Opcode{Name: "MOV R0,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R0, val) // TODO: memory bank
		return err
	}}

	tbl[0xa9] = Opcode{Name: "MOV R1,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R1, val) // TODO: memory bank
		return err
	}}

	tbl[0xaa] = Opcode{Name: "MOV R2,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R2, val) // TODO: memory bank
		return err
	}}

	tbl[0xab] = Opcode{Name: "MOV R3,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R3, val) // TODO: memory bank
		return err
	}}

	tbl[0xac] = Opcode{Name: "MOV R4,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R4, val) // TODO: memory bank
		return err
	}}

	tbl[0xad] = Opcode{Name: "MOV R5,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R5, val) // TODO: memory bank
		return err
	}}

	tbl[0xae] = Opcode{Name: "MOV R6,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R6, val) // TODO: memory bank
		return err
	}}

	tbl[0xaf] = Opcode{Name: "MOV R7,ramaddr", Eval: func(vm *Machine, operands []byte) error {
		loc := operands[0]
		val, err := vm.ReadMem(loc)
		if err != nil {
			return err
		}
		err = vm.WriteMem(LOC_R7, val) // TODO: memory bank
		return err
	}}

	tbl[0xb0] = Opcode{Name: "ANL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb1] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb2] = Opcode{Name: "CPL bitaddr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb3] = Opcode{Name: "CPL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb4] = Opcode{Name: "CJNE A,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb5] = Opcode{Name: "CJNE A,iram addr,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb6] = Opcode{Name: "CJNE @R0,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb7] = Opcode{Name: "CJNE @R1,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb8] = Opcode{Name: "CJNE R0,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xb9] = Opcode{Name: "CJNE R1,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xba] = Opcode{Name: "CJNE R2,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xbb] = Opcode{Name: "CJNE R3,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xbc] = Opcode{Name: "CJNE R4,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xbd] = Opcode{Name: "CJNE R5,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xbe] = Opcode{Name: "CJNE R6,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xbf] = Opcode{Name: "CJNE R7,#data,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc0] = Opcode{Name: "PUSH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc1] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc2] = Opcode{Name: "CLR", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc3] = Opcode{Name: "CLR", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc4] = Opcode{Name: "SWAP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc5] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc6] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc7] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc8] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xc9] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xca] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xcb] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xcc] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xcd] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xce] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xcf] = Opcode{Name: "XCH", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd0] = Opcode{Name: "POP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd1] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd2] = Opcode{Name: "SETB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd3] = Opcode{Name: "SETB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd4] = Opcode{Name: "DA", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd5] = Opcode{Name: "DJNZ", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd6] = Opcode{Name: "XCHD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd7] = Opcode{Name: "XCHD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd8] = Opcode{Name: "DJNZ R0,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xd9] = Opcode{Name: "DJNZ R1,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xda] = Opcode{Name: "DJNZ R2,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xdb] = Opcode{Name: "DJNZ R3,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xdc] = Opcode{Name: "DJNZ R4,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xdd] = Opcode{Name: "DJNZ R5,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xde] = Opcode{Name: "DJNZ R6,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xdf] = Opcode{Name: "DJNZ R7,reladdr", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe0] = Opcode{Name: "MOVX", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe1] = Opcode{Name: "AJMP", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe2] = Opcode{Name: "MOVX", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe3] = Opcode{Name: "MOVX", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe4] = Opcode{Name: "CLR", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe5] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe6] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe7] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe8] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xe9] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xea] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xeb] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xec] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xed] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xee] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xef] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf0] = Opcode{Name: "MOVX", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf1] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf2] = Opcode{Name: "MOVX", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf3] = Opcode{Name: "MOVX", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf4] = Opcode{Name: "CPL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf5] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf6] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf7] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf8] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xf9] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xfa] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xfb] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xfc] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xfd] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xfe] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xff] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	return tbl
}
