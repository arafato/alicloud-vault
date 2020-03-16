package vault

import "fmt"

type provider interface {
	Retrieve() (*TempCredentials, error)
}

func GenerateTempCredentials(config *Config, k *CredentialKeyring) (*TempCredentials, error) {

	creds, err := k.Get(config.ProfileName)
	if err != nil {
		return nil, err
	}

	var p provider

	if config.RoleARN == "" {
		p = LongtermCredsProvider{
			config: config,
			creds:  creds,
		}
	} else {
		p = AssumeRoleProvider{
			config: config,
			creds:  creds,
		}
	}

	return p.Retrieve()
}

func FormatKeyForDisplay(k string) string {
	if len(k) == 24 {
		return fmt.Sprintf("%s************\t", string(k[0:12]))
	}
	// This does not seem to be a well-formed access key but we show it nevertheless
	return fmt.Sprintf("%s\t", k)
}
