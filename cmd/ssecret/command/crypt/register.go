package crypt

import "github.com/spf13/cobra"

func Register(rootCmd *cobra.Command) {

	rootCmd.AddCommand(createCryptCommand())
}
