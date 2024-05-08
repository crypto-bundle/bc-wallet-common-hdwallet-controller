# Bc-wallet-common-hdwallet-controller

## Description

Application for control access to mnemonic derivation wallets. Application support flow for:
* Create and manage mnemonic wallets
* Open and manager wallet sessions
* Prepare and close signature request

HdWallet-controller application is first part of hdwallet applications bundle. 
HdWallet-api application is second part of bundle which present in another repository - [bc-wallet-common-hdwallet-api](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet-api)
Third part - target blockchain plugin. For example:

* [bc-wallet-tron-hdwallet](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet)
* [bc-wallet-ethereum-hdwallet](https://github.com/crypto-bundle/bc-wallet-ethereum-hdwallet)
* [bc-wallet-bitcoin-hdwallet](https://github.com/crypto-bundle/bc-wallet-bitcoin-hdwallet)

## Api
Service has two types of API:
* gRPC-API
  * API documentation [here](docs/api/controller_proto.md)
  * Protobuf descriptions [/pkg/proto/controller_api](/pkg/proto/controller_api)
* NatsRPC API

### Integration

Example of gRPC integration with hdwallet-controller application you can see in integration test [/pkg/grpc/controller](/pkg/grpc/controller)

## Infrastructure dependencies

* **PostgreSQL** as main persistent storage
* **Redis** or redis-protocol compatible cache-server as cache storage
* **Nats** - for communication between processing perimeters 
* **Hashicorp Vault** as service provider of secrets management and as provider of encrypt/decrypt sensitive information flow flow
* Instance of **bc-wallet-common-hdwallet-api** for target blockchain

### PostgreSQL

Database: bc-wallet-tron-hdwallet

Users:
* Controller application:
  * Username: bc-wallet-tron-hdwallet-controller
  * Privileges: SELECT, INSERT, UPDATE
* Migrator application:
  * Username: bc-wallet-tron-hdwallet-migrator
  * Privileges: CREATE, DELETE, DROP, SELECT, INSERT, UPDATE
* Terraformer application:
  * Username: bc-wallet-tron-hdwallet-terraformer
  * Privileges: CREATE, DELETE, DROP, SELECT, INSERT, UPDATE

* Migrations:
  * Migration tool: [bc-wallet-common-migrator](https://github.com/crypto-bundle/bc-wallet-common-migrator) 
  * Migrations path: [/migrations](/migrations)

### Redis

Redis used for store values of next entities:
* Mnemonic Wallet's:
  * Key name template: `<STAGE_PREFIX>__<APPLICATION_NAME>__MNEMONIC-WALLETS.<WALLET_UUID>`
  * Example: `DEV__BC-WALLET-TRON-HDWALLET-CONTROLLER__MNEMONIC-WALLETS.f14cf623-b1f3-40cf-8eb2-e18e51be08c`
  * Expiration time: infinite
* Mnemonic wallets sessions:
  * Key name template: `<STAGE_PREFIX>__<APPLICATION_NAME>__MNEMONIC-WALLETS-SESSIONS.<WALLET_SESSION_UUID>`
  * Example: `DEV__BC-WALLET-TRON-HDWALLET-CONTROLLER__MNEMONIC-WALLETS-SESSIONS.f14cf623-b1f3-40cf-8eb2-e18e51be08c`
  * Expiration time: equal to default wallet session expiration time from `DEFAULT_WALLET_UNLOAD_INTERVAL` env variable
* Signature requests:
  * Key name template: `<STAGE_PREFIX>__<APPLICATION_NAME>__SIGN_REQUESTS.<SIGN_REQUEST_UUID>`
  * Example: `DEV__BC-WALLET-TRON-HDWALLET-CONTROLLER__SIGN-REQUESTS.5dc0da6b-90ac-4fb2-8a1b-6f52f9df5c65`
  * Expiration time: will expire with parent wallet session

### Hashicorp Vault

Application required common and personal bucket with some secrets. Also, required personal auth token.

### Tron example
Example for tron blockchain:

Auth token: bc-wallet-tron-hdwallet-controller
```bash
vault token create -display-name bc-wallet-tron-hdwallet-controller 
```

Buckets:
* crypto-bundle/bc-wallet-common/transit
  * VAULT_COMMON_TRANSIT_KEY
* crypto-bundle/bc-wallet-tron-hdwallet/common
  * POSTGRESQL_DATABASE_NAME
  * VAULT_APP_ENCRYPTION_KEY
* crypto-bundle/bc-wallet-tron-hdwallet/controller
  * NATS_PASSWORD
  * NATS_USER
  * POSTGRESQL_PASSWORD
  * POSTGRESQL_USERNAME
  * REDIS_PASSWORD
  * REDIS_USER
* crypto-bundle/bc-wallet-tron-hdwallet/migrator
  * POSTGRESQL_PASSWORD
  * POSTGRESQL_USERNAME
* crypto-bundle/bc-wallet-tron-hdwallet/terraformer
  * NATS_PASSWORD
  * NATS_USER
  * POSTGRESQL_PASSWORD
  * POSTGRESQL_USERNAME
  * REDIS_PASSWORD
  * REDIS_USER

Also, application required two encryption keys:
* Common for whole crypto-bundle project transit key - **crypto-bundle-bc-wallet-common-transit-key**. 
Value with transit key name will be loaded from `VAULT_COMMON_TRANSIT_KEY` Vault variable, which stored in 
common bucket - **crypto-bundle/bc-wallet-common/transit**.

* Target encryption key for hdwallet-controller and hdwallet-api - **crypto-bundle-bc-wallet-tron-hdwallet**
Value with transit key name will be loaded from `VAULT_APP_ENCRYPTION_KEY` Vault variable, which stored in
common bucket - **crypto-bundle/bc-wallet-tron-hdwallet/common**.

### Bc-wallet-common-hdwallet-api

Repository of hdwallet-api - [bc-wallet-common-hdwallet-api](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet-api)

Application must work in pair with instance bc-wallet-common-hdwallet-api.
For example in case of Tron blockchain:
Instance of bc-wallet-common-hdwallet-controller - **bc-wallet-tron-hdwallet-controller** must work instance 
of bc-wallet-common-hdwallet-api - **bc-wallet-tron-hdwallet-api**. 

Communication between blockchain via gRPC unix-file socket connection. 
You can control socket file path by change `HDWALLET_UNIX_SOCKET_PATH` environment variable.
Also, you can set target blockchain of hdwallet-controller and hdwallet-api via `PROCESSING_NETWORK` env variable.

## Environment variables

Full example of env variables you can see in  [env-controller-example.env](./env-controller-example.env) file.

## Deployment

Currently, support only kubernetes deployment flow via Helm

### Kubernetes
Application must be deployed as part of bc-wallet-<BLOCKCHAIN_NAME>-hdwallet bundle. 
Application must be started as single container in Kubernetes Pod with shared volume. 

You can see example of HELM-chart deployment application in next repositories:
* [bc-wallet-tron-hdwallet-api/deploy/helm/hdwallet](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet/tree/develop/deploy/helm/hdwallet)
* [bc-wallet-ethereum-hdwallet-api/deploy/helm/hdwallet](https://github.com/crypto-bundle/bc-wallet-ethereum-hdwallet/tree/develop/deploy/helm/hdwallet)

## Licence

**bc-wallet-common-hdwallet-controller** is licensed under the [MIT](./LICENSE) License.