package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// @param program - slice of the program starting from the
// position of the current instruction onwards
//
// @param pos - the position of the current instruction
//
// @return sz - size of the instruction (opcode + operands)
// so that the next loop iteration will know
// how many elements to skip
type Operation func(program []byte, pos int) int

func oneByteOp(program []byte, pos int) int {
	return 1
}

func twoByteOp(program []byte, pos int) int {
	operand := program[1]
	fmt.Printf("operand: %02x\n", operand)
	return 2
}

func threeByteOp(program []byte, pos int) int {
	operand := program[1:3]
	fmt.Printf("operand: %02x\n", operand)
	return 3
}

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
	OPCODES[0x00+0x00] = oneByteOp

	// AJMP codeaddr
	OPCODES[0x00+0x01] = twoByteOp

	// LJMP
	OPCODES[0x00+0x02] = threeByteOp

	// RR
	OPCODES[0x00+0x03] = oneByteOp

	// INC
	OPCODES[0x00+0x04] = oneByteOp

	// INC
	OPCODES[0x00+0x05] = twoByteOp

	// INC
	OPCODES[0x00+0x06] = oneByteOp

	// INC
	OPCODES[0x00+0x07] = oneByteOp

	// INC
	OPCODES[0x00+0x08] = oneByteOp

	// INC
	OPCODES[0x00+0x09] = oneByteOp

	// INC
	OPCODES[0x00+0x0A] = oneByteOp

	// INC
	OPCODES[0x00+0x0B] = oneByteOp

	// INC
	OPCODES[0x00+0x0C] = oneByteOp

	// INC
	OPCODES[0x00+0x0D] = oneByteOp

	// INC
	OPCODES[0x00+0x0E] = oneByteOp

	// INC
	OPCODES[0x00+0x0F] = oneByteOp

	// JBC
	OPCODES[0x10+0x00] = threeByteOp

	// ACALL
	OPCODES[0x10+0x01] = twoByteOp

	// LCALL codeaddr
	OPCODES[0x10+0x02] = threeByteOp

	// RRC
	OPCODES[0x10+0x03] = oneByteOp

	// DEC
	OPCODES[0x10+0x04] = oneByteOp

	// DEC
	OPCODES[0x10+0x05] = twoByteOp

	// DEC
	OPCODES[0x10+0x06] = oneByteOp

	// DEC
	OPCODES[0x10+0x07] = oneByteOp

	// DEC
	OPCODES[0x10+0x08] = oneByteOp

	// DEC
	OPCODES[0x10+0x09] = oneByteOp

	// DEC
	OPCODES[0x10+0x0A] = oneByteOp

	// DEC
	OPCODES[0x10+0x0B] = oneByteOp

	// DEC
	OPCODES[0x10+0x0C] = oneByteOp

	// DEC
	OPCODES[0x10+0x0D] = oneByteOp

	// DEC
	OPCODES[0x10+0x0E] = oneByteOp

	// DEC
	OPCODES[0x10+0x0F] = oneByteOp

	// JB
	OPCODES[0x20+0x00] = threeByteOp

	// AJMP
	OPCODES[0x20+0x01] = twoByteOp

	// RET
	OPCODES[0x20+0x02] = oneByteOp

	// RL
	OPCODES[0x20+0x03] = oneByteOp

	// ADD
	OPCODES[0x20+0x04] = twoByteOp

	// ADD
	OPCODES[0x20+0x05] = twoByteOp

	// ADD
	OPCODES[0x20+0x06] = oneByteOp

	// ADD
	OPCODES[0x20+0x07] = oneByteOp

	// ADD
	OPCODES[0x20+0x08] = oneByteOp

	// ADD
	OPCODES[0x20+0x09] = oneByteOp

	// ADD
	OPCODES[0x20+0x0A] = oneByteOp

	// ADD
	OPCODES[0x20+0x0B] = oneByteOp

	// ADD
	OPCODES[0x20+0x0C] = oneByteOp

	// ADD
	OPCODES[0x20+0x0D] = oneByteOp

	// ADD
	OPCODES[0x20+0x0E] = oneByteOp

	// ADD
	OPCODES[0x20+0x0F] = oneByteOp

	// JNB
	OPCODES[0x30+0x00] = threeByteOp

	// ACALL
	OPCODES[0x30+0x01] = twoByteOp

	// RETI
	OPCODES[0x30+0x02] = oneByteOp

	// RLC
	OPCODES[0x30+0x03] = oneByteOp

	// ADDC
	OPCODES[0x30+0x04] = twoByteOp

	// ADDC
	OPCODES[0x30+0x05] = twoByteOp

	// ADDC
	OPCODES[0x30+0x06] = oneByteOp

	// ADDC
	OPCODES[0x30+0x07] = oneByteOp

	// ADDC
	OPCODES[0x30+0x08] = oneByteOp

	// ADDC
	OPCODES[0x30+0x09] = oneByteOp

	// ADDC
	OPCODES[0x30+0x0A] = oneByteOp

	// ADDC
	OPCODES[0x30+0x0B] = oneByteOp

	// ADDC
	OPCODES[0x30+0x0C] = oneByteOp

	// ADDC
	OPCODES[0x30+0x0D] = oneByteOp

	// ADDC
	OPCODES[0x30+0x0E] = oneByteOp

	// ADDC
	OPCODES[0x30+0x0F] = oneByteOp

	// JC reladdr
	OPCODES[0x40+0x00] = twoByteOp

	// AJMP
	OPCODES[0x40+0x01] = twoByteOp

	// ORL
	OPCODES[0x40+0x02] = twoByteOp

	// ORL
	OPCODES[0x40+0x03] = threeByteOp

	// ORL
	OPCODES[0x40+0x04] = twoByteOp

	// ORL
	OPCODES[0x40+0x05] = twoByteOp

	// ORL
	OPCODES[0x40+0x06] = oneByteOp

	// ORL
	OPCODES[0x40+0x07] = oneByteOp

	// ORL
	OPCODES[0x40+0x08] = oneByteOp

	// ORL
	OPCODES[0x40+0x09] = oneByteOp

	// ORL
	OPCODES[0x40+0x0A] = oneByteOp

	// ORL
	OPCODES[0x40+0x0B] = oneByteOp

	// ORL
	OPCODES[0x40+0x0C] = oneByteOp

	// ORL
	OPCODES[0x40+0x0D] = oneByteOp

	// ORL
	OPCODES[0x40+0x0E] = oneByteOp

	// ORL
	OPCODES[0x40+0x0F] = oneByteOp

	// JNC reladdr
	OPCODES[0x50+0x00] = twoByteOp

	// ACALL
	OPCODES[0x50+0x01] = twoByteOp

	// ANL
	OPCODES[0x50+0x02] = twoByteOp

	// ANL
	OPCODES[0x50+0x03] = threeByteOp

	// ANL
	OPCODES[0x50+0x04] = twoByteOp

	// ANL
	OPCODES[0x50+0x05] = twoByteOp

	// ANL
	OPCODES[0x50+0x06] = oneByteOp

	// ANL
	OPCODES[0x50+0x07] = oneByteOp

	// ANL
	OPCODES[0x50+0x08] = oneByteOp

	// ANL
	OPCODES[0x50+0x09] = oneByteOp

	// ANL
	OPCODES[0x50+0x0A] = oneByteOp

	// ANL
	OPCODES[0x50+0x0B] = oneByteOp

	// ANL
	OPCODES[0x50+0x0C] = oneByteOp

	// ANL
	OPCODES[0x50+0x0D] = oneByteOp

	// ANL
	OPCODES[0x50+0x0E] = oneByteOp

	// ANL
	OPCODES[0x50+0x0F] = oneByteOp

	// JZ
	OPCODES[0x60+0x00] = twoByteOp

	// AJMP
	OPCODES[0x60+0x01] = twoByteOp

	// XRL
	OPCODES[0x60+0x02] = twoByteOp

	// XRL
	OPCODES[0x60+0x03] = threeByteOp

	// XRL
	OPCODES[0x60+0x04] = twoByteOp

	// XRL
	OPCODES[0x60+0x05] = twoByteOp

	// XRL
	OPCODES[0x60+0x06] = oneByteOp

	// XRL
	OPCODES[0x60+0x07] = oneByteOp

	// XRL
	OPCODES[0x60+0x08] = oneByteOp

	// XRL
	OPCODES[0x60+0x09] = oneByteOp

	// XRL
	OPCODES[0x60+0x0A] = oneByteOp

	// XRL
	OPCODES[0x60+0x0B] = oneByteOp

	// XRL
	OPCODES[0x60+0x0C] = oneByteOp

	// XRL
	OPCODES[0x60+0x0D] = oneByteOp

	// XRL
	OPCODES[0x60+0x0E] = oneByteOp

	// XRL
	OPCODES[0x60+0x0F] = oneByteOp

	// JNZ reladdr
	OPCODES[0x70+0x00] = twoByteOp

	// ACALL
	OPCODES[0x70+0x01] = twoByteOp

	// ORL
	OPCODES[0x70+0x02] = twoByteOp

	// JMP
	OPCODES[0x70+0x03] = oneByteOp

	// MOV
	OPCODES[0x70+0x04] = twoByteOp

	// MOV
	OPCODES[0x70+0x05] = threeByteOp

	// MOV
	OPCODES[0x70+0x06] = twoByteOp

	// MOV
	OPCODES[0x70+0x07] = twoByteOp

	// MOV
	OPCODES[0x70+0x08] = twoByteOp

	// MOV
	OPCODES[0x70+0x09] = twoByteOp

	// MOV
	OPCODES[0x70+0x0A] = twoByteOp

	// MOV
	OPCODES[0x70+0x0B] = twoByteOp

	// MOV
	OPCODES[0x70+0x0C] = twoByteOp

	// MOV
	OPCODES[0x70+0x0D] = twoByteOp

	// MOV
	OPCODES[0x70+0x0E] = twoByteOp

	// MOV
	OPCODES[0x70+0x0F] = twoByteOp

	// SJMP reladdr
	OPCODES[0x80+0x00] = twoByteOp

	// AJMP
	OPCODES[0x80+0x01] = twoByteOp

	// ANL
	OPCODES[0x80+0x02] = twoByteOp

	// MOVC
	OPCODES[0x80+0x03] = oneByteOp

	// DIV
	OPCODES[0x80+0x04] = oneByteOp

	// MOV
	OPCODES[0x80+0x05] = threeByteOp

	// MOV
	OPCODES[0x80+0x06] = twoByteOp

	// MOV
	OPCODES[0x80+0x07] = twoByteOp

	// MOV
	OPCODES[0x80+0x08] = twoByteOp

	// MOV
	OPCODES[0x80+0x09] = twoByteOp

	// MOV
	OPCODES[0x80+0x0A] = twoByteOp

	// MOV
	OPCODES[0x80+0x0B] = twoByteOp

	// MOV
	OPCODES[0x80+0x0C] = twoByteOp

	// MOV
	OPCODES[0x80+0x0D] = twoByteOp

	// MOV
	OPCODES[0x80+0x0E] = twoByteOp

	// MOV
	OPCODES[0x80+0x0F] = twoByteOp

	// MOV
	OPCODES[0x90+0x00] = threeByteOp

	// ACALL
	OPCODES[0x90+0x01] = twoByteOp

	// MOV
	OPCODES[0x90+0x02] = twoByteOp

	// MOVC
	OPCODES[0x90+0x03] = oneByteOp

	// SUBB
	OPCODES[0x90+0x04] = twoByteOp

	// SUBB
	OPCODES[0x90+0x05] = twoByteOp

	// SUBB
	OPCODES[0x90+0x06] = oneByteOp

	// SUBB
	OPCODES[0x90+0x07] = oneByteOp

	// SUBB
	OPCODES[0x90+0x08] = oneByteOp

	// SUBB
	OPCODES[0x90+0x09] = oneByteOp

	// SUBB
	OPCODES[0x90+0x0A] = oneByteOp

	// SUBB
	OPCODES[0x90+0x0B] = oneByteOp

	// SUBB
	OPCODES[0x90+0x0C] = oneByteOp

	// SUBB
	OPCODES[0x90+0x0D] = oneByteOp

	// SUBB
	OPCODES[0x90+0x0E] = oneByteOp

	// SUBB
	OPCODES[0x90+0x0F] = oneByteOp

	// ORL
	OPCODES[0xA0+0x00] = twoByteOp

	// AJMP
	OPCODES[0xA0+0x01] = twoByteOp

	// MOV
	OPCODES[0xA0+0x02] = twoByteOp

	// INC
	OPCODES[0xA0+0x03] = oneByteOp

	// MUL
	OPCODES[0xA0+0x04] = oneByteOp

	// ?
	OPCODES[0xA0+0x05] = oneByteOp

	// MOV
	OPCODES[0xA0+0x06] = twoByteOp

	// MOV
	OPCODES[0xA0+0x07] = twoByteOp

	// MOV
	OPCODES[0xA0+0x08] = twoByteOp

	// MOV
	OPCODES[0xA0+0x09] = twoByteOp

	// MOV
	OPCODES[0xA0+0x0A] = twoByteOp

	// MOV
	OPCODES[0xA0+0x0B] = twoByteOp

	// MOV
	OPCODES[0xA0+0x0C] = twoByteOp

	// MOV
	OPCODES[0xA0+0x0D] = twoByteOp

	// MOV
	OPCODES[0xA0+0x0E] = twoByteOp

	// MOV
	OPCODES[0xA0+0x0F] = twoByteOp

	// ANL
	OPCODES[0xB0+0x00] = twoByteOp

	// ACALL
	OPCODES[0xB0+0x01] = twoByteOp

	// CPL bitaddr
	OPCODES[0xB0+0x02] = twoByteOp

	// CPL
	OPCODES[0xB0+0x03] = oneByteOp

	// CJNE A,#data,reladdr
	OPCODES[0xB0+0x04] = threeByteOp

	// CJNE A,iram addr,reladdr
	OPCODES[0xB0+0x05] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x06] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x07] = threeByteOp

	// CJNE R0,#data,reladdr
	OPCODES[0xB0+0x08] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x09] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x0A] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x0B] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x0C] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x0D] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x0E] = threeByteOp

	// CJNE
	OPCODES[0xB0+0x0F] = threeByteOp

	// PUSH
	OPCODES[0xC0+0x00] = twoByteOp

	// AJMP
	OPCODES[0xC0+0x01] = twoByteOp

	// CLR
	OPCODES[0xC0+0x02] = twoByteOp

	// CLR
	OPCODES[0xC0+0x03] = oneByteOp

	// SWAP
	OPCODES[0xC0+0x04] = oneByteOp

	// XCH
	OPCODES[0xC0+0x05] = twoByteOp

	// XCH
	OPCODES[0xC0+0x06] = oneByteOp

	// XCH
	OPCODES[0xC0+0x07] = oneByteOp

	// XCH
	OPCODES[0xC0+0x08] = oneByteOp

	// XCH
	OPCODES[0xC0+0x09] = oneByteOp

	// XCH
	OPCODES[0xC0+0x0A] = oneByteOp

	// XCH
	OPCODES[0xC0+0x0B] = oneByteOp

	// XCH
	OPCODES[0xC0+0x0C] = oneByteOp

	// XCH
	OPCODES[0xC0+0x0D] = oneByteOp

	// XCH
	OPCODES[0xC0+0x0E] = oneByteOp

	// XCH
	OPCODES[0xC0+0x0F] = oneByteOp

	// POP
	OPCODES[0xD0+0x00] = twoByteOp

	// ACALL
	OPCODES[0xD0+0x01] = twoByteOp

	// SETB
	OPCODES[0xD0+0x02] = twoByteOp

	// SETB
	OPCODES[0xD0+0x03] = oneByteOp

	// DA
	OPCODES[0xD0+0x04] = oneByteOp

	// DJNZ
	OPCODES[0xD0+0x05] = threeByteOp

	// XCHD
	OPCODES[0xD0+0x06] = oneByteOp

	// XCHD
	OPCODES[0xD0+0x07] = oneByteOp

	// DJNZ R0,reladdr
	OPCODES[0xD0+0x08] = twoByteOp

	// DJNZ R1,reladdr
	OPCODES[0xD0+0x09] = twoByteOp

	// DJNZ R2,reladdr
	OPCODES[0xD0+0x0A] = twoByteOp

	// DJNZ R3,reladdr
	OPCODES[0xD0+0x0B] = twoByteOp

	// DJNZ R4,reladdr
	OPCODES[0xD0+0x0C] = twoByteOp

	// DJNZ R5,reladdr
	OPCODES[0xD0+0x0D] = twoByteOp

	// DJNZ R6,reladdr
	OPCODES[0xD0+0x0E] = twoByteOp

	// DJNZ R7,reladdr
	OPCODES[0xD0+0x0F] = twoByteOp

	// MOVX
	OPCODES[0xE0+0x00] = oneByteOp

	// AJMP
	OPCODES[0xE0+0x01] = twoByteOp

	// MOVX
	OPCODES[0xE0+0x02] = oneByteOp

	// MOVX
	OPCODES[0xE0+0x03] = oneByteOp

	// CLR
	OPCODES[0xE0+0x04] = oneByteOp

	// MOV
	OPCODES[0xE0+0x05] = twoByteOp

	// MOV
	OPCODES[0xE0+0x06] = oneByteOp

	// MOV
	OPCODES[0xE0+0x07] = oneByteOp

	// MOV
	OPCODES[0xE0+0x08] = oneByteOp

	// MOV
	OPCODES[0xE0+0x09] = oneByteOp

	// MOV
	OPCODES[0xE0+0x0A] = oneByteOp

	// MOV
	OPCODES[0xE0+0x0B] = oneByteOp

	// MOV
	OPCODES[0xE0+0x0C] = oneByteOp

	// MOV
	OPCODES[0xE0+0x0D] = oneByteOp

	// MOV
	OPCODES[0xE0+0x0E] = oneByteOp

	// MOV
	OPCODES[0xE0+0x0F] = oneByteOp

	// MOVX
	OPCODES[0xF0+0x00] = oneByteOp

	// ACALL
	OPCODES[0xF0+0x01] = twoByteOp

	// MOVX
	OPCODES[0xF0+0x02] = oneByteOp

	// MOVX
	OPCODES[0xF0+0x03] = oneByteOp

	// CPL
	OPCODES[0xF0+0x04] = oneByteOp

	// MOV
	OPCODES[0xF0+0x05] = twoByteOp

	// MOV
	OPCODES[0xF0+0x06] = oneByteOp

	// MOV
	OPCODES[0xF0+0x07] = oneByteOp

	// MOV
	OPCODES[0xF0+0x08] = oneByteOp

	// MOV
	OPCODES[0xF0+0x09] = oneByteOp

	// MOV
	OPCODES[0xF0+0x0A] = oneByteOp

	// MOV
	OPCODES[0xF0+0x0B] = oneByteOp

	// MOV
	OPCODES[0xF0+0x0C] = oneByteOp

	// MOV
	OPCODES[0xF0+0x0D] = oneByteOp

	// MOV
	OPCODES[0xF0+0x0E] = oneByteOp

	// MOV
	OPCODES[0xF0+0x0F] = oneByteOp

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
