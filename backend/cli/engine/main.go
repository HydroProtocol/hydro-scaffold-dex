package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"github.com/HydroProtocol/hydro-box-dex/backend/dex_engine"
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
