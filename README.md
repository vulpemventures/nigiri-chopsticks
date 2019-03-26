# Nigiri chopsticks

This is an API passthrough that simply proxies requests to the underlying services.
It expects an electrum REST server and an optional RPC server for faucet and custom broadcasting services.

## Usage

Clone the repo:

```bash
$ git clone git@github.com:vulpmeventures/nigiri-chopsticks.git
```

Enter the folder project and install:

```bash
nigiri-chopsticks $ bash scripts/install
```

Build for Linux x64:

```bash
nigiri-chopsticks $ bash scripts/build linux amd64
```

Build for Mac:

```bash
nigiri-chopsticks $ bash scripts/build darwin amd64
```

Run:

```bash
nigiri-chopsticks $ ./build/nigiri-chopsticks-darwin-amd64
```

## Routes and Customization

The web server starts at default address `localhost:3000` with the following routes:

* `/faucet` if faucet is enabled, to send funds to an address
* all [esplora](https://github.com/blockstream/esplora/blob/master/API.md) HTTP API endpoints

**Note:**  
If mining is enabled, the esplora broadcast endpoint is wrapped so that a block is mined just after the transaction is published to get it confirmed; this is useful when running in regtest network.  
All requests to chopsticks are (optionally) logged using a logger inspired by [negroni](https://github.com/urfave/negroni) package.

To customize server urls and ports use flags when running the binary:

* `--addr` server listening address (default `localhost:3000`)
* `--btc-addr` btc RPC server listening address (default `localhost:19001`)
* `--btc-cookie` btc RPC server user and password (default `admin1:123`)
* `--liquid-addr` liquid RPC server listening address (default `localhost:18884`)
* `--electrs-addr` electrs HTTP server listening address (default `localhost:3002`)
* `--use-tls` specify using either `http` or `https` (default `true`)
* `--use-faucet` to have a /faucet endpoint available for sending funds
* `--use-mining` to have the esplora /broadcast endpoint wrapped so that a block is mined after the transaction 
is published
* `--use-logger` to log every request/response 
