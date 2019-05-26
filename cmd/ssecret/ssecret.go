package main

import (
	"context"
	"fmt"

	"github.com/dssysolyatin/ssecret/ui"
	"github.com/spf13/cobra"
)

var createConfigCmd = &cobra.Command{
	Use: "config",
	RunE: func(cmd *cobra.Command, args []string) error {
		con := ui.NewConsoleUI()
		// 1. Get master password
		_, err := con.ReadPassword(context.Background(), "Read password: ")
		if err != nil {
			return err
		}

		return nil
	},
}

func main() {
	//st, err := storage.NewDirStorage("test")
	if err := createConfigCmd.Execute(); err != nil {
		fmt.Printf("%s", err)
	}
	// 1. master password
	// 2. generator password
	//pbkdf2.Key()
	// 1. generate
}
