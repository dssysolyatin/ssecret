package generate

import (
	"context"
	"io/ioutil"

	"github.com/dssysolyatin/ssecret/di"
	"github.com/dssysolyatin/ssecret/generator"
	"github.com/spf13/cobra"
)

func createGenerateSecretCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "secret_key [output file]",
		Long: `Generate secret key`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g := generator.NewRandSecretKeyGenerator()
			sk, err := g.Generate()
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(args[0], sk, 0400); err != nil {
				return err
			}

			ui := di.NewDI().GetUI()
			ui.Print(context.Background(), "Secret key is generated successfully")

			return nil
		},
	}
}
