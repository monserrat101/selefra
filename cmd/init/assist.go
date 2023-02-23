package init

import (
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/cmd/version"
)

func initHeaderOutput(providers []string) {
	for i := range providers {
		cli_ui.Successln(providers[i] + " [✔]\n")
	}
	cli_ui.Successf(`Running with selefra-cli %s

	This command will walk you through creating a new Selefra project

	Enter a value or leave blank to accept the (default), and press <ENTER>.
	Press ^C at any time to quit.`, version.Version)
}
