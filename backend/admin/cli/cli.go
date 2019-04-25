package admincli

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"sort"
)

func NewDexCli() *cli.App {
	admin := NewAdmin(os.Getenv("ADMIN_API_URL"), nil, nil)

	app := cli.NewApp()
	app.Usage = "A tool to manage hydro dex"
	app.Version = "0.0.1"
	app.Name = "hydro-dex-ctl"

	var marketID string
	var baseTokenAddress string
	var quoteTokenAddress string
	var isPublish string

	var minOrderSize string
	var pricePrecision string
	var priceDecimals string
	var amountDecimals string
	var makerFeeRate string
	var takerFeeRate string
	var gasUsedEstimation string

	var limit string
	var offset string
	var status string

	newMarketFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "baseTokenAddress",
			Destination: &baseTokenAddress,
		},
		cli.StringFlag{
			Name:        "quoteTokenAddress",
			Destination: &quoteTokenAddress,
		},
		cli.StringFlag{
			Name:        "minOrderSize",
			Destination: &minOrderSize,
		},
		cli.StringFlag{
			Name:        "pricePrecision",
			Destination: &pricePrecision,
		},
		cli.StringFlag{
			Name:        "priceDecimals",
			Destination: &priceDecimals,
		},
		cli.StringFlag{
			Name:        "amountDecimals",
			Destination: &amountDecimals,
		},
		cli.StringFlag{
			Name:        "makerFeeRate",
			Destination: &makerFeeRate,
		},
		cli.StringFlag{
			Name:        "takerFeeRate",
			Destination: &takerFeeRate,
		},
		cli.StringFlag{
			Name:        "gasUsedEstimation",
			Destination: &gasUsedEstimation,
		},
	}

	marketUpdateFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "minOrderSize",
			Destination: &minOrderSize,
		},
		cli.StringFlag{
			Name:        "pricePrecision",
			Destination: &pricePrecision,
		},
		cli.StringFlag{
			Name:        "priceDecimals",
			Destination: &priceDecimals,
		},
		cli.StringFlag{
			Name:        "amountDecimals",
			Destination: &amountDecimals,
		},
		cli.StringFlag{
			Name:        "makerFeeRate",
			Destination: &makerFeeRate,
		},
		cli.StringFlag{
			Name:        "takerFeeRate",
			Destination: &takerFeeRate,
		},
		cli.StringFlag{
			Name:        "gasUsedEstimation",
			Destination: &gasUsedEstimation,
		},
		cli.StringFlag{
			Name:        "isPublish",
			Destination: &isPublish,
		},
	}

	orderListFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "marketID",
			Destination: &marketID,
		},
		cli.StringFlag{
			Name:        "limit",
			Destination: &limit,
		},
		cli.StringFlag{
			Name:        "offset",
			Destination: &offset,
		},
		cli.StringFlag{
			Name:        "status",
			Destination: &status,
		},
	}

	balanceListFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "limit",
			Destination: &limit,
		},
		cli.StringFlag{
			Name:        "offset",
			Destination: &offset,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "market",
			Description: "manage markets",
			Usage:       "Manage markets. (create, update)",
			Subcommands: []cli.Command{
				{
					Name:        "list",
					Usage:       "List markets",
					Description: "List markets",
					Action: func(c *cli.Context) error {
						printIfErr(admin.ListMarkets())
						return nil
					},
				},
				{
					Name:        "new",
					Usage:       "Create a market",
					Description: "Create a market",
					Flags:       newMarketFlags,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)

						if len(marketID) == 0 || len(baseTokenAddress) == 0 || len(quoteTokenAddress) == 0 {
							fmt.Println("require flag marketID, usage: hydro-dex-cli market new marketId --baseTokenAddress=xxx --quoteTokenAddress=xxx")
							return nil
						}

						printIfErr(admin.NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation))
						return nil
					},
				},
				{
					Name:        "update",
					Usage:       "Update a market",
					Description: "Update a market",
					Flags:       marketUpdateFlags,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							fmt.Println("require flag marketID, usage: hydro-dex-cli market update marketID [flags]")
							return nil
						}

						printIfErr(admin.UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish))
						return nil
					},
				},
				{
					Name:        "publish",
					Usage:       "Publish a market",
					Description: "Publish a market",
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							fmt.Println("require flag marketID, usage: hydro-dex-cli market publish marketID")
							return nil
						}

						printIfErr(admin.PublishMarket(marketID))
						return nil
					},
				},
				{
					Name:        "unpublish",
					Usage:       "Unpublish a market",
					Description: "Unpublish a market",
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							fmt.Println("require flag marketID, usage: hydro-dex-cli market unpublish marketID")
							return nil
						}

						printIfErr(admin.UnPublishMarket(marketID))
						return nil
					},
				},
				{
					Name:        "changeFees",
					Usage:       "Change maker fee and taker fee of a market",
					Description: "Change maker fee and taker fee of a market",
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						makerFee := c.Args().Get(1)
						takerFee := c.Args().Get(2)

						if len(marketID) == 0 || len(makerFee) == 0 || len(takerFee) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli market changeFees marketID makerFee takerFee")
							return nil
						}

						printIfErr(admin.UpdateMarketFee(marketID, makerFee, takerFee))
						return nil
					},
				},
			},
		},
		{
			Name:  "address",
			Usage: "Get info of an address",
			Subcommands: cli.Commands{
				{
					Name:  "orders",
					Usage: "Get address orders",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(0)
						if len(address) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli address orders address")
							return nil
						}

						if len(marketID) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli address balances address --marketID=xxx")
							return nil
						}

						printIfErr(admin.ListAccountOrders(marketID, address, limit, offset, status))
						return nil
					},
				},
				{
					Name:  "balances",
					Usage: "Get address balances",
					Flags: balanceListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(0)
						if len(address) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli address balances address")
							return nil
						}

						printIfErr(admin.ListAccountBalances(address, limit, offset))
						return nil
					},
				},
				{
					Name:  "trades",
					Usage: "Get address trades",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(0)
						if len(address) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli address trades address")
							return nil
						}

						if len(marketID) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli address trades address --marketID=xxx")
							return nil
						}

						printIfErr(admin.ListAccountTrades(marketID, address, limit, offset, status))
						return nil
					},
				},
			},
		},
		{
			Name:  "order",
			Usage: "Manage order. (cancel)",
			Subcommands: cli.Commands{
				{
					Name:  "cancel",
					Usage: "Cancel an order by orderID",
					Action: func(c *cli.Context) error {
						orderID := c.Args().Get(0)
						if len(orderID) == 0 {
							fmt.Println("missing arguments, usage: hydro-dex-cli order cancel orderID")
							return nil
						}

						printIfErr(admin.CancelOrder(orderID))
						return nil
					},
				},
			},
		},
		{
			Name:  "engine",
			Usage: "Manage hydro dex engine",
			Subcommands: cli.Commands{
				{
					Name:  "restart",
					Usage: "",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						printIfErr(admin.RestartEngine())
						return nil
					},
				},
			},
		},
		{
			Name:  "status",
			Usage: "Get current status of the ",
			Action: func(c *cli.Context) error {
				printIfErr(admin.Status())
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app

}

func printIfErr(ret []byte, err error) {
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(ret))
	}
}
