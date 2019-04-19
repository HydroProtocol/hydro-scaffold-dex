# Run Hydro Box Dex on other ethereum network

## Step0: Cleanup

If you have a dex running, stop it first. Or you can move to the next sectio.

The following command can help you to do that. Be careful, it also clean all data. Backup first if you need the data.

	docker-compose down -v

## Step1: Select network

There are already some docker-compose templates for special envs.
Currently, we support ropsten and mainnet. Choose an network, and copy the template docker-compose file to root directory.For example 

	cp ./envs/ropsten/docker-compose.yaml ./docker-compose-ropsten.yaml

## Step2: Prepare a relayer address

You need to provide an address as relayer on the corresponding network.

Update the docker-compose file you just create. Use values to replace `___CHANGE_ME___` strings.

## Step3: Start the service

Use docker-compose to run.

	docker-compose -f docker-compose-ropsten.yaml up

## Step4: Use

	There is only one market available at first. You can add your own market by using the admin-api and admin-cli. Learn more about them here.

Open `http://localhost:3000/` on your browser. You need to setup an address to trade.	