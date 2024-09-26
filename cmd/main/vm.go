package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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

// @param program - slice of the program starting from the
// position of the current instruction onwards
//
// @param pos - the position of the current instruction
//
// @return sz - size of the instruction (opcode + operands)
// so that the next loop iteration will know
// how many elements to skip
type Operation func(program []byte, pos int) int

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: ./vm <binary>")
		fmt.Println("usage: ./vm examples/blink.bin")
		return
	}

	OPCODES_NAME := make(map[byte]string)
	OPCODES_NAME[0x00+0x00] = "NOP"
	OPCODES_NAME[0x00+0x01] = "AJMP codeaddr"
	OPCODES_NAME[0x00+0x02] = "LJMP codeaddr"
	OPCODES_NAME[0x00+0x03] = "RR"
	OPCODES_NAME[0x00+0x04] = "INC"
	OPCODES_NAME[0x00+0x05] = "INC"
	OPCODES_NAME[0x00+0x06] = "INC"
	OPCODES_NAME[0x00+0x07] = "INC"
	OPCODES_NAME[0x00+0x08] = "INC"
	OPCODES_NAME[0x00+0x09] = "INC"
	OPCODES_NAME[0x00+0x0a] = "INC"
	OPCODES_NAME[0x00+0x0b] = "INC"
	OPCODES_NAME[0x00+0x0c] = "INC"
	OPCODES_NAME[0x00+0x0d] = "INC"
	OPCODES_NAME[0x00+0x0e] = "INC"
	OPCODES_NAME[0x00+0x0f] = "INC"

	OPCODES_NAME[0x10+0x00] = "JBC"
	OPCODES_NAME[0x10+0x01] = "ACALL page0"
	OPCODES_NAME[0x10+0x02] = "LCALL codeaddr"
	OPCODES_NAME[0x10+0x03] = "RRC"
	OPCODES_NAME[0x10+0x04] = "DEC"
	OPCODES_NAME[0x10+0x05] = "DEC"
	OPCODES_NAME[0x10+0x06] = "DEC"
	OPCODES_NAME[0x10+0x07] = "DEC"
	OPCODES_NAME[0x10+0x08] = "DEC"
	OPCODES_NAME[0x10+0x09] = "DEC"
	OPCODES_NAME[0x10+0x0a] = "DEC"
	OPCODES_NAME[0x10+0x0b] = "DEC"
	OPCODES_NAME[0x10+0x0c] = "DEC"
	OPCODES_NAME[0x10+0x0d] = "DEC"
	OPCODES_NAME[0x10+0x0e] = "DEC"
	OPCODES_NAME[0x10+0x0f] = "DEC"

	OPCODES_NAME[0x20+0x00] = "JB"
	OPCODES_NAME[0x20+0x01] = "AJMP"
	OPCODES_NAME[0x20+0x02] = "RET"
	OPCODES_NAME[0x20+0x03] = "RL"
	OPCODES_NAME[0x20+0x04] = "ADD"
	OPCODES_NAME[0x20+0x05] = "ADD"
	OPCODES_NAME[0x20+0x06] = "ADD"
	OPCODES_NAME[0x20+0x07] = "ADD"
	OPCODES_NAME[0x20+0x08] = "ADD"
	OPCODES_NAME[0x20+0x09] = "ADD"
	OPCODES_NAME[0x20+0x0a] = "ADD"
	OPCODES_NAME[0x20+0x0b] = "ADD"
	OPCODES_NAME[0x20+0x0c] = "ADD"
	OPCODES_NAME[0x20+0x0d] = "ADD"
	OPCODES_NAME[0x20+0x0e] = "ADD"
	OPCODES_NAME[0x20+0x0f] = "ADD"

	OPCODES_NAME[0x30+0x00] = "JNB"
	OPCODES_NAME[0x30+0x01] = "ACALL"
	OPCODES_NAME[0x30+0x02] = "RETI"
	OPCODES_NAME[0x30+0x03] = "RLC"
	OPCODES_NAME[0x30+0x04] = "ADDC"
	OPCODES_NAME[0x30+0x05] = "ADDC"
	OPCODES_NAME[0x30+0x06] = "ADDC"
	OPCODES_NAME[0x30+0x07] = "ADDC"
	OPCODES_NAME[0x30+0x08] = "ADDC"
	OPCODES_NAME[0x30+0x09] = "ADDC"
	OPCODES_NAME[0x30+0x0a] = "ADDC"
	OPCODES_NAME[0x30+0x0b] = "ADDC"
	OPCODES_NAME[0x30+0x0c] = "ADDC"
	OPCODES_NAME[0x30+0x0d] = "ADDC"
	OPCODES_NAME[0x30+0x0e] = "ADDC"
	OPCODES_NAME[0x30+0x0f] = "ADDC"

	OPCODES_NAME[0x40+0x00] = "JC reladdr"
	OPCODES_NAME[0x40+0x01] = "AJMP"
	OPCODES_NAME[0x40+0x02] = "ORL"
	OPCODES_NAME[0x40+0x03] = "ORL"
	OPCODES_NAME[0x40+0x04] = "ORL"
	OPCODES_NAME[0x40+0x05] = "ORL"
	OPCODES_NAME[0x40+0x06] = "ORL"
	OPCODES_NAME[0x40+0x07] = "ORL"
	OPCODES_NAME[0x40+0x08] = "ORL"
	OPCODES_NAME[0x40+0x09] = "ORL"
	OPCODES_NAME[0x40+0x0a] = "ORL"
	OPCODES_NAME[0x40+0x0b] = "ORL"
	OPCODES_NAME[0x40+0x0c] = "ORL"
	OPCODES_NAME[0x40+0x0d] = "ORL"
	OPCODES_NAME[0x40+0x0e] = "ORL"
	OPCODES_NAME[0x40+0x0f] = "ORL"

	OPCODES_NAME[0x50+0x00] = "JNC reladdr"
	OPCODES_NAME[0x50+0x01] = "ACALL"
	OPCODES_NAME[0x50+0x02] = "ANL"
	OPCODES_NAME[0x50+0x03] = "ANL"
	OPCODES_NAME[0x50+0x04] = "ANL"
	OPCODES_NAME[0x50+0x05] = "ANL"
	OPCODES_NAME[0x50+0x06] = "ANL"
	OPCODES_NAME[0x50+0x07] = "ANL"
	OPCODES_NAME[0x50+0x08] = "ANL"
	OPCODES_NAME[0x50+0x09] = "ANL"
	OPCODES_NAME[0x50+0x0a] = "ANL"
	OPCODES_NAME[0x50+0x0b] = "ANL"
	OPCODES_NAME[0x50+0x0c] = "ANL"
	OPCODES_NAME[0x50+0x0d] = "ANL"
	OPCODES_NAME[0x50+0x0e] = "ANL"
	OPCODES_NAME[0x50+0x0f] = "ANL"

	OPCODES_NAME[0x60+0x00] = "JZ reladdr"
	OPCODES_NAME[0x60+0x01] = "AJMP"
	OPCODES_NAME[0x60+0x02] = "XRL"
	OPCODES_NAME[0x60+0x03] = "XRL"
	OPCODES_NAME[0x60+0x04] = "XRL"
	OPCODES_NAME[0x60+0x05] = "XRL"
	OPCODES_NAME[0x60+0x06] = "XRL"
	OPCODES_NAME[0x60+0x07] = "XRL"
	OPCODES_NAME[0x60+0x08] = "XRL"
	OPCODES_NAME[0x60+0x09] = "XRL"
	OPCODES_NAME[0x60+0x0a] = "XRL"
	OPCODES_NAME[0x60+0x0b] = "XRL"
	OPCODES_NAME[0x60+0x0c] = "XRL"
	OPCODES_NAME[0x60+0x0d] = "XRL"
	OPCODES_NAME[0x60+0x0e] = "XRL"
	OPCODES_NAME[0x60+0x0f] = "XRL"

	OPCODES_NAME[0x70+0x00] = "JNZ reladdr"
	OPCODES_NAME[0x70+0x01] = "ACALL"
	OPCODES_NAME[0x70+0x02] = "ORL"
	OPCODES_NAME[0x70+0x03] = "JMP"
	OPCODES_NAME[0x70+0x04] = "MOV"
	OPCODES_NAME[0x70+0x05] = "MOV"
	OPCODES_NAME[0x70+0x06] = "MOV"
	OPCODES_NAME[0x70+0x07] = "MOV"
	OPCODES_NAME[0x70+0x08] = "MOV"
	OPCODES_NAME[0x70+0x09] = "MOV R1,#data"
	OPCODES_NAME[0x70+0x0a] = "MOV R2,#data"
	OPCODES_NAME[0x70+0x0b] = "MOV R3,#data"
	OPCODES_NAME[0x70+0x0c] = "MOV"
	OPCODES_NAME[0x70+0x0d] = "MOV"
	OPCODES_NAME[0x70+0x0e] = "MOV"
	OPCODES_NAME[0x70+0x0f] = "MOV"

	OPCODES_NAME[0x80+0x00] = "SJMP reladdr"
	OPCODES_NAME[0x80+0x01] = "AJMP"
	OPCODES_NAME[0x80+0x02] = "ANL"
	OPCODES_NAME[0x80+0x03] = "MOVC"
	OPCODES_NAME[0x80+0x04] = "DIV"
	OPCODES_NAME[0x80+0x05] = "MOV"
	OPCODES_NAME[0x80+0x06] = "MOV"
	OPCODES_NAME[0x80+0x07] = "MOV"
	OPCODES_NAME[0x80+0x08] = "MOV"
	OPCODES_NAME[0x80+0x09] = "MOV"
	OPCODES_NAME[0x80+0x0a] = "MOV"
	OPCODES_NAME[0x80+0x0b] = "MOV"
	OPCODES_NAME[0x80+0x0c] = "MOV"
	OPCODES_NAME[0x80+0x0d] = "MOV"
	OPCODES_NAME[0x80+0x0e] = "MOV"
	OPCODES_NAME[0x80+0x0f] = "MOV"

	OPCODES_NAME[0x90+0x01] = "ACALL"
	OPCODES_NAME[0x90+0x02] = "MOV"
	OPCODES_NAME[0x90+0x00] = "MOV"
	OPCODES_NAME[0x90+0x03] = "MOVC"
	OPCODES_NAME[0x90+0x04] = "SUBB"
	OPCODES_NAME[0x90+0x05] = "SUBB"
	OPCODES_NAME[0x90+0x06] = "SUBB"
	OPCODES_NAME[0x90+0x07] = "SUBB"
	OPCODES_NAME[0x90+0x08] = "SUBB"
	OPCODES_NAME[0x90+0x09] = "SUBB"
	OPCODES_NAME[0x90+0x0a] = "SUBB"
	OPCODES_NAME[0x90+0x0b] = "SUBB"
	OPCODES_NAME[0x90+0x0c] = "SUBB"
	OPCODES_NAME[0x90+0x0d] = "SUBB"
	OPCODES_NAME[0x90+0x0e] = "SUBB"
	OPCODES_NAME[0x90+0x0f] = "SUBB"

	OPCODES_NAME[0xa0+0x00] = "ORL"
	OPCODES_NAME[0xa0+0x01] = "AJMP"
	OPCODES_NAME[0xa0+0x02] = "MOV"
	OPCODES_NAME[0xa0+0x03] = "INC"
	OPCODES_NAME[0xa0+0x04] = "MUL"
	OPCODES_NAME[0xa0+0x05] = "?"
	OPCODES_NAME[0xa0+0x06] = "MOV"
	OPCODES_NAME[0xa0+0x07] = "MOV"
	OPCODES_NAME[0xa0+0x08] = "MOV"
	OPCODES_NAME[0xa0+0x09] = "MOV"
	OPCODES_NAME[0xa0+0x0a] = "MOV"
	OPCODES_NAME[0xa0+0x0b] = "MOV"
	OPCODES_NAME[0xa0+0x0c] = "MOV"
	OPCODES_NAME[0xa0+0x0d] = "MOV"
	OPCODES_NAME[0xa0+0x0e] = "MOV"
	OPCODES_NAME[0xa0+0x0f] = "MOV"

	OPCODES_NAME[0xb0+0x00] = "ANL"
	OPCODES_NAME[0xb0+0x01] = "ACALL"
	OPCODES_NAME[0xb0+0x02] = "CPL bitaddr"
	OPCODES_NAME[0xb0+0x03] = "CPL"
	OPCODES_NAME[0xb0+0x04] = "CJNE A,#data,reladdr"
	OPCODES_NAME[0xb0+0x05] = "CJNE A,iram addr,reladdr"
	OPCODES_NAME[0xb0+0x06] = "CJNE @R0,#data,reladdr"
	OPCODES_NAME[0xb0+0x07] = "CJNE @R1,#data,reladdr"
	OPCODES_NAME[0xb0+0x08] = "CJNE R0,#data,reladdr"
	OPCODES_NAME[0xb0+0x09] = "CJNE R1,#data,reladdr"
	OPCODES_NAME[0xb0+0x0a] = "CJNE R2,#data,reladdr"
	OPCODES_NAME[0xb0+0x0b] = "CJNE R3,#data,reladdr"
	OPCODES_NAME[0xb0+0x0c] = "CJNE R4,#data,reladdr"
	OPCODES_NAME[0xb0+0x0d] = "CJNE R5,#data,reladdr"
	OPCODES_NAME[0xb0+0x0e] = "CJNE R6,#data,reladdr"
	OPCODES_NAME[0xb0+0x0f] = "CJNE R7,#data,reladdr"

	OPCODES_NAME[0xc0+0x00] = "PUSH"
	OPCODES_NAME[0xc0+0x01] = "AJMP"
	OPCODES_NAME[0xc0+0x02] = "CLR"
	OPCODES_NAME[0xc0+0x03] = "CLR"
	OPCODES_NAME[0xc0+0x04] = "SWAP"
	OPCODES_NAME[0xc0+0x05] = "XCH"
	OPCODES_NAME[0xc0+0x06] = "XCH"
	OPCODES_NAME[0xc0+0x07] = "XCH"
	OPCODES_NAME[0xc0+0x08] = "XCH"
	OPCODES_NAME[0xc0+0x09] = "XCH"
	OPCODES_NAME[0xc0+0x0a] = "XCH"
	OPCODES_NAME[0xc0+0x0b] = "XCH"
	OPCODES_NAME[0xc0+0x0c] = "XCH"
	OPCODES_NAME[0xc0+0x0d] = "XCH"
	OPCODES_NAME[0xc0+0x0e] = "XCH"
	OPCODES_NAME[0xc0+0x0f] = "XCH"

	OPCODES_NAME[0xd0+0x00] = "POP"
	OPCODES_NAME[0xd0+0x01] = "ACALL"
	OPCODES_NAME[0xd0+0x02] = "SETB"
	OPCODES_NAME[0xd0+0x03] = "SETB"
	OPCODES_NAME[0xd0+0x04] = "DA"
	OPCODES_NAME[0xd0+0x05] = "DJNZ"
	OPCODES_NAME[0xd0+0x06] = "XCHD"
	OPCODES_NAME[0xd0+0x07] = "XCHD"
	OPCODES_NAME[0xd0+0x08] = "DJNZ R0,reladdr"
	OPCODES_NAME[0xd0+0x09] = "DJNZ R1,reladdr"
	OPCODES_NAME[0xd0+0x0a] = "DJNZ R2,reladdr"
	OPCODES_NAME[0xd0+0x0b] = "DJNZ R3,reladdr"
	OPCODES_NAME[0xd0+0x0c] = "DJNZ R4,reladdr"
	OPCODES_NAME[0xd0+0x0d] = "DJNZ R5,reladdr"
	OPCODES_NAME[0xd0+0x0e] = "DJNZ R6,reladdr"
	OPCODES_NAME[0xd0+0x0f] = "DJNZ R7,reladdr"

	OPCODES_NAME[0xe0+0x00] = "MOVX"
	OPCODES_NAME[0xe0+0x01] = "AJMP"
	OPCODES_NAME[0xe0+0x02] = "MOVX"
	OPCODES_NAME[0xe0+0x03] = "MOVX"
	OPCODES_NAME[0xe0+0x04] = "CLR"
	OPCODES_NAME[0xe0+0x05] = "MOV"
	OPCODES_NAME[0xe0+0x06] = "MOV"
	OPCODES_NAME[0xe0+0x07] = "MOV"
	OPCODES_NAME[0xe0+0x08] = "MOV"
	OPCODES_NAME[0xe0+0x09] = "MOV"
	OPCODES_NAME[0xe0+0x0a] = "MOV"
	OPCODES_NAME[0xe0+0x0b] = "MOV"
	OPCODES_NAME[0xe0+0x0c] = "MOV"
	OPCODES_NAME[0xe0+0x0d] = "MOV"
	OPCODES_NAME[0xe0+0x0e] = "MOV"
	OPCODES_NAME[0xe0+0x0f] = "MOV"

	OPCODES_NAME[0xf0+0x00] = "MOVX"
	OPCODES_NAME[0xf0+0x01] = "ACALL"
	OPCODES_NAME[0xf0+0x02] = "MOVX"
	OPCODES_NAME[0xf0+0x03] = "MOVX"
	OPCODES_NAME[0xf0+0x04] = "CPL"
	OPCODES_NAME[0xf0+0x05] = "MOV"
	OPCODES_NAME[0xf0+0x06] = "MOV"
	OPCODES_NAME[0xf0+0x07] = "MOV"
	OPCODES_NAME[0xf0+0x08] = "MOV"
	OPCODES_NAME[0xf0+0x09] = "MOV"
	OPCODES_NAME[0xf0+0x0a] = "MOV"
	OPCODES_NAME[0xf0+0x0b] = "MOV"
	OPCODES_NAME[0xf0+0x0c] = "MOV"
	OPCODES_NAME[0xf0+0x0d] = "MOV"
	OPCODES_NAME[0xf0+0x0e] = "MOV"
	OPCODES_NAME[0xf0+0x0f] = "MOV"

	OPCODES := make(map[byte]Operation)

	// NOP
	OPCODES[0x00+0x00] = func(program []byte, pos int) int {
		return 1
	}

	// AJMP codeaddr
	OPCODES[0x00+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// LJMP
	OPCODES[0x00+0x02] = func(program []byte, pos int) int {
		operand := program[1:3]
		fmt.Printf("operand: %02x\n", operand)
		return 3
	}

	// RR
	OPCODES[0x00+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// INC
	OPCODES[0x00+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// INC
	OPCODES[0x00+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JBC
	OPCODES[0x10+0x00] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// ACALL
	OPCODES[0x10+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// LCALL codeaddr
	OPCODES[0x10+0x02] = func(program []byte, pos int) int {
		operand := program[1:3]
		fmt.Printf("operand: %02x\n", operand)
		return 3
	}

	// RRC
	OPCODES[0x10+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DEC
	OPCODES[0x10+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// DEC
	OPCODES[0x10+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JB
	OPCODES[0x20+0x00] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// AJMP
	OPCODES[0x20+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// RET
	OPCODES[0x20+0x02] = func(program []byte, pos int) int {
		return 1
	}

	// RL
	OPCODES[0x20+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ADD
	OPCODES[0x20+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ADD
	OPCODES[0x20+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// ADD
	OPCODES[0x20+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JNB
	OPCODES[0x30+0x00] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// ACALL
	OPCODES[0x30+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// RETI
	OPCODES[0x30+0x02] = func(program []byte, pos int) int {
		return 1
	}

	// RLC
	OPCODES[0x30+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ADDC
	OPCODES[0x30+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ADDC
	OPCODES[0x30+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// ADDC
	OPCODES[0x30+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JC reladdr
	OPCODES[0x40+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// AJMP
	OPCODES[0x40+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ORL
	OPCODES[0x40+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ORL
	OPCODES[0x40+0x03] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// ORL
	OPCODES[0x40+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ORL
	OPCODES[0x40+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ORL
	OPCODES[0x40+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0x40+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JNC reladdr
	OPCODES[0x50+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ACALL
	OPCODES[0x50+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ANL
	OPCODES[0x50+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ANL
	OPCODES[0x50+0x03] = func(program []byte, pos int) int {
		return 3
	}

	// ANL
	OPCODES[0x50+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ANL
	OPCODES[0x50+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ANL
	OPCODES[0x50+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// ANL
	OPCODES[0x50+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JZ
	OPCODES[0x60+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// AJMP
	OPCODES[0x60+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// XRL
	OPCODES[0x60+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// XRL
	OPCODES[0x60+0x03] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// XRL
	OPCODES[0x60+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// XRL
	OPCODES[0x60+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// XRL
	OPCODES[0x60+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// XRL
	OPCODES[0x60+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// JNZ reladdr
	OPCODES[0x70+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ACALL
	OPCODES[0x70+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ORL
	OPCODES[0x70+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// JMP
	OPCODES[0x70+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0x70+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x05] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// MOV
	OPCODES[0x70+0x06] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x07] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x08] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x09] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x0A] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x0B] = func(program []byte, pos int) int {
		// program[0] is the instruction
		// therefore, the index 1 is the operand
		operand := program[1]
		fmt.Printf("operand:  %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x0C] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x0D] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x0E] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x70+0x0F] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// SJMP reladdr
	OPCODES[0x80+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// AJMP
	OPCODES[0x80+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ANL
	OPCODES[0x80+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOVC
	OPCODES[0x80+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// DIV
	OPCODES[0x80+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0x80+0x05] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// MOV
	OPCODES[0x80+0x06] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x07] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x08] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x09] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x0A] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x0B] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x0C] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x0D] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x0E] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x80+0x0F] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x90+0x00] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// ACALL
	OPCODES[0x90+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0x90+0x02] = func(program []byte, pos int) int {
		return 2
	}

	// MOVC
	OPCODES[0x90+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x04] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// SUBB
	OPCODES[0x90+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// SUBB
	OPCODES[0x90+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// SUBB
	OPCODES[0x90+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// ORL
	OPCODES[0xA0+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// AJMP
	OPCODES[0xA0+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// INC
	OPCODES[0xA0+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// MUL
	OPCODES[0xA0+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// ?
	OPCODES[0xA0+0x05] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xA0+0x06] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x07] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x08] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x09] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x0A] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x0B] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x0C] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x0D] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x0E] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xA0+0x0F] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ANL
	OPCODES[0xB0+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ACALL
	OPCODES[0xB0+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// CPL bitaddr
	OPCODES[0xB0+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// CPL
	OPCODES[0xB0+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// CJNE A,#data,reladdr
	OPCODES[0xB0+0x04] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE A,iram addr,reladdr
	OPCODES[0xB0+0x05] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x06] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x07] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE R0,#data,reladdr
	OPCODES[0xB0+0x08] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x09] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x0A] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x0B] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x0C] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x0D] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x0E] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// CJNE
	OPCODES[0xB0+0x0F] = func(program []byte, pos int) int {
		// TODO: operand
		return 3
	}

	// PUSH
	OPCODES[0xC0+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// AJMP
	OPCODES[0xC0+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// CLR
	OPCODES[0xC0+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// CLR
	OPCODES[0xC0+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// SWAP
	OPCODES[0xC0+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// XCH
	OPCODES[0xC0+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// XCH
	OPCODES[0xC0+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// POP
	OPCODES[0xD0+0x00] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// ACALL
	OPCODES[0xD0+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// SETB
	OPCODES[0xD0+0x02] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// SETB
	OPCODES[0xD0+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// DA
	OPCODES[0xD0+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// DJNZ
	OPCODES[0xD0+0x05] = func(program []byte, pos int) int {
		// TODO: oper
		return 3
	}

	// XCHD
	OPCODES[0xD0+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// XCHD
	OPCODES[0xD0+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// DJNZ R0,reladdr
	OPCODES[0xD0+0x08] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R1,reladdr
	OPCODES[0xD0+0x09] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R2,reladdr
	OPCODES[0xD0+0x0A] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R3,reladdr
	OPCODES[0xD0+0x0B] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R4,reladdr
	OPCODES[0xD0+0x0C] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R5,reladdr
	OPCODES[0xD0+0x0D] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R6,reladdr
	OPCODES[0xD0+0x0E] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// DJNZ R7,reladdr
	OPCODES[0xD0+0x0F] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOVX
	OPCODES[0xE0+0x00] = func(program []byte, pos int) int {
		return 1
	}

	// AJMP
	OPCODES[0xE0+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOVX
	OPCODES[0xE0+0x02] = func(program []byte, pos int) int {
		return 1
	}

	// MOVX
	OPCODES[0xE0+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// CLR
	OPCODES[0xE0+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xE0+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xE0+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	// MOVX
	OPCODES[0xF0+0x00] = func(program []byte, pos int) int {
		return 1
	}

	// ACALL
	OPCODES[0xF0+0x01] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOVX
	OPCODES[0xF0+0x02] = func(program []byte, pos int) int {
		return 1
	}

	// MOVX
	OPCODES[0xF0+0x03] = func(program []byte, pos int) int {
		return 1
	}

	// CPL
	OPCODES[0xF0+0x04] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x05] = func(program []byte, pos int) int {
		operand := program[1]
		fmt.Printf("operand: %02x\n", operand)
		return 2
	}

	// MOV
	OPCODES[0xF0+0x06] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x07] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x08] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x09] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x0A] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x0B] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x0C] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x0D] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x0E] = func(program []byte, pos int) int {
		return 1
	}

	// MOV
	OPCODES[0xF0+0x0F] = func(program []byte, pos int) int {
		return 1
	}

	fileName := os.Args[1]

	log.Printf("Processing file %s\n", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s\n", fileName, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	program := make([]byte, stat.Size())
	_, err = bufio.NewReader(f).Read(program)
	if err != nil && err != io.EOF {
		log.Fatalf("error when reading file: %s\n", err)
	}

	log.Printf("File %s is %d bytes\n", fileName, stat.Size())

	for pos := 0; pos < len(program); {
		opcode := program[pos]

		name, ok := OPCODES_NAME[opcode]
		if !ok {
			fmt.Printf("op %d not in OPCODES_NAME\n", opcode)
			continue
		}

		parseFn, ok := OPCODES[opcode]
		if !ok {
			fmt.Printf("op %d not in OPCODES\n", opcode)
			continue
		}

		fmt.Printf("pos %d op %d (0x%02X) = %s\n", pos, opcode, opcode, name)

		sz := parseFn(program[pos:], pos)

		fmt.Println()

		pos = pos + sz
	}
}
