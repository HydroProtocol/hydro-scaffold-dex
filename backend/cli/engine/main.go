package main

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/dex_engine"
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"os"
)

func run() int {
	ctx, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	dex_engine.Run(ctx)
	return 0
}

func main() {
	os.Exit(run())
}
