package vault

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

type AssumeRoleProvider struct {
	config *Config
	creds  *Credentials
}

func (p AssumeRoleProvider) Retrieve() (*TempCredentials, error) {

	client, err := sts.NewClientWithAccessKey(p.config.Region, p.creds.AccessKeyID, p.creds.SecretAccessKey)
	request := sts.CreateAssumeRoleRequest()
	request.RoleArn = p.config.RoleARN
	request.RoleSessionName = p.config.RoleSessionName
	request.DurationSeconds = requests.NewInteger(p.config.AssumeRoleDuration)
	request.Scheme = "https"
	response, err := client.AssumeRole(request)
	if err != nil {
		return nil, err
	}

	return &TempCredentials{
		Creds: &Credentials{
			AccessKeyID:     response.Credentials.AccessKeyId,
			SecretAccessKey: response.Credentials.AccessKeySecret,
		},
		StsToken: response.Credentials.SecurityToken,
		// format 2015-04-09T11:52:19Z, see https://www.alibabacloud.com/help/doc-detail/28763.htm
		Duration: response.Credentials.Expiration,
	}, nil
}
