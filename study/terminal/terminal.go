package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func main() {
	fmt.Println("Your password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
	password := string(bytePassword)
	fmt.Println() // it's necessary to add a new line after user's input
	fmt.Printf("Your password has leaked, it is '%s'", password)
}
