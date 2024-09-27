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
type OpcodeParser func(program []byte, pos int) int

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

type Opcode struct {
	Name  string
	Parse OpcodeParser
}

func getLookupTable() map[byte]Opcode {
	tbl := make(map[byte]Opcode)

	tbl[0x00] = Opcode{Name: "NOP", Parse: oneByteOp}
	tbl[0x01] = Opcode{Name: "AJMP codeaddr", Parse: twoByteOp}
	tbl[0x02] = Opcode{Name: "LJMP codeaddr", Parse: threeByteOp}
	tbl[0x03] = Opcode{Name: "RR", Parse: oneByteOp}
	tbl[0x04] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x05] = Opcode{Name: "INC", Parse: twoByteOp}
	tbl[0x06] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x07] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x08] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x09] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x0a] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x0b] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x0c] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x0d] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x0e] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x0f] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0x10] = Opcode{Name: "JBC", Parse: threeByteOp}
	tbl[0x11] = Opcode{Name: "ACALL page0", Parse: twoByteOp}
	tbl[0x12] = Opcode{Name: "LCALL codeaddr", Parse: threeByteOp}
	tbl[0x13] = Opcode{Name: "RRC", Parse: oneByteOp}
	tbl[0x14] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x15] = Opcode{Name: "DEC", Parse: twoByteOp}
	tbl[0x16] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x17] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x18] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x19] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x1a] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x1b] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x1c] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x1d] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x1e] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x1f] = Opcode{Name: "DEC", Parse: oneByteOp}
	tbl[0x20] = Opcode{Name: "JB", Parse: threeByteOp}
	tbl[0x21] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0x22] = Opcode{Name: "RET", Parse: oneByteOp}
	tbl[0x23] = Opcode{Name: "RL", Parse: oneByteOp}
	tbl[0x24] = Opcode{Name: "ADD", Parse: twoByteOp}
	tbl[0x25] = Opcode{Name: "ADD", Parse: twoByteOp}
	tbl[0x26] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x27] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x28] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x29] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x2a] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x2b] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x2c] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x2d] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x2e] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x2f] = Opcode{Name: "ADD", Parse: oneByteOp}
	tbl[0x30] = Opcode{Name: "JNB", Parse: threeByteOp}
	tbl[0x31] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0x32] = Opcode{Name: "RETI", Parse: oneByteOp}
	tbl[0x33] = Opcode{Name: "RLC", Parse: oneByteOp}
	tbl[0x34] = Opcode{Name: "ADDC", Parse: twoByteOp}
	tbl[0x35] = Opcode{Name: "ADDC", Parse: twoByteOp}
	tbl[0x36] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x37] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x38] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x39] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x3a] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x3b] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x3c] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x3d] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x3e] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x3f] = Opcode{Name: "ADDC", Parse: oneByteOp}
	tbl[0x40] = Opcode{Name: "JC reladdr", Parse: twoByteOp}
	tbl[0x41] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0x42] = Opcode{Name: "ORL", Parse: twoByteOp}
	tbl[0x43] = Opcode{Name: "ORL", Parse: threeByteOp}
	tbl[0x44] = Opcode{Name: "ORL", Parse: twoByteOp}
	tbl[0x45] = Opcode{Name: "ORL", Parse: twoByteOp}
	tbl[0x46] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x47] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x48] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x49] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x4a] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x4b] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x4c] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x4d] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x4e] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x4f] = Opcode{Name: "ORL", Parse: oneByteOp}
	tbl[0x50] = Opcode{Name: "JNC reladdr", Parse: twoByteOp}
	tbl[0x51] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0x52] = Opcode{Name: "ANL", Parse: twoByteOp}
	tbl[0x53] = Opcode{Name: "ANL", Parse: threeByteOp}
	tbl[0x54] = Opcode{Name: "ANL", Parse: twoByteOp}
	tbl[0x55] = Opcode{Name: "ANL", Parse: twoByteOp}
	tbl[0x56] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x57] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x58] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x59] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x5a] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x5b] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x5c] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x5d] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x5e] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x5f] = Opcode{Name: "ANL", Parse: oneByteOp}
	tbl[0x60] = Opcode{Name: "JZ reladdr", Parse: twoByteOp}
	tbl[0x61] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0x62] = Opcode{Name: "XRL", Parse: twoByteOp}
	tbl[0x63] = Opcode{Name: "XRL", Parse: threeByteOp}
	tbl[0x64] = Opcode{Name: "XRL", Parse: twoByteOp}
	tbl[0x65] = Opcode{Name: "XRL", Parse: twoByteOp}
	tbl[0x66] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x67] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x68] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x69] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x6a] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x6b] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x6c] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x6d] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x6e] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x6f] = Opcode{Name: "XRL", Parse: oneByteOp}
	tbl[0x70] = Opcode{Name: "JNZ reladdr", Parse: twoByteOp}
	tbl[0x71] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0x72] = Opcode{Name: "ORL", Parse: twoByteOp}
	tbl[0x73] = Opcode{Name: "JMP", Parse: oneByteOp}
	tbl[0x74] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x75] = Opcode{Name: "MOV", Parse: threeByteOp}
	tbl[0x76] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x77] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x78] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x79] = Opcode{Name: "MOV R1,#data", Parse: twoByteOp}
	tbl[0x7a] = Opcode{Name: "MOV R2,#data", Parse: twoByteOp}
	tbl[0x7b] = Opcode{Name: "MOV R3,#data", Parse: twoByteOp}
	tbl[0x7c] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x7d] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x7e] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x7f] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x80] = Opcode{Name: "SJMP reladdr", Parse: twoByteOp}
	tbl[0x81] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0x82] = Opcode{Name: "ANL", Parse: twoByteOp}
	tbl[0x83] = Opcode{Name: "MOVC", Parse: oneByteOp}
	tbl[0x84] = Opcode{Name: "DIV", Parse: oneByteOp}
	tbl[0x85] = Opcode{Name: "MOV", Parse: threeByteOp}
	tbl[0x86] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x87] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x88] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x89] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x8a] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x8b] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x8c] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x8d] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x8e] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x8f] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x90] = Opcode{Name: "MOV", Parse: threeByteOp}
	tbl[0x91] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0x92] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0x93] = Opcode{Name: "MOVC", Parse: oneByteOp}
	tbl[0x94] = Opcode{Name: "SUBB", Parse: twoByteOp}
	tbl[0x95] = Opcode{Name: "SUBB", Parse: twoByteOp}
	tbl[0x96] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x97] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x98] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x99] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x9a] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x9b] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x9c] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x9d] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x9e] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0x9f] = Opcode{Name: "SUBB", Parse: oneByteOp}
	tbl[0xa0] = Opcode{Name: "ORL", Parse: twoByteOp}
	tbl[0xa1] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0xa2] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xa3] = Opcode{Name: "INC", Parse: oneByteOp}
	tbl[0xa4] = Opcode{Name: "MUL", Parse: oneByteOp}
	tbl[0xa5] = Opcode{Name: "?", Parse: oneByteOp}
	tbl[0xa6] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xa7] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xa8] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xa9] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xaa] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xab] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xac] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xad] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xae] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xaf] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xb0] = Opcode{Name: "ANL", Parse: twoByteOp}
	tbl[0xb1] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0xb2] = Opcode{Name: "CPL bitaddr", Parse: twoByteOp}
	tbl[0xb3] = Opcode{Name: "CPL", Parse: oneByteOp}
	tbl[0xb4] = Opcode{Name: "CJNE A,#data,reladdr", Parse: threeByteOp}
	tbl[0xb5] = Opcode{Name: "CJNE A,iram addr,reladdr", Parse: threeByteOp}
	tbl[0xb6] = Opcode{Name: "CJNE @R0,#data,reladdr", Parse: threeByteOp}
	tbl[0xb7] = Opcode{Name: "CJNE @R1,#data,reladdr", Parse: threeByteOp}
	tbl[0xb8] = Opcode{Name: "CJNE R0,#data,reladdr", Parse: threeByteOp}
	tbl[0xb9] = Opcode{Name: "CJNE R1,#data,reladdr", Parse: threeByteOp}
	tbl[0xba] = Opcode{Name: "CJNE R2,#data,reladdr", Parse: threeByteOp}
	tbl[0xbb] = Opcode{Name: "CJNE R3,#data,reladdr", Parse: threeByteOp}
	tbl[0xbc] = Opcode{Name: "CJNE R4,#data,reladdr", Parse: threeByteOp}
	tbl[0xbd] = Opcode{Name: "CJNE R5,#data,reladdr", Parse: threeByteOp}
	tbl[0xbe] = Opcode{Name: "CJNE R6,#data,reladdr", Parse: threeByteOp}
	tbl[0xbf] = Opcode{Name: "CJNE R7,#data,reladdr", Parse: threeByteOp}
	tbl[0xc0] = Opcode{Name: "PUSH", Parse: twoByteOp}
	tbl[0xc1] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0xc2] = Opcode{Name: "CLR", Parse: twoByteOp}
	tbl[0xc3] = Opcode{Name: "CLR", Parse: oneByteOp}
	tbl[0xc4] = Opcode{Name: "SWAP", Parse: oneByteOp}
	tbl[0xc5] = Opcode{Name: "XCH", Parse: twoByteOp}
	tbl[0xc6] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xc7] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xc8] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xc9] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xca] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xcb] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xcc] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xcd] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xce] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xcf] = Opcode{Name: "XCH", Parse: oneByteOp}
	tbl[0xd0] = Opcode{Name: "POP", Parse: twoByteOp}
	tbl[0xd1] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0xd2] = Opcode{Name: "SETB", Parse: twoByteOp}
	tbl[0xd3] = Opcode{Name: "SETB", Parse: oneByteOp}
	tbl[0xd4] = Opcode{Name: "DA", Parse: oneByteOp}
	tbl[0xd5] = Opcode{Name: "DJNZ", Parse: threeByteOp}
	tbl[0xd6] = Opcode{Name: "XCHD", Parse: oneByteOp}
	tbl[0xd7] = Opcode{Name: "XCHD", Parse: oneByteOp}
	tbl[0xd8] = Opcode{Name: "DJNZ R0,reladdr", Parse: twoByteOp}
	tbl[0xd9] = Opcode{Name: "DJNZ R1,reladdr", Parse: twoByteOp}
	tbl[0xda] = Opcode{Name: "DJNZ R2,reladdr", Parse: twoByteOp}
	tbl[0xdb] = Opcode{Name: "DJNZ R3,reladdr", Parse: twoByteOp}
	tbl[0xdc] = Opcode{Name: "DJNZ R4,reladdr", Parse: twoByteOp}
	tbl[0xdd] = Opcode{Name: "DJNZ R5,reladdr", Parse: twoByteOp}
	tbl[0xde] = Opcode{Name: "DJNZ R6,reladdr", Parse: twoByteOp}
	tbl[0xdf] = Opcode{Name: "DJNZ R7,reladdr", Parse: twoByteOp}
	tbl[0xe0] = Opcode{Name: "MOVX", Parse: oneByteOp}
	tbl[0xe1] = Opcode{Name: "AJMP", Parse: twoByteOp}
	tbl[0xe2] = Opcode{Name: "MOVX", Parse: oneByteOp}
	tbl[0xe3] = Opcode{Name: "MOVX", Parse: oneByteOp}
	tbl[0xe4] = Opcode{Name: "CLR", Parse: oneByteOp}
	tbl[0xe5] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xe6] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xe7] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xe8] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xe9] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xea] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xeb] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xec] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xed] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xee] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xef] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xf0] = Opcode{Name: "MOVX", Parse: oneByteOp}
	tbl[0xf1] = Opcode{Name: "ACALL", Parse: twoByteOp}
	tbl[0xf2] = Opcode{Name: "MOVX", Parse: oneByteOp}
	tbl[0xf3] = Opcode{Name: "MOVX", Parse: oneByteOp}
	tbl[0xf4] = Opcode{Name: "CPL", Parse: oneByteOp}
	tbl[0xf5] = Opcode{Name: "MOV", Parse: twoByteOp}
	tbl[0xf6] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xf7] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xf8] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xf9] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xfa] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xfb] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xfc] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xfd] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xfe] = Opcode{Name: "MOV", Parse: oneByteOp}
	tbl[0xff] = Opcode{Name: "MOV", Parse: oneByteOp}

	return tbl
}

var LOOKUP_TABLE map[byte]Opcode = getLookupTable()

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: ./vm <binary>")
		fmt.Println("usage: ./vm examples/blink.bin")
		return
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
		byteCode := program[pos]

		op, ok := LOOKUP_TABLE[byteCode]
		if !ok {
			fmt.Printf("op %d not in LOOKUP_TABLE\n", byteCode)
			continue
		}

		fmt.Printf("pos %d op %d (0x%02X) = %s\n", pos, byteCode, byteCode, op.Name)

		sz := op.Parse(program[pos:], pos)

		fmt.Println()

		pos = pos + sz
	}
}
