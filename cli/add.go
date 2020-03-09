package cli

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/arafato/alicloud-vault/helper"
	"github.com/arafato/alicloud-vault/vault"
)

type AddCommandInput struct {
	ProfileName string
	Keyring     *vault.CredentialKeyring
	FromEnv     bool
}

func ConfigureAddCommand(app *kingpin.Application) {
	input := AddCommandInput{}

	cmd := app.Command("add", "Adds credentials, prompts if none provided")
	cmd.Arg("profile", "Name of the profile").
		Required().
		StringVar(&input.ProfileName)

	cmd.Flag("env", "Read the credentials from the environment").
		BoolVar(&input.FromEnv)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.Keyring = &vault.CredentialKeyring{Keyring: keyringImpl}
		AddCommand(app, input)
		return nil
	})
}

func AddCommand(app *kingpin.Application, input AddCommandInput) {
	var accessKeyID, secretKey string

	if input.FromEnv {
		if accessKeyID = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKeyID == "" {
			app.Fatalf("Missing value for ALICLOUD_ACCESS_KEY")
			return
		}
		if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
			app.Fatalf("Missing value for ALICLOUD_SECRET_KEY")
			return
		}
	} else {
		var err error
		if accessKeyID, err = helper.TerminalPrompt("Enter Access Key ID: "); err != nil {
			app.Fatalf(err.Error())
			return
		}
		if secretKey, err = helper.TerminalPrompt("Enter Secret Key: "); err != nil {
			app.Fatalf(err.Error())
			return
		}
	}

	creds := vault.NewCredentials(accessKeyID, secretKey)

	if err := input.Keyring.Set(input.ProfileName, creds); err != nil {
		app.Fatalf(err.Error())
		return
	}

	if err := configLoader.AddNewProfile(input.ProfileName); err != nil {
		app.Fatalf(err.Error())
		return
	}

	fmt.Printf("Added credentials to profile %q in vault\n", input.ProfileName)
}
