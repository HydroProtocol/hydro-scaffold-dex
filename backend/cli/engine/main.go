package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/cli"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/dex_engine"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"os"
)

func run() int {
	ctx, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	dex_engine.Run(ctx, utils.StartMetrics)
	return 0
}

func main() {
	os.Exit(run())
}
