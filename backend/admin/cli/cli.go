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

	marketFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "marketID",
			Destination: &marketID,
		},
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

	app.Commands = []cli.Command{
		{
			Name:        "market",
			Description: "manage markets",
			Usage:       "Manage markets. (create, update)",
			Subcommands: []cli.Command{
				{
					Name:        "new",
					Usage:       "Create a market",
					Description: "Create a market",
					Flags:       marketFlags,
					Action: func(c *cli.Context) error {
						if len(marketID) == 0 || len(baseTokenAddress) == 0 || len(quoteTokenAddress) == 0 {
							return fmt.Errorf("require flag marketID, usage: hydro-dex-cli market new --marketId=xxx --baseTokenAddress=xxx --quoteTokenAddress=xxx")
						}

						return admin.NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation)
					},
				},
				{
					Name:        "update",
					Usage:       "Update a market",
					Description: "Update a market",
					Flags:       marketUpdateFlags,
					Action: func(c *cli.Context) error {
						if len(marketID) == 0 || len(baseTokenAddress) == 0 || len(quoteTokenAddress) == 0 {
							return fmt.Errorf("require flag marketID, usage: hydro-dex-cli market new --marketId=xxx --baseTokenAddress=xxx --quoteTokenAddress=xxx")
						}

						return admin.UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish)
					},
				},
				{
					Name:        "publish",
					Usage:       "Publish a market",
					Description: "Publish a market",
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(1)
						if len(marketID) == 0 {
							return fmt.Errorf("require flag marketID, usage: hydro-dex-cli market publish --marketID=xxx")
						}

						return admin.PublishMarket(marketID)
					},
				},
				{
					Name:        "unpublish",
					Usage:       "Unpublish a market",
					Description: "Unpublish a market",
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(1)
						if len(marketID) == 0 {
							return fmt.Errorf("require flag marketID, usage: hydro-dex-cli market unpublish xxx")
						}

						return admin.UnPublishMarket(marketID)
					},
				},

				{
					Name:        "changeFees",
					Usage:       "Change maker fee and taker fee of a market",
					Description: "Change maker fee and taker fee of a market",
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(1)
						makerFee := c.Args().Get(2)
						takerFee := c.Args().Get(3)

						if len(marketID) == 0 || len(makerFee) == 0 || len(takerFee) == 0 {
							return fmt.Errorf("missing arguments, usage: hydro-dex-cli market marketID makerFee takerFee")
						}

						return admin.UpdateMarketFee(marketID, makerFee, takerFee)
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
						address := c.Args().Get(1)

						return admin.ListAccountOrders(address, limit, offset, status)
					},
				},
				{
					Name:  "balances",
					Usage: "Get address balances",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(1)

						return admin.ListAccountBalances(address, limit, offset)
					},
				},
				{
					Name: "trades",
					Usage: "Get address trades",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(1)

						return admin.ListAccountTrades(address, limit, offset, status)
					},
				},
			},
		},
		{
			Name: "order",
			Usage: "Manage order. (cancel)",
			Subcommands: cli.Commands{
				{
					Name:  "cancel",
					Usage: "Cancel an order by orderID",
					Action: func(c *cli.Context) error {
						orderID := c.Args().Get(1)
						if len(orderID) == 0 {
							return fmt.Errorf("missing arguments, usage: hydro-dex-cli order orderID")
						}

						return admin.CancelOrder(orderID)
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
						option := c.Args().Get(1)
						if len(option) == 0 {
							return fmt.Errorf("missing arguments, usage: hydro-dex-cli order orderID")
						}

						return admin.RestartEngine()
					},
				},
			},
		},
		{
			Name:  "status",
			Usage: "Get current status of the ",
			Action: func(c *cli.Context) error {
				option := c.Args().Get(1)
				if len(option) == 0 {
					return fmt.Errorf("missing arguments, usage: hydro-dex-cli order orderID")
				}

				return admin.Status()
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app

}
