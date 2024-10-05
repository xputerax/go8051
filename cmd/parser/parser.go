package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	parser "aimandaniel.com/go8051/parser"
)

var LOOKUP_TABLE map[byte]parser.Opcode = parser.LOOKUP_TABLE

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: ./parser <binary>")
		fmt.Println("usage: ./parser examples/blink.bin")
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
			fmt.Printf("op %#02x not in LOOKUP_TABLE\n", byteCode)
			continue
		}

		sz, operands, err := op.Parse(program[pos:], pos)
		if err != nil {
			fmt.Printf("error encountered when parsing opcode %#02x: %s", err)
			continue
		}

		if len(operands) == 0 {
			fmt.Printf("pos %d op %d (%#02x) = %s\n", pos, byteCode, byteCode, op.Name)
		} else {
			operandsStr := ""
			for i, operand := range operands {
				operandsStr += fmt.Sprintf("%#02x", operand)
				if i != len(operands)-1 { // not the last operand, append space
					operandsStr += " "
				}
			}

			fmt.Printf("pos %d op %d (%#02x) %s = %s\n", pos, byteCode, byteCode, operandsStr, op.Name)
		}

		pos = pos + sz
	}
}
