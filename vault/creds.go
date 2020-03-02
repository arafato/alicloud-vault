package vault

type Credential struct {
	AccessKeyID     string
	SecretAccessKey string
}

func NewCredential(accessKeyID, secretAccessKey string) *Credential {
	return &Credential{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}
