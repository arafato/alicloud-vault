# Alicloud Vault
Alicloud Vault is a tool to securely store and access Alibaba Cloud credentials in a development environment.

Alicloud Vault stores RAM credentials in your operating system's secure keystore and then generates temporary credentials from those to expose to your shell and applications. It's designed to be complementary to the Aliyun CLI tools, and is aware of your config file and profiles in ~/.aliyun/config.

Full release and documentation will follow very soon.

## Installing

You can install Alicloud Vault:
- by downloading the [latest release](https://github.com/arafato/alicloud-vault/releases/latest)
- on Arch Linux with the [AUR](https://aur.archlinux.org/packages/alicloud-vault/): `yay -S alicloud-vault`

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

# Execute a command (using temporary credentials), note that you need to explicitly define the access key, secret and token as flag (enclosed with ') since aliyun is not aware of env variables
# Environment variables that are not enclosed with ' are not automatically expanded based on the new session context but take the values from the current session.  
$ alicloud-vault exec jonsmith -- aliyun --profile johnsmith --access-key-id '$ALICLOUD_ACCESS_KEY' --access-key-secret '$ALICLOUD_SECRET_KEY' --sts-token '$ALICLOUD_STS_TOKEN' oss ls
bucket_1
bucket_2

# Export environment variables to new shell context and call Aliyun CLI sequentially
$ alicloud-vault exec jonsmith
$ aliyun --profile johnsmith --access-key-id $ALICLOUD_ACCESS_KEY --access-key-secret $ALICLOUD_SECRET_KEY --sts-token $ALICLOUD_STS_TOKEN oss ls


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
`alicloud-vault` is tightly integrated with `Aliyun` CLI and requires a matching profile in ~/.aliyun/config. This allows you to seamlessly use your Aliyun CLI configuration together with `alicloud-vault` and `aliyun` CLI. If you create a new profile in the `alicloud-vault` it will automatically create a matching profile in ~/.aliyun/config if it does not exist yet.
It will read configuration data such as `ram_role_arn` directly from this profile. Values defined here take precedence over environment variables. *Access Key ID* and *Access Key Secret* are always read from the keychain, obviously.

Note that if you do not specify `ram_role_arn` alicloud-vault will export your long-term credentials to your current shell if you execute `alicloud-vault exec <profilename>`. You can also force this behavior with the flag `--no-session`.

Attributes that need to be specified for StsToken mode are
- name (Required) - the name of your profile. This needs to match with your alicloud-vault profile name.
- mode (Required) - Make sure to specify as `StsToken``.
- ram_role_arn (Required) - the full resource id of the role this profile should assume.  
- ram_session_name (Required) - the name of the role session.
- expired_seconds (Required) - The TTL of the session token. Minium duration 900 seconds, maximum duration 3600 seconds.
- region_id (Optional) - if not specified the default region of any service endpoint that is to be used
- language (Optional) - the language of the CLI. "en" and  "cn" are currently supported.
- site (Optional) - INTL account (intl) or domestic account (domestic)  

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

Alicloud Vault supports the following **environments variables** for profile configuration and temporary credential creation. Note that values defined in ~/.aliyun/config always take precendence.

- `ALICLOUD_REGION` - the region of any service endpoint that is to be used
- `ALICLOUD_ROLE_ARN` - the full resource id of the role this profile should assume
- `ALICLOUD_ROLE_SESSION_NAME` - the name of the role session
- `ALICLOUD_ASSUME_ROLE_TTL` - The TTL of the session token. Minium duration 900 seconds, maximum duration 3600 seconds.

## Development

The [macOS release builds](https://github.com/arafato/alicloud-vault/releases/latest) are currently not code-signed so there will be extra prompts in Keychain.

If you are developing or compiling the alicloud-vault binary yourself, you can [generate a self-signed certificate](https://support.apple.com/en-au/guide/keychain-access/kyca8916/mac) by accessing Keychain Access > Certificate Assistant > Create Certificate > Code Signing Certificate. You can then sign your binary with:

    $ go build .
    $ codesign --sign "Name of my certificate" ./alicloud-vault

## References and Inspiration
- https://github.com/99designs/aws-vault 