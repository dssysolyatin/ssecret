package main

import (
	"fmt"

	"github.com/dssysolyatin/ssecret/cmd/ssecret/command/crypt"

	"github.com/dssysolyatin/ssecret/cmd/ssecret/command/generate"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "ssecret",
	}

	generate.Register(rootCmd)
	crypt.Register(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s\n", err)
	}
}
