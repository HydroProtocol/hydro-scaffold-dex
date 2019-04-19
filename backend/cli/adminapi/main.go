package main

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/admin/cli"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"os"
)

func main (){
	app := admincli.NewDexCli()
	err := app.Run(os.Args)
	if err != nil {
		utils.Error("error: %v" ,err)
	}
}
