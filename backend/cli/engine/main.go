package main

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/dex_engine"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
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

	dex_engine.Run(ctx, utils.StartMetrics)
	return 0
}

func main() {
	os.Exit(run())
}
