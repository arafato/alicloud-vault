package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/arafato/alicloud-vault/vault"
	"gopkg.in/alecthomas/kingpin.v2"
)

type LsCommandInput struct {
	Keyring         *vault.CredentialKeyring
	OnlyProfiles    bool
	OnlyCredentials bool
}

func ConfigureListCommand(app *kingpin.Application) {
	input := LsCommandInput{}

	cmd := app.Command("list", "List profiles, along with their creation time")
	cmd.Alias("ls")

	cmd.Flag("profiles", "Show only the profile names").
		BoolVar(&input.OnlyProfiles)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.Keyring = &vault.CredentialKeyring{Keyring: keyringImpl}
		app.FatalIfError(LsCommand(input), "")
		return nil
	})
}

func LsCommand(input LsCommandInput) error {

	credentialsNames, err := input.Keyring.CredentialsKeys()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 25, 4, 2, ' ', 0)

	if !input.OnlyProfiles {
		fmt.Fprintln(w, "Profile\tAccessKeyId\tCreated\t")
		fmt.Fprintln(w, "=======\t===========\t========\t")
	}
	for _, profileName := range credentialsNames {
		if input.OnlyProfiles {
			fmt.Printf("%s\n", profileName)
			continue
		}

		fmt.Fprintf(w, "%s\t", profileName)
		creds, err := input.Keyring.Get(profileName)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, vault.FormatKeyForDisplay(creds.AccessKeyID))
		fmt.Fprintf(w, "%s\t\n", string(creds.Created))
	}
	fmt.Fprintf(w, "\n")
	return w.Flush()
}
