const Web3 = require("web3");
const Proxy = artifacts.require("./Proxy.sol");
const HybridExchange = artifacts.require("./HybridExchange.sol");
const TestToken = artifacts.require("./helper/TestToken.sol");
const WethToken = artifacts.require("./helper/WethToken.sol");
const BigNumber = require("bignumber.js");

BigNumber.config({ EXPONENTIAL_AT: 1000 });

const getWeb3 = () => {
  const myWeb3 = new Web3(web3.currentProvider);
  return myWeb3;
};

const newContract = async (contract, ...args) => {
  const c = await contract.new(...args);
  const w = getWeb3();
  const instance = new w.eth.Contract(contract.abi, c.address);
  return instance;
};

const newContractAt = (contract, address) => {
  const w = getWeb3();
  const instance = new w.eth.Contract(contract.abi, address);
  return instance;
};

module.exports = async () => {
  let hot, exchange, proxy;
  try {
    const testAddresses = web3.eth.accounts.slice(1, 6);
    const owner = web3.eth.accounts[0];
    const relayer = web3.eth.accounts[9];
    const maker = web3.eth.accounts[8];

    console.log("owner", owner);
    console.log("relayer", relayer);
    console.log("maker", maker);
    console.log("testAddresses", testAddresses);

    const bigAllowance =
      "0xf000000000000000000000000000000000000000000000000000000000000000";

    hot = await newContract(TestToken, "HydroToken", "Hot", 18);
    console.log("Hydro Token address", web3.toChecksumAddress(hot._address));

    proxy = await newContract(Proxy);
    console.log("Proxy address", web3.toChecksumAddress(proxy._address));

    exchange = await newContract(HybridExchange, proxy._address, hot._address);
    console.log(
      "HybridExchange address",
      web3.toChecksumAddress(exchange._address)
    );

    await Proxy.at(proxy._address).addAddress(exchange._address);
    console.log("Proxy add exchange into whitelist");

    usd = await newContract(TestToken, "USD TOKEN", "USD", 18);
    console.log("USD TOKEN address", web3.toChecksumAddress(usd._address));

    weth = await newContract(WethToken, "Wrapped Ethereum", "WETH", 18);
    console.log(
      "Wrapped Ethereum TOKEN address",
      web3.toChecksumAddress(weth._address)
    );

    const approveAllToken = async (address, tokens) => {
      for (let i = 0; i < tokens.length; i++) {
        const token = tokens[i];
        const tokenName = await token.methods.name().call();

        await token.methods
          .approve(proxy._address, bigAllowance)
          .send({ from: address });
        console.log(address, `${tokenName} approved`);
      }
    };

    const giveCoinsTo = async (address, amount, tokens) => {
      for (let i = 0; i < tokens.length; i++) {
        const token = tokens[i];
        const tokenName = await token.methods.name().call();

        await token.methods
          .transfer(address, `${amount}000000000000000000`)
          .send({ from: owner });
        console.log(address, `${amount} ${tokenName} received`);
      }
    };

    const wrapETH = async (address, amount) => {
      await weth.methods
        .deposit()
        .send({ from: address, value: `${amount}000000000000000000` });
      console.log(address, `${amount} WETH deposited`);
    };

    await Promise.all(
      testAddresses.map(async u => {
        await approveAllToken(u, [usd, hot, weth]);
        await giveCoinsTo(u, "100000", [hot, usd]);
        await wrapETH(u, "1000");
      })
    );

    await approveAllToken(maker, [usd, hot, weth]);
    await giveCoinsTo(maker, "100000", [hot, usd]);
    await wrapETH(maker, "1000");
    await approveAllToken(relayer, [usd, hot, weth]);

    process.exit(0);
  } catch (e) {
    console.log(e);
  }
};
