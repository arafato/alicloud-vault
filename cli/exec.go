package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/99designs/keyring"
	"github.com/arafato/alicloud-vault/vault"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ExecCommandInput struct {
	ProfileName     string
	Command         string
	Args            []string
	Keyring         keyring.Keyring
	Config          vault.Config
	SessionDuration int
	noSession       bool
}

func ConfigureExecCommand(app *kingpin.Application) {
	input := ExecCommandInput{}

	cmd := app.Command("exec", "Executes a command with AWS credentials in the environment")

	cmd.Flag("duration", "Duration of the temporary or assume-role session. Defaults to 1h").
		Short('d').
		IntVar(&input.SessionDuration)

	cmd.Flag("no-session", "Skip creating STS session with AssumeRole and use long-term credentials").
		Short('n').
		BoolVar(&input.noSession)

	cmd.Arg("profile", "Name of the profile").
		Required().
		HintAction(getProfileNames).
		StringVar(&input.ProfileName)

	cmd.Arg("cmd", "Command to execute, defaults to $SHELL").
		Default(os.Getenv("SHELL")).
		StringVar(&input.Command)

	cmd.Arg("args", "Command arguments").
		StringsVar(&input.Args)

	cmd.Action(func(c *kingpin.ParseContext) error {
		input.Keyring = keyringImpl
		input.Config.AssumeRoleDuration = input.SessionDuration
		app.FatalIfError(ExecCommand(input), "exec")
		return nil
	})
}

func ExecCommand(input ExecCommandInput) error {

	configLoader.BaseConfig = input.Config
	config, err := configLoader.LoadProfile(input.ProfileName)
	if err != nil {
		return err
	}

	credKeyring := &vault.CredentialKeyring{Keyring: input.Keyring}

	if input.noSession {
		config.RoleARN = ""
	}

	creds, err := vault.GenerateTempCredentials(config, credKeyring)
	if err != nil {
		return fmt.Errorf("Error getting temporary credentials for %s: %w", input.ProfileName, err)
	}

	env := environ(os.Environ())

	env.Unset("ALICLOUD_ACCESS_KEY")
	env.Unset("ALICLOUD_SECRET_KEY")

	if config.Region != "" {
		log.Printf("Setting subprocess env: ALICLOUD_REGION=%s", config.Region)
		env.Set("ALICLOUD_REGION", config.Region)
	}

	log.Println("Setting subprocess env: ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY")
	env.Set("ALICLOUD_ACCESS_KEY", creds.Creds.AccessKeyID)
	env.Set("ALICLOUD_SECRET_KEY", creds.Creds.SecretAccessKey)
	env.Set("ALICLOUD_STS_TOKEN", creds.StsToken)
	env.Set("ALICLOUD_SESSION_EXPIRATION", creds.Duration)

	// Manually expanding all ALICLOUD_* environments variables that are enclosed with ''.
	// Environment variables that are not enclosed with '' expand to empty strings if they have not been defined yet in the current shell.
	input.Args = expandEnvironmentVariables(input.Args, env)

	err = execSyscall(input.Command, input.Args, env)

	if err != nil {
		return fmt.Errorf("Error execing process: %w", err)
	}

	return nil
}

// environ is a slice of strings representing the environment, in the form "key=value".
type environ []string

// Unset an environment variable by key
func (e *environ) Unset(key string) {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			(*e)[i] = (*e)[len(*e)-1]
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

func expandEnvironmentVariables(args, env []string) []string {
	var expanded []string
	for _, item := range args {
		if strings.HasPrefix(item[1:], "ALICLOUD_") {
			expanded = append(expanded, getValue(env, item[1:]))
		} else {
			expanded = append(expanded, item)
		}
	}
	return expanded
}

func getValue(env []string, key string) string {
	for _, item := range env {
		if strings.HasPrefix(item, key) {
			return strings.SplitN(item, "=", 2)[1]
		}
	}
	return ""
}

// Set adds an environment variable, replacing any existing ones of the same key
func (e *environ) Set(key, val string) {
	e.Unset(key)
	*e = append(*e, key+"="+val)
}

func execCmd(command string, args []string, env []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start command: %v", err)
	}

	go func() {
		for {
			sig := <-sigChan
			cmd.Process.Signal(sig)
		}
	}()

	if err := cmd.Wait(); err != nil {
		cmd.Process.Signal(os.Kill)
		return fmt.Errorf("Failed to wait for command termination: %v", err)
	}

	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	os.Exit(waitStatus.ExitStatus())
	return nil
}

func supportsExecSyscall() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd"
}

func execSyscall(command string, args []string, env []string) error {
	if !supportsExecSyscall() {
		return execCmd(command, args, env)
	}

	argv0, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	argv := make([]string, 0, 1+len(args))
	argv = append(argv, command)
	argv = append(argv, args...)

	return syscall.Exec(argv0, argv, env)
}
