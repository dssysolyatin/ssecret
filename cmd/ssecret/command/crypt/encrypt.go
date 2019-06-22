package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"io/ioutil"

	"github.com/dssysolyatin/ssecret/di"
	"github.com/dssysolyatin/ssecret/utils/contextu"
	"github.com/dssysolyatin/ssecret/utils/filepathu"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
)

func createEncryptCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "encrypt [secret-key path] [output-file]",
		Long: `Generate secret key`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := contextu.WithCancelCtrlC(context.Background())
			path, err := filepathu.Abs(args[0])
			if err != nil {
				return err
			}

			secretKey, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			ui := di.NewDI().GetUI()
			encryptPassword, err := ui.ReadPassword(ctx, "Enter password: ")
			if err != nil {
				return err
			}

			text, err := ui.ReadString(ctx, "Enter text: ")
			if err != nil {
				return err
			}

			secretKey = append(secretKey, []byte(encryptPassword)...)
			secretKey = pbkdf2.Key(secretKey, nil, 4096, 32, sha1.New)

			c, err := aes.NewCipher(secretKey)
			if err != nil {
				return err
			}

			gcm, err := cipher.NewGCM(c)
			if err != nil {
				return err
			}

			nonce := make([]byte, gcm.NonceSize())
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return err
			}

			gcm.Seal(nonce, nonce, []byte(text), nil)
			return nil
		},
	}
}
