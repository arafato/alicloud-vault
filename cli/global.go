package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	// "github.com/99designs/alicloud-vault/prompt"
	// "github.com/99designs/alicloud-vault/vault"
	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultKeyringName = "alicloud-vault"
)

var (
	keyringImpl keyring.Keyring
	// alicloudConfigFile *vault.ConfigFile
	// configLoader       *vault.ConfigLoader
	// promptsAvailable   = prompt.Available()
)

var GlobalFlags struct {
	Debug        bool
	Backend      string
	PromptDriver string
	KeychainName string
	PassDir      string
	PassCmd      string
	PassPrefix   string
}

func ConfigureGlobals(app *kingpin.Application) {
	backendsAvailable := []string{}
	for _, backendType := range keyring.AvailableBackends() {
		backendsAvailable = append(backendsAvailable, string(backendType))
	}

	app.Flag("debug", "Show debugging output").
		BoolVar(&GlobalFlags.Debug)

	app.Flag("backend", fmt.Sprintf("Secret backend to use %v", backendsAvailable)).
		Envar("ALICLOUD_VAULT_BACKEND").
		EnumVar(&GlobalFlags.Backend, backendsAvailable...)

		// app.Flag("prompt", fmt.Sprintf("Prompt driver to use %v", promptsAvailable)).
		// 	Default("terminal").
		// 	Envar("ALICLOUD_VAULT_PROMPT").
		// 	EnumVar(&GlobalFlags.PromptDriver, promptsAvailable...)

	app.Flag("prompt", "").
		Default("terminal").
		Envar("ALICLOUD_VAULT_PROMPT").
		StringVar(&GlobalFlags.PromptDriver)

	app.Flag("keychain", "Name of macOS keychain to use, if it doesn't exist it will be created").
		Default("alicloud-vault").
		Envar("ALICLOUD_VAULT_KEYCHAIN_NAME").
		StringVar(&GlobalFlags.KeychainName)

	app.Flag("pass-dir", "Pass password store directory").
		Envar("ALICLOUD_VAULT_PASS_PASSWORD_STORE_DIR").
		StringVar(&GlobalFlags.PassDir)

	app.Flag("pass-cmd", "Name of the pass executable").
		Envar("ALICLOUD_VAULT_PASS_CMD").
		StringVar(&GlobalFlags.PassCmd)

	app.Flag("pass-prefix", "Prefix to prepend to the item path stored in pass").
		Envar("ALICLOUD_VAULT_PASS_PREFIX").
		StringVar(&GlobalFlags.PassPrefix)

	app.PreAction(func(c *kingpin.ParseContext) (err error) {
		if !GlobalFlags.Debug {
			log.SetOutput(ioutil.Discard)
		} else {
			keyring.Debug = true
		}
		log.Printf("alicloud-vault %s", app.Model().Version)
		if keyringImpl == nil {
			var allowedBackends []keyring.BackendType
			if GlobalFlags.Backend != "" {
				allowedBackends = append(allowedBackends, keyring.BackendType(GlobalFlags.Backend))
			}
			keyringImpl, err = keyring.Open(keyring.Config{
				ServiceName:              "alicloud-vault",
				AllowedBackends:          allowedBackends,
				KeychainName:             GlobalFlags.KeychainName,
				FileDir:                  "~/.alicloudvault/keys/",
				FilePasswordFunc:         fileKeyringPassphrasePrompt,
				PassDir:                  GlobalFlags.PassDir,
				PassCmd:                  GlobalFlags.PassCmd,
				PassPrefix:               GlobalFlags.PassPrefix,
				LibSecretCollectionName:  "alicloudvault",
				KWalletAppID:             "alicloud-vault",
				KWalletFolder:            "alicloud-vault",
				KeychainTrustApplication: true,
				WinCredPrefix:            "alicloud-vault",
			})
			if err != nil {
				return err
			}
		}

		// alicloudConfigFile, err = vault.LoadConfigFromEnv()
		// configLoader = &vault.ConfigLoader{File: alicloudConfigFile}

		return err
	})
}

func fileKeyringPassphrasePrompt(prompt string) (string, error) {
	if password := os.Getenv("ALICLOUD_VAULT_FILE_PASSPHRASE"); password != "" {
		return password, nil
	}

	fmt.Fprintf(os.Stderr, "%s: ", prompt)
	b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(b), nil
}

func getProfileNames() []string {
	return nil
	// return alicloudConfigFile.ProfileNames()
}
