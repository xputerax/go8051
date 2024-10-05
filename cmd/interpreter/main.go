package main

import (
	"fmt"

	vm "aimandaniel.com/go8051/vm"
)

func main() {
	m := vm.NewMachine()

	// ni kalau receive raw instruction/byte code
	// TODO: pass assembly, software akan translate jadi byte code, feed masuk VM
	err := m.Feed([]byte{0x24, 0xFF})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
}
