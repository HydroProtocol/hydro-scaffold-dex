package main

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"os"
)

func run() int {
	_, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	//todo
	//launcher.Run(ctx)
	return 0
}

func main() {
	os.Exit(run())
}
