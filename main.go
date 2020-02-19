package main

import (
	"os"

	"github.com/arafato/alicloud-vault/cli"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version is provided at compile time
var Version = "dev"

func main() {
	run(os.Args[1:], os.Exit)
}

func run(args []string, exit func(int)) {
	app := kingpin.New(
		`alicloud-vault`,
		`A vault for securely storing and accessing Alibaba Cloud credentials in development environments.`,
	)

	app.ErrorWriter(os.Stderr)
	app.Writer(os.Stdout)
	app.Version(Version)
	app.Terminate(exit)

	cli.ConfigureGlobals(app)

	kingpin.MustParse(app.Parse(args))
}
