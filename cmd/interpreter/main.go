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

	tbl[0x03] = Opcode{Name: "RR", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x04] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x05] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x06] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x07] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x08] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x09] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x0a] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x0b] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x0c] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x0d] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x0e] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x0f] = Opcode{Name: "INC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x10] = Opcode{Name: "JBC", Eval: func(vm *Machine, operands []byte) error {
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

	tbl[0x13] = Opcode{Name: "RRC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x14] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x15] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x16] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x17] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x18] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x19] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x1a] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x1b] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x1c] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x1d] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x1e] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x1f] = Opcode{Name: "DEC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
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

	tbl[0x23] = Opcode{Name: "RL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x24] = Opcode{Name: "ADD A,#data", Eval: func(vm *Machine, operands []byte) error {
		log.Println("performing ADD A,#data")
		A := vm.Registers.ACC
		A += operands[0]
		vm.Registers.ACC = A
		return nil
	}}

	tbl[0x25] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x26] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x27] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x28] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x29] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x2a] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x2b] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x2c] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x2d] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x2e] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x2f] = Opcode{Name: "ADD", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x30] = Opcode{Name: "JNB", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x31] = Opcode{Name: "ACALL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x32] = Opcode{Name: "RETI", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x33] = Opcode{Name: "RLC", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
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

	tbl[0x42] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x43] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x44] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x45] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x46] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x47] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x48] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x49] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4a] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4b] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4c] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4d] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4e] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x4f] = Opcode{Name: "ORL", Eval: func(vm *Machine, operands []byte) error {
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

	tbl[0x74] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x75] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x76] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x77] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x78] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x79] = Opcode{Name: "MOV R1,#data", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x7a] = Opcode{Name: "MOV R2,#data", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x7b] = Opcode{Name: "MOV R3,#data", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x7c] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x7d] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x7e] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x7f] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
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

	tbl[0x86] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x87] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x88] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x89] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x8a] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x8b] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x8c] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x8d] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x8e] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0x8f] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
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

	tbl[0xa6] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa7] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa8] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xa9] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xaa] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xab] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xac] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xad] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xae] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
	}}

	tbl[0xaf] = Opcode{Name: "MOV", Eval: func(vm *Machine, operands []byte) error {
		// TODO: implement
		return nil
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
