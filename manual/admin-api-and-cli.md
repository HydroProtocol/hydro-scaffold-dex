Hydro DEX's provide an Admin API: a RESTful interface for operating and configuring your DEX.

Hydro also provides a basic CLI to make configuring your DEX simple and easy. This document:

- Summarizes the [key details of the Admin API](https://github.com/HydroProtocol/hydro-scaffold-dex/blob/master/manual/admin-api-and-cli.md#admin-api)
- Provides a [guide for CLI functions](https://github.com/HydroProtocol/hydro-scaffold-dex/blob/master/manual/admin-api-and-cli.md#cli-guide-admin-cli)

*Note that because this API controls important fundamental elements of Hydro dex, it is important to secure this API against unwanted access.*

***

# Configuring Your Hydro Relayer

## Admin API

### Supported Content Types

The Admin API accepts `application/json` types on every endpoint

### Information routes

#### List all markets

```
GET /markets
```

##### Response

```json
{	
	"status": "success",
	"data": [
		{
			"id": "HOT-DAI",
			"baseTokenSymbol": "HOT",
			"BaseTokenName": "HOT",
			"baseTokenAddress": "0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218",
			"baseTokenDecimals": 18,
			"quoteTokenSymbol": "DAI",
			"QuoteTokenName": "DAI",
			"quoteTokenAddress": "0xbc3524faa62d0763818636d5e400f112279d6cc0",
			"quoteTokenDecimals": 18,
			"minOrderSize": "0.001",
			"pricePrecision": 5,
			"priceDecimals": 5,
			"amountDecimals": 5,
			"makerFeeRate": "0.003",
			"takerFeeRate": "0.001",
			"gasUsedEstimation": 1,
			"isPublished": true
		}
	]
}
```

#### Create a market

```
POST /markets
```

##### Request body

```js
{
	"id": "HOT-WETH",                                                  // required
	"baseTokenAddress": "0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218",  // required
	"quoteTokenAddress": "0xbc3524faa62d0763818636d5e400f112279d6cc0", // required
	"minOrderSize": "0.001",                                           // optional default 0.01
	"pricePrecision": 5,                                               // optional default 5
	"priceDecimals": 5,                                                // optional default 5
	"amountDecimals": 5,                                               // optional default 5
	"makerFeeRate": "0.003",                                           // optional default 0.003
	"takerFeeRate": "0.001",                                           // optional default 0.001
	"gasUsedEstimation": 1,                                            // optional default 190000
	"isPublished": true                                                // optional default false
}
```

##### Response on success

```json
{
	"status": "success"
}
```

##### Response on fail

```json
{
	"status": "fail",
	"error_message": "reason"
}
```

#### Approve market tokens

```
POST /markets/approve?marketID=HOT-WETH
```

##### Response on success

```json
{
	"status": "success"
}
```

##### Response on fail

```json
{
	"status": "fail",
	"error_message": "reason"
}
```

#### Update a market

```
PUT /markets
```

##### Request body

```json
{
	"id": "HOT-WETH",
	"minOrderSize": "0.001",
	"isPublished": true
}
```

##### Response on success

```json
{
	"status": "success"
}
```

##### Response on fail

```json
{
	"status": "fail",
	"error_message": "reason"
}
```

***

## CLI Guide (admin-cli)

If you are using docker-compose to run your hydro relayer, you can login into the admin service by entering: 

		docker-compose exec admin sh

This enters the Admin CLI. Once you are logged in, you can use the commands detailed below to configure your DEX. To exit the CLI, type `exit`

### Commands

#### Show help

```
hydro-dex-ctl help
```

#### Get dex status

```
hydro-dex-ctl status
```

#### Manage markets

```
hydro-dex-ctl market help
```

#### Get all markets
```
hydro-dex-ctl market list
```

#### Create a new market

When creating a new market, you can choose to either:

- use default options for the majority of the parameters
- specify all parameters

To use default options, you only need to specify the base and quote token addresses for your trading pair. You can always edit these parameters later.

```
// Default market creation: specify the token addresses for your trading pair
hydro-dex-ctl market new HOT-WWW \
  --baseTokenAddress=0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218 \
  --quoteTokenAddress=0xbc3524faa62d0763818636d5e400f112279d6cc0

// create a new market and specify all attributes
hydro-dex-ctl market new HOT-WWW \
  --baseTokenAddress=0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218 \
  --quoteTokenAddress=0xbc3524faa62d0763818636d5e400f112279d6cc0 \
  --minOrderSize=0.1 \
  --pricePrecision=5 \
  --priceDecimals=5 \
  --amountDecimals=5 \
  --makerFeeRate=0.001 \
  --takerFeeRate=0.002 \
  --gasUsedEstimation=150000

```

#### Approve market tokens.

In order to make trades with a new token pair, the relayer must first set token allowance permissions for the new market. To do this, enter:

```
hydro-dex-ctl market approve HOT-WETH
```

In the example above, the Relayer is approving both HOT and WETH (if not already approved) for trading.

#### Update a market

To change parameters in an existing market, specify the market and the parameter you wish to modify.

```
hydro-dex-ctl market update HOT-WWW --amountDecimals=3
```

#### Publish a market

The market selection area on the frontend will only show markets that are published.

```
hydro-dex-ctl market publish HOT-WETH
```

#### Unpublish a market

Unpublishing a market will make the market no exist in the frontend market selection area.

```
hydro-dex-ctl market unpublish HOT-WETH

```
#### Change fee structure for a market

You can modify the fee structure for each market. Hydro supports assymetric fee structures - so you can have the order maker and order taker automatically pay different fees.

```
// set the HOT-WETH market makerFee to 0.1% and takerFee to 0.3%

hydro-dex-ctl market changeFees HOT-WETH "0.001" "0.003"
```
