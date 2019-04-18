package main

import (
	"testing"
)

func TestCli(t *testing.T) {
	app := newDexCli()
	app.Run([]string{"hydro-dex-cli", "market"})
	app.Run([]string{"hydro-dex-cli", "address"})
	app.Run([]string{"hydro-dex-cli", "engine"})
	app.Run([]string{"hydro-dex-cli", "status"})
}
