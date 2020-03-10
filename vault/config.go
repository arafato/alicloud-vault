package vault

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aliyun/aliyun-cli/config"
)

const configFile = "config.json"

// Config is a collection of configuration options for creating temporary credentials
type Config struct {
	// ProfileName specifies the name of the profile config
	ProfileName string

	// Region is the Alicloud region
	Region string

	// AssumeRole config
	RoleARN         string
	RoleSessionName string

	// Specifies the wanted duration for credentials generated with AssumeRole
	AssumeRoleDuration int
}

type ConfigLoader struct {
	BaseConfig   Config
	AliyunConfig config.Configuration
}

func (cl *ConfigLoader) Init() error {
	w := new(bytes.Buffer)
	var err error
	cl.AliyunConfig, err = config.LoadConfiguration(config.GetConfigPath()+"/"+configFile, w)
	if err != nil {
		cl.AliyunConfig = config.NewConfiguration()
		err = config.SaveConfiguration(cl.AliyunConfig)
	}

	return err
}

// Init loads the profile from the config file and environment variables into config
func (cl *ConfigLoader) LoadProfile(profileName string) (*Config, error) {
	cl.populateFromEnv(&cl.BaseConfig)
	err := cl.populateFromConfigFile(&cl.BaseConfig, profileName)
	if err != nil {
		return nil, err
	}

	return &cl.BaseConfig, nil
}

func (cl *ConfigLoader) DeleteProfile(profileName string) error {
	newConfig := config.NewConfiguration()
	for _, p := range cl.AliyunConfig.Profiles {
		if p.Name != profileName {
			newConfig.PutProfile(p)
		}
	}
	err := config.SaveConfiguration(newConfig)
	return err
}

func (cl *ConfigLoader) AddNewProfile(name string) error {
	p, exists := cl.AliyunConfig.GetProfile(name)
	if !exists {
		p.Mode = "StsToken"
		cl.AliyunConfig.PutProfile(p)
		err := config.SaveConfiguration(cl.AliyunConfig)
		if err == nil {
			fmt.Printf("Created new profile '%s' in Aliyun CLI config \n", name)
		}
		return err
	}

	return nil
}

func (cl *ConfigLoader) GetProfileNames() []string {
	var profileNames []string
	for _, p := range cl.AliyunConfig.Profiles {
		profileNames = append(profileNames, p.Name)
	}

	return profileNames
}

func (cl *ConfigLoader) populateFromConfigFile(configuration *Config, profileName string) error {
	w := new(bytes.Buffer)
	profile, err := config.LoadProfile(config.GetConfigPath()+"/"+configFile, w, profileName)
	if err != nil {
		return err
	}

	if configuration.ProfileName == "" {
		configuration.ProfileName = profile.Name
	}
	if configuration.Region == "" {
		configuration.Region = profile.RegionId
	}
	if configuration.RoleARN == "" {
		configuration.RoleARN = profile.RamRoleArn
	}
	if configuration.RoleSessionName == "" {
		configuration.RoleSessionName = profile.RoleSessionName
	}
	if configuration.AssumeRoleDuration == 0 {
		configuration.AssumeRoleDuration = profile.ExpiredSeconds
	}

	return nil
}

func (cl *ConfigLoader) populateFromEnv(profile *Config) {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" && profile.Region == "" {
		log.Printf("Using region %q from ALICLOUD_REGION", region)
		profile.Region = region
	}
	if roleARN := os.Getenv("ALICLOUD_ROLE_ARN"); roleARN != "" && profile.RoleARN == "" {
		log.Printf("Using role_arn %q from ALICLOUD_ROLE_ARN", roleARN)
		profile.RoleARN = roleARN
	}
	if roleSessionName := os.Getenv("ALICLOUD_ROLE_SESSION_NAME"); roleSessionName != "" && profile.RoleSessionName == "" {
		log.Printf("Using role_session_name %q from ALICLOUD_ROLE_SESSION_NAME", roleSessionName)
		profile.RoleSessionName = roleSessionName
	}
	if assumeRoleDuration := os.Getenv("ALICLOUD_ASSUME_ROLE_TTL"); assumeRoleDuration != "" && profile.AssumeRoleDuration == 0 {
		var err error
		profile.AssumeRoleDuration, err = strconv.Atoi(assumeRoleDuration)
		if err == nil {
			log.Printf("Using duration_seconds %q from ALICLOUD_ASSUME_ROLE_TTL", profile.AssumeRoleDuration)
		}
	}
}
