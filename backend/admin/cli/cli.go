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
	app.Usage = "hydro-dex-cli COMMAND"
	app.Version = "0.0.1"
	app.Name = "dex admin management tool"

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
			Usage:       "",
			Subcommands: []cli.Command{
				{
					Name:        "new",
					Usage:       "",
					Description: "create a market",
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
					Usage:       "",
					Description: "update a market",
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
					Usage:       "",
					Description: "publish a market",
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
					Usage:       "",
					Description: "unpublish a market",
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
					Usage:       "",
					Description: "change market fees",
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
			Usage: "",
			Subcommands: cli.Commands{
				{
					Name:  "orders",
					Usage: "",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(1)

						return admin.ListAccountOrders(address, limit, offset, status)
					},
				},
				{
					Name:  "balances",
					Usage: "",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(1)

						return admin.ListAccountBalances(address, limit, offset)
					},
				},
				{
					Name:  "trades",
					Usage: "",
					Flags: orderListFlags,
					Action: func(c *cli.Context) error {
						address := c.Args().Get(1)

						return admin.ListAccountTrades(address, limit, offset, status)
					},
				},
			},
		},
		{
			Name:  "order",
			Usage: "",
			Subcommands: cli.Commands{
				{
					Name:  "cancel",
					Usage: "",
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
			Usage: "",
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
			Usage: "",
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
