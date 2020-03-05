package vault

import "log"

type provider interface {
	Retrieve() *TempCredentials
}

func GenerateTempCredentials(config *Config, k *CredentialKeyring) (*TempCredentials, error) {

	log.Printf("Looking up keyring for '%s'", config.ProfileName)
	creds, err := k.Get(config.ProfileName)
	if err != nil {
		return nil, err
	}

	p := AssumeRoleProvider{
		config: config,
		c:      creds,
	}

	return p.Retrieve()
}
