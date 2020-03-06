# Alicloud Vault
Alicloud Vault is a tool to securely store and access Alibaba Cloud credentials in a development environment.

Alicloud Vault stores RAM credentials in your operating system's secure keystore and then generates temporary credentials from those to expose to your shell and applications. It's designed to be complementary to the Aliyun CLI tools, and is aware of your config file and profiles in ~/.aliyun/config.

Full release and documentation will follow very soon.

## Installing

You can install Alicloud Vault:
- by downloading the [latest release](https://github.com/arafato/alicloud-vault/releases/latest)

## Vaulting Backends

The supported vaulting backends are:

* [macOS Keychain](https://support.apple.com/en-au/guide/keychain-access/welcome/mac)
* [Windows Credential Manager](https://support.microsoft.com/en-au/help/4026814/windows-accessing-credential-manager)
* Secret Service ([Gnome Keyring](https://wiki.gnome.org/Projects/GnomeKeyring), [KWallet](https://kde.org/applications/system/org.kde.kwalletmanager5))
* [KWallet](https://kde.org/applications/system/org.kde.kwalletmanager5)
* [Pass](https://www.passwordstore.org/)
* Encrypted file

Use the `--backend` flag or `ALICLOUD_VAULT_BACKEND` environment variable to specify.

## Basic Usage

```bash
# Store Alicloud credentials for the "jonsmith" profile
$ alicloud-vault add jonsmith
Enter Access Key Id: ABDCDEFDASDASF
Enter Secret Key: %%%

# Execute a command (using temporary credentials), note that you need to explicitly define the keys and token as flag since aliyun is not aware of env variables
$ alicloud-vault exec jonsmith -- aliyun --profile johnsmith --access-key-id $ALICLOUD_ACCESS_KEY --access-key-secret $ALICLOUD_SECRET_KEY --sts-token $ALICLOUD_STS_TOKEN oss ls
bucket_1
bucket_2

# List credentials
$ alicloud-vault ls
Profile                  AccessKeyId               Created
=======                  ===========               ========
johnsmith                LTAI7RuRJSPz************  2020-03-06
```

## How it works
`alicloud-vault` uses Alibaba Cloud's STS service to generate temporary credentials via the [`AssumeRole`](https://www.alibabacloud.com/help/doc-detail/28763.htm) API call. These expire in a short period of time, so the risk of leaking credentials is reduced. Note that not all services support STS. You can find the currently supported services here: https://www.alibabacloud.com/help/doc-detail/135527.htm

Alicloud Vault then exposes the temporary credentials to the sub-process through environment variables in the following way
   ```bash
   $ alicloud-vault exec jonsmith -- env | grep ALICLOUD
   ALICLOUD_VAULT=jonsmith
   ALICLOUD_REGION=us-east-1
   ALICLOUD_ACCESS_KEY=%%%
   ALICLOUD_SECRET_KEY=%%%
   ALICLOUD_STS_TOKEN=%%%
   ALICLOUD_SESSION_EXPIRATION=2020-03-06T10:02:33Z
   ```

## Profiles
`alicloud-vault` is tightly integrated with `Aliyun` CLI and requires a matching profile in ~/.aliyun/config. This allows you to seamlessly use your Aliyun CLI configuration together with `alicloud-vault` and `aliyun` CLI.
It will read configuration data such as `ram_role_arn` directly from this profile. Values defined here take precedence over environment variables. *Access Key ID* and *Access Key Secret* are always read from the keychain, obviously. 
Example:

```
{
	"name": "johnsmith",
	"mode": "StsToken",
	"access_key_id": "",
	"access_key_secret": "",
	"sts_token": "",
	"ram_role_name": "",
	"ram_role_arn": "acs:ram::5509671337805201:role/johnsmith-role",
	"ram_session_name": "johnsmith_session",
	"private_key": "",
	"key_pair_name": "",
	"expired_seconds": 900,
	"verified": "",
	"region_id": "eu-central-1",
	"output_format": "json",
	"language": "en",
	"site": "intl",
	"retry_timeout": 0,
	"retry_count": 0
}
```

## Development

The [macOS release builds](https://github.com/arafato/alicloud-vault/releases/latest) are currently not code-signed so there will be extra prompts in Keychain.

If you are developing or compiling the alicloud-vault binary yourself, you can [generate a self-signed certificate](https://support.apple.com/en-au/guide/keychain-access/kyca8916/mac) by accessing Keychain Access > Certificate Assistant > Create Certificate > Code Signing Certificate. You can then sign your binary with:

    $ go build .
    $ codesign --sign "Name of my certificate" ./alicloud-vault

## References and Inspiration
- https://github.com/99designs/aws-vault 