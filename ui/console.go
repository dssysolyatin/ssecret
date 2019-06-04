package ui

import (
	"context"
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var _ UI = (*console)(nil)

type console struct {
}

func NewConsoleUI() *console {
	return &console{}
}

func (c *console) ReadPassword(ctx context.Context, desc string) (string, error) {
	fmt.Print(desc)
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Printf("\n")
	return string(passwd), err
}

func (c *console) Print(ctx context.Context, output string) error {
	fmt.Println(output)
	return nil
}
