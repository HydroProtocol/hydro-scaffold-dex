# Operate your hydro dex

After experiencing the basic box feature, you may have a question: How to add a trading market? How to change the fee rates or a market? Fortunately, we have prepared a suite of api and a command line tool for you to modify the state of hydro dex.

Hydor dex Admin API provides a RESTful interface for administration and configuration of markets.

Because this API describes a control of Hydro dex, it is important to secure this API against unwanted access. 

## Admin api

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

```json
{
	"id": "HOT-WETH",
	"baseTokenSymbol": "HOT",
	"BaseTokenName": "HOT",
	"baseTokenAddress": "0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218",
	"baseTokenDecimals": 18,
	"quoteTokenSymbol": "WETH",
	"QuoteTokenName": "Wrapped Ethereum",
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


## hytdro-dex-ctl (admin-cli)

If you are using docker-compose to run hydro dex. You can login into the admin service via `docker-compose exec admin sh`, the `hydro-dex-ctl` binary has already included in it.

### Commands

Show help

```
hydro-dex-ctl help
```

Get dex status.

```
hydor-dex-ctl status
```

Manage markets

```
hydor-dex-ctl market help

hydor-dex-ctl market list

hydor-dex-ctl market new $marketID --baseTokenAddress=”” --quoteTokenAddress=””

hydor-dex-ctl market update $marketID --editableAttr=value

hydor-dex-ctl market publish $marketID

hydor-dex-ctl market unpublish $marketID

hydor-dex-ctl market changeFees $marketID "asMakerFee" "asTakerFee"
```