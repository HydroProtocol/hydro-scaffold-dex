package main

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/admin/cli"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"os"
)

func main() {
	cli := admincli.NewDexCli()
	err := cli.Run(os.Args)
	if err != nil {
		utils.Error(err.Error())
	}
}
