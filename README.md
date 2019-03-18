# Nigiri chopsticks

A simple web server written in golang that proxies requests to [*Nigiri*](https://github.com/vulpemventures/nigiri.git) services and expose 2 native endpoints:

* `POST /faucet` faucet endpoint that expects a receiving address in the request body `{"address":<receiving_address}`.
* `POST /broadcast` endpoint that pushes a signed transaction to the network and mines a block to get it confirmed.

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

* `/faucet` includes `/send` and `/broadcast` endpoints
* `/esplora` includes all *electrs* API endpoints

To customize server urls and ports use flags when running the binary:

* `--host` server listening address (default `localhost:3000`)
* `--btc-cookie` btc RPC server user and password (default `admin1:123`)
* `--btc-rpc-addr` btc RPC server listening address. If not specified rpc will be disabled (default `localhost:19001`)
* `--liquid-prc-addr` liquid RPC server listening address. If not specified rpc will be disabled (default `localhost:18884`)
* `--rest-addr` electrs REST HTTP server listening address (default `localhost:3002`)
* `--grpc-addr` electrs gRPC server listening address
* `--use-tls` specify using either `http` or `https` (default `http` - TODO)
**IDEA**
* `--btc-disable-wallet` specify bitcoin.conf disablewallet flag
* `--liquid-disable-wallet` specify liquid.conf disablewallet flag 
