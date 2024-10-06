package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	vm "aimandaniel.com/go8051/vm"
)

func dumpVm(vm *vm.Machine) {
	fmt.Printf("Memory=%+v\n", vm.Data[:0x80])
	fmt.Printf("Registers=%+v\n", vm.Data[0x80:])
	fmt.Printf("A=%#08b %#02x %d\n", vm.Data[0xE0], vm.Data[0xE0], vm.Data[0xE0])
	fmt.Printf("PC=%#08b %#02x %d\n", vm.PC, vm.PC, vm.PC)
	fmt.Printf("SP=%#08b %#02x %d\n", vm.SP, vm.SP, vm.SP)
}

func main() {
	vm := vm.NewMachine()
	// Create a channel to listen for Ctrl+C (SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to indicate when to exit
	doneChan := make(chan bool)

	// Start a goroutine to handle signals
	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		doneChan <- true
	}()

	// Start the REPL loop
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("REPL started. Type something and press Enter. Press Ctrl+C to exit.")

Repl:
	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')

		if input == "" || input == "\n" || input == " " {
			continue
		}

		if strings.HasPrefix(input, "/quit") || strings.HasPrefix(input, "/exit") || strings.HasPrefix(input, "/q") {
			fmt.Println("quitting REPL...\n")
			return
		} else if strings.HasPrefix(input, "/dump") {
			dumpVm(vm)
			continue
		}

		parts := strings.Split(strings.Trim(input, "\n"), " ")
		partsByte := make([]byte, len(parts))

		for i, part := range parts {
			b, err := hex.DecodeString(part)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				goto Repl
			}
			partsByte[i] = b[0]
		}

		fmt.Printf("executing instruction %+v\n", partsByte)

		err := vm.Feed(partsByte)
		if err != nil {
			fmt.Printf("vm error: %s\n", err)
		}

		dumpVm(vm)
		fmt.Println()

		select {
		case <-doneChan:
			return
		default:
		}
	}
}
