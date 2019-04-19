package main

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/admin/api"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"os"
)

func run() int {
	ctx, stop := context.WithCancel(context.Background())

	go cli.WaitExitSignal(stop)
	adminapi.StartServer(ctx)

	return 0
}

func main() {
	os.Exit(run())
}
