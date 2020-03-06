package cli

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/arafato/alicloud-vault/vault"

	"gopkg.in/alecthomas/kingpin.v2"
)

type RotateCommandInput struct {
	ProfileName string
	Keyring     *vault.CredentialKeyring
}

func ConfigureRotateCommand(app *kingpin.Application) {
	input := RotateCommandInput{}

	cmd := app.Command("rotate", "Rotates credentials")

	cmd.Arg("profile", "Name of the profile").
		Required().
		HintAction(getProfileNames).
		StringVar(&input.ProfileName)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.Keyring = &vault.CredentialKeyring{Keyring: keyringImpl}
		app.FatalIfError(RotateCommand(input), "rotate")
		return nil
	})
}

func RotateCommand(input RotateCommandInput) error {

	// Get the existing credentials access key ID
	oldCreds, err := input.Keyring.Get(input.ProfileName)
	if err != nil {
		return err
	}

	log.Printf("Rotating access key '%s' for profile '%s' \n", vault.FormatKeyForDisplay(oldCreds.AccessKeyID), input.ProfileName)

	// Create a new access key
	config, err := configLoader.LoadProfile(input.ProfileName)
	if err != nil {
		return err
	}
	client, err := ram.NewClientWithAccessKey(config.Region, oldCreds.AccessKeyID, oldCreds.SecretAccessKey)
	request := ram.CreateCreateAccessKeyRequest()
	request.Scheme = "https"
	request.UserName = input.ProfileName
	response, err := client.CreateAccessKey(request)
	if err != nil {
		return err
	}
	newCreds := vault.NewCredentials(response.AccessKey.AccessKeyId, response.AccessKey.AccessKeySecret)

	fmt.Printf("Created new access key %s\n", vault.FormatKeyForDisplay(newCreds.AccessKeyID))
	err = input.Keyring.Set(input.ProfileName, newCreds)
	if err != nil {
		return fmt.Errorf("Error storing new access key %s: %w", vault.FormatKeyForDisplay(newCreds.AccessKeyID), err)
	}

	// Use new credentials to delete old access key
	fmt.Printf("Deleting old access key %s\n", oldCreds.AccessKeyID)

	client, err = ram.NewClientWithAccessKey(config.Region, newCreds.AccessKeyID, newCreds.SecretAccessKey)

	request2 := ram.CreateDeleteAccessKeyRequest()
	request2.Scheme = "HTTPS"
	request2.UserAccessKeyId = oldCreds.AccessKeyID
	request2.UserName = input.ProfileName

	_, err = client.DeleteAccessKey(request2)
	if err != nil {
		return fmt.Errorf("Can't delete old access key %s: %w", oldCreds.AccessKeyID, err)
	}

	fmt.Printf("Deleted old access key %s\n", oldCreds.AccessKeyID)
	fmt.Println("Finished rotating access key")

	return nil
}
