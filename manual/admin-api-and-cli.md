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

#### Show help

```
hydro-dex-ctl help
```

#### Get dex status.

```
hydor-dex-ctl status
```

#### Manage markets

```
hydor-dex-ctl market help
```

#### Get all markets
```
hydor-dex-ctl market list
```

#### Create a new market

```
// create a market, just set the token addresses, use default parameters for the other attributes.
hydor-dex-ctl market new HOT-WWW \
  --baseTokenAddress=0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218 \
  --quoteTokenAddress=0xbc3524faa62d0763818636d5e400f112279d6cc0

// create a market with full attributes
hydor-dex-ctl market new HOT-WWW \
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

#### Update a market

```
hydor-dex-ctl market update HOT-WWW --amountDecimals=3
```

#### Set a market pulbished. 

The market will exist in frontend markets select.

```
hydor-dex-ctl market publish HOT-WETH
```

#### Set a market unpulbished. 

The market will no longer exist in frontend markets select.

```
hydor-dex-ctl market unpublish HOT-WETH

```
#### Change fees of a market.

```
// set HOT-WETH market makerFee to 0.1%, takerFee to 0.3%

hydor-dex-ctl market changeFees HOT-WETH "0.001" "0.003"
```