package main

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/admin/cli"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func main() {
	cli := admincli.NewDexCli()
	err := cli.Run(os.Args)
	if err != nil {
		utils.Errorf(err.Error())
	}
}
