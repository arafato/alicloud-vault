package cli

import (
	"fmt"

	"github.com/arafato/alicloud-vault/helper"
	"github.com/arafato/alicloud-vault/vault"
	"gopkg.in/alecthomas/kingpin.v2"
)

type RemoveCommandInput struct {
	ProfileName     string
	AliyunCliProfil bool
	Keyring         *vault.CredentialKeyring
}

func ConfigureRemoveCommand(app *kingpin.Application) {
	input := RemoveCommandInput{}

	cmd := app.Command("remove", "Removes profiles from keyring")
	cmd.Alias("rm")

	cmd.Arg("profile", "Name of the profile").
		Required().
		HintAction(getProfileNames).
		StringVar(&input.ProfileName)

	cmd.Flag("aliyun", "Delete the according Aliyun CLI profile").
		Short('a').
		BoolVar(&input.AliyunCliProfil)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.Keyring = &vault.CredentialKeyring{Keyring: keyringImpl}
		RemoveCommand(app, input)
		return nil
	})
}

func RemoveCommand(app *kingpin.Application, input RemoveCommandInput) {
	r, err := helper.TerminalPrompt(fmt.Sprintf("Delete credentials for profile %q? (Y|n)", input.ProfileName))
	if err != nil {
		app.Fatalf(err.Error())
		return
	} else if r == "N" || r == "n" {
		return
	}

	if err := input.Keyring.Remove(input.ProfileName); err != nil {
		app.Fatalf(err.Error())
		return
	}
	fmt.Printf("Deleted profile from Keyvault.\n")

	if input.AliyunCliProfil {
		if err := configLoader.DeleteProfile(input.ProfileName); err != nil {
			app.Fatalf(err.Error())
			return
		}
		fmt.Printf("Deleted profile from Aliyun CLI config.\n")
	}
}
