package vault

type Credential struct {
	AccessKeyID     string
	SecretAccessKey string
	StsToken        string
}

func NewCredential(accessKeyID, secretAccessKey, stsToken string) *Credential {
	return &Credential{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		StsToken:        stsToken,
	}
}
