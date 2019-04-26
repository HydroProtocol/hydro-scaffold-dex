# Run Hydro Box Dex on other ethereum network

we have prepared serveral docker-compose files for different networks, the docker-compose file names are in format like: docker-compose-***eth\_node\_type***[-source].yaml

***eth\_node\_type*** can be:

- localhost - which reprents a localhost eth node
- ropsten - eth ropsten network
- mainnet - eth mainnet network

If file name has surfix: `-source`, the docker images are built from local source code, otherwise, they are pulled from [https://hub.docker.com/u/hydroprotocolio](https://hub.docker.com/u/hydroprotocolio).

The docker-compose files we provide are:

- docker-compose-localhost.yaml
- docker-compose-localhost-source.yaml
- docker-compose-ropsten.yaml
- docker-compose-ropsten-source.yaml
- docker-compose-mainnet.yaml - **TODO**
- docker-compose-mainnet-source.yaml - **TODO**

## Step0: Cleanup

If you have a dex running, stop it first. Or you can move to the next section.

The following command can help you to do that. Be careful, it also clean all data. Backup first if you need them.
	
Remember to use the corresponding config file to down the services.
If you start the services by `docker-compose -f docker-compose-ropsten.yaml up` command,
you should use `docker-compose -f docker-compose-ropsten.yaml down` to stop them.

	docker-compose [-f config_file.yaml] down -v
	

If the following command shows nothing docker containers existent, you are good to go.

	docker ps -a | grep hydro-box-dex

## Step1: Prepare a relayer address

A Hydro Relayer needs to send matching orders to the Hydro Protocol Smart Contracts for settlement. It needs to provide the private key of the relayer address to sign the matching transactions.

For `localhost` environment, we have already prepare all things for you, include the relayer address and private key. Please move to the next section.

For the other environments, you shoud prepare your own relayer address. The values are configured via environment variables. We set these valus default to `___CHANGE_ME___`, you can search this string in the corresponding docker-compose file to locate.

1) Set `HSK_RELAYER_PK` environment variables. The value should be your relayer private key(without `0x` prefix).

2) Set `HSK_RELAYER_ADDRESS` environment variables. The value should be your relayer address(with `0x` prefix, two places in total).

3) Hydro protocol require all relayer address has all quote token approved. It's beacuse when the taker side is `sell`, relayer will be a delegater for quote token between makers and taker. It is designed to allow taker to pay fee without quote approved.

## Step2: Start the service

Use docker-compose to start the service.

```shell
# run HydroBoxDex on ropsten  
docker-compose -f docker-compose-ropsten.yaml up

# or run HydroBoxDex using local source code on ropsten
docker-compose -f docker-compose-ropsten-source.yaml up
```

## Step3: Use

There is only one market available at first. You can add your own market by using the admin-api and admin-cli. Learn more about them [here](./admin-api-and-cli.md).

Open `http://localhost:3000/` on your browser. For localhost, we have already prepared an address for you to test. For the other environments, you need to setup an address to trade.