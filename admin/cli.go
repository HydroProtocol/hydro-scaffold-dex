package main

import (
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/urfave/cli"
	"os"
	"sort"
)

func main() {
	process(os.Args)
}

func process(args []string) {
	app := cli.NewApp()
	app.Usage = "to manage dex backend service"
	app.Version = "0.0.1"
	app.Name = "dex api admin"

	admin := NewAdmin(os.Getenv("ADMIN_API_URL"))
	app.Commands = []cli.Command{
		{
			Name:    "new-market",
			Aliases: []string{"nmt"},
			Usage:   "",
			Action: func(c *cli.Context) error {
				data := c.Args().Get(0)
				return admin.NewMarket(data)
			},
		},
		{
			Name:    "edit-market",
			Aliases: []string{"emt"},
			Usage:   "",
			Action: func(c *cli.Context) error {
				data := c.Args().Get(0)
				return admin.EditMarket(data)
			},
		},
		{
			Name:    "cancel-order",
			Aliases: []string{"cor"},
			Usage:   "",
			Action: func(c *cli.Context) error {
				data := c.Args().Get(0)
				return admin.CancelOrder(data)
			},
		},
		{
			Name:    "list-order",
			Aliases: []string{"lor"},
			Usage:   "",
			Action: func(c *cli.Context) error {
				data := c.Args().Get(0)
				return admin.ListAccountOrders(data)
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(args)
	if err != nil {
		utils.Error("error: %v", err)
	}
}
