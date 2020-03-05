package vault

import "time"

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Created         string
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
		Created:         time.Now().Format("2006-01-02"),
	}
}

func NewTempCredentials(creds *Credentials, stsToken string, duration string) *TempCredentials {
	return &TempCredentials{
		Creds:    creds,
		StsToken: stsToken,
		Duration: duration,
	}
}
