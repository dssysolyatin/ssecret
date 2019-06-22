package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var _ UI = (*console)(nil)

type console struct {
	reader *bufio.Reader
}

func NewConsoleUI() *console {
	return &console{reader: bufio.NewReader(os.Stdin)}
}

func (c *console) ReadPassword(ctx context.Context, desc string) (string, error) {
	fmt.Print(desc)
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Printf("\n")
	return string(passwd), err
}

func (c *console) ReadString(ctx context.Context, desc string) (string, error) {
	fmt.Print(desc)
	return c.reader.ReadString('\n')
}

func (c *console) Print(ctx context.Context, output string) error {
	fmt.Println(output)
	return nil
}
