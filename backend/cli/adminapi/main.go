package main

import (
	"context"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/admin/api"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/cli"
	_ "github.com/joho/godotenv/autoload"
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
