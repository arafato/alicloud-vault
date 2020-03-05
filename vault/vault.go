package vault

type provider interface {
	Retrieve() *TempCredentials
}

func GenerateTempCredentials(config *Config, k *CredentialKeyring) (*TempCredentials, error) {

	creds, err := k.Get(config.ProfileName)
	if err != nil {
		return nil, err
	}

	// This part can be extended in future versions to support different kinds of providers if needed.
	p := AssumeRoleProvider{
		config: config,
		creds:  creds,
	}

	return p.Retrieve()
}
