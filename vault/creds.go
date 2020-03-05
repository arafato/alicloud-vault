package vault

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type TempCredentials struct {
	Creds    *Credentials
	StsToken string
	Duration string
}

func NewCredentials(accessKeyID, secretAccessKey string) *Credentials {
	return &Credentials{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}

func NewTempCredentials(creds *Credentials, stsToken string, duration string) *TempCredentials {
	return &TempCredentials{
		Creds:    creds,
		StsToken: stsToken,
		Duration: duration,
	}
}
