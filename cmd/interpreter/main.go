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

type Machine struct {
	Registers Register
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
		return fmt.Errorf("opcode '%02x' does not exist in OPCODES")
	}

	log.Printf("executing instruction '%02X' (%s) with operand '%v'", opcode, op.Name, operands)

	log.Printf("BEFORE: %+v\n", m.Registers)

	evalErr := op.Eval(m, operands)
	if evalErr != nil {
		return fmt.Errorf("VM eval error: %s", evalErr)
	}

	log.Printf("AFTER: %+v\n", m.Registers)

	return nil
}

func main() {
	m := Machine{
		Registers: Register{},
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
	tbl[0x00+0x00] = Opcode{Name: "NOP", Eval: func(vm *Machine, operands []byte) error {
		fmt.Println("PERFORMING NOP")
		return nil
	}}

	tbl[0x24] = Opcode{Name: "ADD A,#data", Eval: func(vm *Machine, operands []byte) error {
		log.Println("performing ADD A,#data")
		A := vm.Registers.ACC
		A += operands[0]
		vm.Registers.ACC = A
		return nil
	}}

	return tbl
}
