package generate

import "github.com/spf13/cobra"

func Register(rootCmd *cobra.Command) {
	generateCmd := &cobra.Command{
		Use: "generate",
	}

	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(createGenerateSecretCommand())
}
