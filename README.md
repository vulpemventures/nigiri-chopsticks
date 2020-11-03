# Nigiri chopsticks

This is an API passthrough that simply proxies requests to the underlying services.
It expects an electrum REST server and an optional RPC server for faucet and custom broadcasting services.

## Usage

Clone the repo:

```bash
$ git clone git@github.com:vulpmeventures/nigiri-chopsticks.git
```

Run tests:

```
$ bash scripts/test local
```

To run tests locally you must have a running Nigiri instance.

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

- `/faucet` if faucet is enabled, to send funds to an address
  - example:
  ```
  $ curl -X POST --data '{"address": "2MsnWskyHaHvcZUHA4gnR3G95EnUmZQjzM8"}' http://localhost:3000/faucet
  # 142a3ef44ddb3f9f395bcd87d73d71ccc8d90566a219f8d57ff0a71822617413
  ```
- `/mint` (only for Liquid chain) if faucet is enabled, to issue an asset and sent all issuance amount to an address
  - example:
  ```
  $ curl -X POST --data '{"address": "ert1q90dz89u8eudeswzynl3p2jke564ejc2cnfcwuq", "quantity": 1000, "name": "TokenName", "ticker":"TKN"}' http://localhost:3000/mint
  # {"asset":"2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3","txId":"7aed7d7f6b4193875e28036728fd360785324f85dfd84d2951cc2b18ea6c2718"}
  ```
- `/registry` (only for Liquid chain) if faucet is enabled, to get extra info about one or more assets like `name` and `ticker`
  ```
  $ curl -X POST --data '{"assets": ["2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3"]}' http://localhost:3000/registry
  # [{"asset":"2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3","contract":{"name":"test","ticker":"TST"},"issuance_txin":{"txid":"a0891447adb288e5a49fa10ede7016788a1b3a175cfb423eb133e45f6cefca84","vin":0},"name":"test","ticker":"TST"
  ```
- all [esplora](https://github.com/blockstream/esplora/blob/master/API.md) HTTP API endpoints

**Note:**  
If mining is enabled, the esplora broadcast endpoint is wrapped so that a block is mined just after the transaction is published to get it confirmed; this is useful when running in regtest network.  
All requests to chopsticks are (optionally) logged using a logger inspired by [negroni](https://github.com/urfave/negroni) package.

To customize server urls and ports use flags when running the binary:

- `--addr` server listening address (default `localhost:3000`)
- `--btc-addr` btc RPC server listening address (default `localhost:19001`)
- `--btc-cookie` btc RPC server user and password (default `admin1:123`)
- `--liquid-addr` liquid RPC server listening address (default `localhost:18884`)
- `--electrs-addr` electrs HTTP server listening address (default `localhost:3002`)
- `--use-tls` specify using either `http` or `https` (default `true`)
- `--use-faucet` to have a /faucet endpoint available for sending funds
- `--use-mining` to have the esplora /broadcast endpoint wrapped so that a block is mined after the transaction
  is published
- `--use-logger` to log every request/response
- `--registry-path` to set the path for the asset registry db (default to current directory - Liquid only)
