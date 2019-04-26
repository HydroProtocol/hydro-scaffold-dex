# Run on other ethereum network and run from source

we have prepared serveral docker-compose files for different networks, the docker-compose file names are in format like: docker-compose-*eth_node_type*[-source].yaml

*eth_node_type* can be:

- localhost, which reprents a localhost Ethereum node
- ropsten, Ethereum ropsten network
- mainnet, Ethereum mainnet network

if file name has suffix: `-source`, the docker images are built from local source code, otherwise, they are pulled from https://hub.docker.com/u/hydroprotocolio .

the docker-compose files we provide are:

- docker-compose-localhost.yaml
- docker-compose-localhost-source.yaml
- docker-compose-ropsten.yaml
- docker-compose-ropsten-source.yaml
- docker-compose-mainnet.yaml
- docker-compose-mainnet-source.yaml

## Step0: Cleanup

If you have a dex running, stop it first. Or you can move to the next sectio.

The following command can help you to do that. Be careful, it also clean all data. Backup first if you need the data.

	docker-compose down -v

## Step1: Prepare a relayer address

You need to update relayer address and privateKey for corresponding network.

Replace `___CHANGE_ME___` strings in the docker-file that you wanna run.

## Step2: Start the service

Use docker-compose to start the service.

```shell
# run HydroBoxDex on ropsten  
docker-compose -f docker-compose-ropsten.yaml up

# or run HydroBoxDex using local source code on ropsten  
docker-compose -f docker-compose-ropsten-source.yaml up
```

## Step3: Use

Open `http://localhost:3000/` on your browser. You need to setup an address to trade.	

There are 3 markets available at first: `HOT-WETH`, `HOT-DAI` and `WETH-DAI`. You can add your own market by using the admin-api and admin-cli. Learn more about them [here](./admin-api-and-cli.md).