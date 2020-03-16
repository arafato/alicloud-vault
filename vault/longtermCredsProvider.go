package vault

type LongtermCredsProvider struct {
	config *Config
	creds  *Credentials
}

func (p LongtermCredsProvider) Retrieve() (*TempCredentials, error) {

	return &TempCredentials{
		Creds: &Credentials{
			AccessKeyID:     p.creds.AccessKeyID,
			SecretAccessKey: p.creds.SecretAccessKey,
		},
	}, nil
}
