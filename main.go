package main

import (
	"os"

	"gitlab.com/tokend/traefik-cop/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
