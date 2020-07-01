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

	cmd.Flag("credentials", "Show only the profiles with stored credential").
		BoolVar(&input.OnlyCredentials)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.Keyring = &vault.CredentialKeyring{Keyring: keyringImpl}
		app.FatalIfError(LsCommand(input), "")
		return nil
	})
}

func LsCommand(input LsCommandInput) error {

	profileNames := configLoader.GetProfileNames()
	credentialsNames, err := input.Keyring.CredentialsKeys()
	if err != nil {
		return err
	}

	if input.OnlyCredentials {
		for _, c := range credentialsNames {
			fmt.Printf("%s\n", c)
		}
		return nil
	}

	if input.OnlyProfiles {
		for _, p := range profileNames {
			fmt.Printf("%s\n", p)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 25, 4, 2, ' ', 0)
	fmt.Fprintln(w, "Profile\tCredentials\t")
	fmt.Fprintln(w, "=======\t===========\t")
	for _, profileName := range profileNames {

		fmt.Fprintf(w, "%s\t", profileName)
		hasCred, err := input.Keyring.Has(profileName)
		if err != nil {
			return err
		}
		if hasCred {
			fmt.Fprintf(w, "%s\t\n", profileName)
		} else {
			fmt.Fprintf(w, "-\t\n")
		}
	}

	// show credentials that don't have profiles
	for _, credName := range credentialsNames {
		contains := false
		for _, profile := range profileNames {
			if credName == profile {
				contains = true
				break
			}
		}
		if !contains {
			fmt.Fprintf(w, "-\t")
			fmt.Fprintf(w, "%s\t\n", credName)
		}
	}

	fmt.Fprintf(w, "\n")
	return w.Flush()
}
