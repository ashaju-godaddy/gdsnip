package main

import (
	"fmt"
	"os"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/commands"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, tui.FormatError(err.Error()))
		os.Exit(1)
	}
}
