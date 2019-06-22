package crypt

import "github.com/spf13/cobra"

func Register(rootCmd *cobra.Command) {
	generateCmd := &cobra.Command{
		Use: "crypt",
	}

	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(createEncryptCommand())
}
