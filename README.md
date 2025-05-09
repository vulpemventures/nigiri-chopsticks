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
  - example for Bitcoin:
  ```bash
  $ curl -X POST --data '{"address": "2MsnWskyHaHvcZUHA4gnR3G95EnUmZQjzM8", "amount": 0.02}' http://localhost:3000/faucet
  ```
  - example for Liquid
  ```bash
  $ curl -X POST --data '{"address": "2MsnWskyHaHvcZUHA4gnR3G95EnUmZQjzM8", "asset": "2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3", "amount": 0.02}' http://localhost:3000/faucet
  ```
- `/mint` (only for Liquid chain) if faucet is enabled, to issue an asset and sent all issuance amount to an address
  - example:
  ```bash
  $ curl -X POST --data '{"address": "ert1q90dz89u8eudeswzynl3p2jke564ejc2cnfcwuq", "quantity": 1000, "name": "TokenName", "ticker":"TKN"}' http://localhost:3000/mint
  ```
- `/registry` (only for Liquid chain) if faucet is enabled, to get extra info about one or more assets like `name` and `ticker`
  ```bash
  $ curl -X POST --data '{"assets": ["2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3"]}' http://localhost:3000/registry
  # [{"asset":"2dcf5a8834645654911964ec3602426fd3b9b4017554d3f9c19403e7fc1411d3","contract":{"name":"test","ticker":"TST"},"issuance_txin":{"txid":"a0891447adb288e5a49fa10ede7016788a1b3a175cfb423eb133e45f6cefca84","vin":0},"name":"test","ticker":"TST"
  ```
- all [esplora](https://github.com/blockstream/esplora/blob/master/API.md) HTTP API endpoints

**Note:**  
If mining is enabled, the esplora broadcast endpoint is wrapped so that a block is mined just after the transaction is published to get it confirmed; this is useful when running in regtest network.  
All requests to chopsticks are (optionally) logged using a logger inspired by [negroni](https://github.com/urfave/negroni) package.


## Configuration: Command-Line Flags

### Network Selection
- `--chain` - Choose between bitcoin or liquid networks
- `--wallet-name` - Specify wallet name (default: empty string `""`)

### Server Configuration
- `--addr` - Server listening address (default: localhost:3000)
- `--use-tls` - Enable/disable HTTPS (default: true)
- `--use-logger` - Log all request/response activity

### Bitcoin Node Connection
- `--rpc-addr` - Bitcoin RPC server address (default: localhost:19001)
- `--btc-cookie` - Bitcoin RPC credentials (default: admin1:123)

### Liquid Node Connection
- `--liquid-addr` - Liquid RPC server address (default: localhost:18884)
- `--registry-path` - Path for asset registry database (Liquid only, defaults to current directory)

### Electrum Server Connection
- `--electrs-addr` - Electrs HTTP server address (default: localhost:3002)

### Testing Utilities
- `--use-faucet` - Enable /faucet endpoint for requesting funds
- `--use-mining` - Enable automatic block mining after transaction broadcast
