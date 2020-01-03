# Run on other ethereum network and run from source

we have prepared serveral docker-compose files for different networks, the docker-compose file names are in format like: docker-compose-***eth\_node\_type***[-source].yaml

***eth\_node\_type*** can be:

- localhost - which reprents a localhost eth node
- ropsten - eth ropsten network
- mainnet - eth mainnet network

If file name has suffix: `-source`, the docker images are built from local source code, otherwise, they are pulled from [https://hub.docker.com/u/hydroprotocolio](https://hub.docker.com/u/hydroprotocolio).

The docker-compose files we provide are:

- docker-compose.yaml (default docker compose file for localhost)
- docker-compose-localhost-source.yaml
- docker-compose-ropsten.yaml
- docker-compose-ropsten-source.yaml
- docker-compose-mainnet.yaml
- docker-compose-mainnet-source.yaml

## Step 1: Cleanup

If you have a dex running, stop it first. Or you can move to the next section.

The following command can help you to do that. Be careful, it also clean all data. Backup first if you need them.
	
Remember to use the corresponding config file to down the services.
If you start the services by `docker-compose -f docker-compose-ropsten.yaml up` command,
you should use `docker-compose -f docker-compose-ropsten.yaml down` to stop them.

	docker-compose [-f config_file.yaml] down -v
	

If the following command shows nothing docker containers existent, you are good to go.

	docker ps -a | grep hydro-scaffold-dex

## Step 2: Prepare Your Network

**If you are deploying to the Ethereum Mainnet, Ropsten or localhost skip to Step 3** as these networks have already been prepared.

To use this DEX on a custom network you must first deploy the required [Hydro Protocol v1.1 smart contracts](https://github.com/HydroProtocol/protocol/tree/v1.1). Ensure you are using the correct protocol branch, tagged **v1.1**, as other branches may be incompatible with the current DEX version. This can be done using the provided [deploy script](https://github.com/HydroProtocol/protocol/blob/v1.1/scripts/deploy.js) or following the steps below.

1. Deploy `Proxy.sol`, `TestToken.sol` and `HybridExchange.sol` to your Network. You must pass the addresses of Proxy and TestToken to the HybridExchange's constructor.

2. Update the contract address entries in your docker-compose.yaml with the new addresses.

        HSK_HYBRID_EXCHANGE_ADDRESS=new_HybridExchange_Address
        HSK_PROXY_ADDRESS=new_Proxy_Address
        HSK_HYDRO_TOKEN_ADDRESS=new_TestToken_Address

3. Call `addAddress` on the Proxy smart contract to register the exchange.

        addAddress(new_HybridExchange_Address)
        
   (This is done automatically if you used the provided deploy script)

## Step 3: Prepare a relayer address

A Hydro Relayer needs to send matching orders to the Hydro Protocol Smart Contracts for settlement. It needs to provide the private key of the relayer address to sign the matching transactions.

For `localhost` environment, we have already prepare all things for you, include the relayer address and private key. Please move to the next section.

For the other environments, you shoud prepare your own relayer address. The values are configured via environment variables. We set these valus default to `___CHANGE_ME___`, you can search this string in the corresponding docker-compose file to locate.

1) Set `HSK_RELAYER_PK` environment variables. The value should be your relayer private key(without `0x` prefix).

2) Set `HSK_RELAYER_ADDRESS` environment variables. The value should be your relayer address(with `0x` prefix).

3) Make sure there are some Ether in this relayer address. The relayer is responsible for sending the transaction to the Ethereum network, so relayer should have some Ether to pay gas.

4) Hydro protocol require all relayer address has all quote token approved. It's beacuse when the taker side is `sell`, relayer will be a delegater for quote token between makers and taker. It is designed to allow taker to pay fee without quote approved. You can operate follow this [manual](admin-api-and-cli.md#approve-market-tokens-1) to approve tokens. This also requires some Ether to pay gas.

## Step 4: Start the service

Use docker-compose to start the service.

```shell
# run HydroScaffoldDex on ropsten  
docker-compose -f docker-compose-ropsten.yaml up

# or run HydroScaffoldDex using local source code on ropsten
docker-compose -f docker-compose-ropsten-source.yaml up
```

## Step 5: Use

Open `http://localhost:3000/` on your browser. 

For localhost, we have already prepared an address in the frontend wallet for you to trade. For the other environments, you need to setup an address by yourself.

There are 3 markets available at first: `HOT-WETH`, `HOT-DAI` and `WETH-DAI`. You can add your own market by using the admin-api and admin-cli. Learn more about them [here](./admin-api-and-cli.md).
