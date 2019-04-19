# Operate your hydro dex

Admin api is used to change the states of hydro dex. You can use it to add markets, changes fees and so on.The admin api server should not be exposed and should be well protected. An easy way to keep it safe is never start the admin api server on `0.0.0.0`. Always run it on localhost and use admin cli to interative with it.

## Admin cli

### Step 0: Enter the admin docker container

```
docker-compose exec admin sh
```

### Step 1: run commands

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

Get info about an address

```
hydor-dex-ctl address "address" --orders --balances --trades
```

Cancel a special order

```
hydor-dex-ctl order cancel "orderID"
```

Restart engine

```
hydor-dex-ctl engine restart
```

## Admin api

Beside the hydor-dex-ctl command, you can also use rest api to manage your dex.

See details at [here](./backend/admin/api/server.go).
