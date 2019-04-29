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

	//var limit string
	//var offset string
	//var status string

	newMarketFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "baseTokenAddress",
			Usage:       "Required",
			Destination: &baseTokenAddress,
		},
		cli.StringFlag{
			Name:        "quoteTokenAddress",
			Usage:       "Required",
			Destination: &quoteTokenAddress,
		},
		cli.StringFlag{
			Name:        "minOrderSize",
			Usage:       "Optional",
			Destination: &minOrderSize,
		},
		cli.StringFlag{
			Name:        "pricePrecision",
			Usage:       "Optional",
			Destination: &pricePrecision,
		},
		cli.StringFlag{
			Name:        "priceDecimals",
			Usage:       "Optional",
			Destination: &priceDecimals,
		},
		cli.StringFlag{
			Name:        "amountDecimals",
			Usage:       "Optional",
			Destination: &amountDecimals,
		},
		cli.StringFlag{
			Name:        "makerFeeRate",
			Usage:       "Optional",
			Destination: &makerFeeRate,
		},
		cli.StringFlag{
			Name:        "takerFeeRate",
			Usage:       "Optional",
			Destination: &takerFeeRate,
		},
		cli.StringFlag{
			Name:        "gasUsedEstimation",
			Usage:       "Optional",
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
	//
	//orderListFlags := []cli.Flag{
	//	cli.StringFlag{
	//		Name:        "marketID",
	//		Destination: &marketID,
	//	},
	//	cli.StringFlag{
	//		Name:        "limit",
	//		Destination: &limit,
	//	},
	//	cli.StringFlag{
	//		Name:        "offset",
	//		Destination: &offset,
	//	},
	//	cli.StringFlag{
	//		Name:        "status",
	//		Destination: &status,
	//	},
	//}
	//
	//balanceListFlags := []cli.Flag{
	//	cli.StringFlag{
	//		Name:        "limit",
	//		Destination: &limit,
	//	},
	//	cli.StringFlag{
	//		Name:        "offset",
	//		Destination: &offset,
	//	},
	//}

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
					Name:  "new",
					Usage: "Create a market",
					Description: `
    Example 1): create a market, just set the token addresses, use default parameters for the other attributes.

    hydor-dex-ctl market new HOT-WWW \
  		--baseTokenAddress=0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218 \
  		--quoteTokenAddress=0xbc3524faa62d0763818636d5e400f112279d6cc0

    Example 2): create a market with full attributes

    hydor-dex-ctl market new HOT-WWW \
        --baseTokenAddress=0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218 \
        --quoteTokenAddress=0xbc3524faa62d0763818636d5e400f112279d6cc0 \
        --minOrderSize=0.1 \
        --pricePrecision=5 \
        --priceDecimals=5 \
        --amountDecimals=5 \
        --makerFeeRate=0.001 \
        --takerFeeRate=0.002 \
        --gasUsedEstimation=150000`,
					Flags: newMarketFlags,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)

						if len(marketID) == 0 || len(baseTokenAddress) == 0 || len(quoteTokenAddress) == 0 {
							return cli.ShowSubcommandHelp(c)
						}

						printIfErr(admin.NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation))
						return nil
					},
				},
				{
					Name:  "update",
					Usage: "Update a market",
					Description: `
    Example:
    
    hydor-dex-ctl market update HOT-WWW --amountDecimals=3`,
					Flags: marketUpdateFlags,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							return cli.ShowSubcommandHelp(c)
						}

						printIfErr(admin.UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish))
						return nil
					},
				},
				{
					Name:  "publish",
					Usage: "Publish a market",
					Description: `
    Example: publish market 'HOT-WWW'
    
    hydor-dex-ctl market publish HOT-WWW`,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							return cli.ShowSubcommandHelp(c)
						}

						printIfErr(admin.PublishMarket(marketID))
						return nil
					},
				},
				{
					Name:  "approve",
					Usage: "Approve a market",
					Description: `
    Example: approve market 'HOT-WWW'
    
    hydor-dex-ctl market approve HOT-WWW`,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							return cli.ShowSubcommandHelp(c)
						}

						printIfErr(admin.ApproveMarket(marketID))
						return nil
					},
				},
				{
					Name:  "unpublish",
					Usage: "Unpublish a market",
					Description: `
    Example: unpublish market 'HOT-WWW'

    hydor-dex-ctl market unpublish HOT-WWW`,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						if len(marketID) == 0 {
							return cli.ShowSubcommandHelp(c)
						}

						printIfErr(admin.UnPublishMarket(marketID))
						return nil
					},
				},
				{
					Name:  "changeFees",
					Usage: "Change maker fee and taker fee of a market",
					Description: `
    Example:

    hydor-dex-ctl market changeFees HOT-WETH "0.001" "0.003"`,
					Action: func(c *cli.Context) error {
						marketID = c.Args().Get(0)
						makerFee := c.Args().Get(1)
						takerFee := c.Args().Get(2)

						if len(marketID) == 0 || len(makerFee) == 0 || len(takerFee) == 0 {
							return cli.ShowSubcommandHelp(c)
						}

						printIfErr(admin.UpdateMarketFee(marketID, makerFee, takerFee))
						return nil
					},
				},
			},
		},
		//{
		//	Name:  "address",
		//	Usage: "Get info of an address",
		//	Subcommands: cli.Commands{
		//		{
		//			Name:  "orders",
		//			Usage: "Get address orders",
		//			Flags: orderListFlags,
		//			Action: func(c *cli.Context) error {
		//				address := c.Args().Get(0)
		//				if len(address) == 0 {
		//					fmt.Println("missing arguments, usage: hydro-dex-cli address orders address --marketID=xxx")
		//					return nil
		//				}
		//
		//				if len(marketID) == 0 {
		//					fmt.Println("missing arguments, usage: hydro-dex-cli address orders address --marketID=xxx")
		//					return nil
		//				}
		//
		//				printIfErr(admin.ListAccountOrders(marketID, address, limit, offset, status))
		//				return nil
		//			},
		//		},
		//		{
		//			Name:  "balances",
		//			Usage: "Get address balances",
		//			Flags: balanceListFlags,
		//			Action: func(c *cli.Context) error {
		//				address := c.Args().Get(0)
		//				if len(address) == 0 {
		//					fmt.Println("missing arguments, usage: hydro-dex-cli address balances address")
		//					return nil
		//				}
		//
		//				printIfErr(admin.ListAccountBalances(address, limit, offset))
		//				return nil
		//			},
		//		},
		//		{
		//			Name:  "trades",
		//			Usage: "Get address trades",
		//			Flags: orderListFlags,
		//			Action: func(c *cli.Context) error {
		//				address := c.Args().Get(0)
		//				if len(address) == 0 {
		//					fmt.Println("missing arguments, usage: hydro-dex-cli address trades address --marketID=xxx")
		//					return nil
		//				}
		//
		//				if len(marketID) == 0 {
		//					fmt.Println("missing arguments, usage: hydro-dex-cli address trades address --marketID=xxx")
		//					return nil
		//				}
		//
		//				printIfErr(admin.ListAccountTrades(marketID, address, limit, offset, status))
		//				return nil
		//			},
		//		},
		//	},
		//},
		//{
		//	Name:  "order",
		//	Usage: "Manage order. (cancel)",
		//	Subcommands: cli.Commands{
		//		{
		//			Name:  "cancel",
		//			Usage: "Cancel an order by orderID",
		//			Action: func(c *cli.Context) error {
		//				orderID := c.Args().Get(0)
		//				if len(orderID) == 0 {
		//					fmt.Println("missing arguments, usage: hydro-dex-cli order cancel orderID")
		//					return nil
		//				}
		//
		//				printIfErr(admin.CancelOrder(orderID))
		//				return nil
		//			},
		//		},
		//	},
		//},
		//{
		//	Name:  "engine",
		//	Usage: "Manage hydro dex engine",
		//	Subcommands: cli.Commands{
		//		{
		//			Name:  "restart",
		//			Usage: "",
		//			Flags: orderListFlags,
		//			Action: func(c *cli.Context) error {
		//				printIfErr(admin.RestartEngine())
		//				return nil
		//			},
		//		},
		//	},
		//},
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
