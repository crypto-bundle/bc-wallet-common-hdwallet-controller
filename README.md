# bc-wallet-common-hdwallet-controller

## Description

Service of management owned hd-wallets - mnemonic wallets for Tron blockchain. 
Application purpose - create new mnemonic wallet, sign transaction, get derivation address

## Api

Service has a gRPC-api - api doc [here](docs/api/controller_proto.md)

## Dependencies
* PostgreSQL
* Redis or redis-protocol compatible cache-server
* Nats
* Instance of bc-wallet-common-hdwallet-api for target blockchain

### PostgreSQL

### Redis

### Nats

### Bc-wallet-common-hdwallet-api

## Environment variables
Example of env variables - [env-controller-example.env](./env-controller-example.env)

## Deployment

You can see example of HELM-chart deployment application in next repositories:
* [bc-wallet-tron-hdwallet-api/deploy/helm/hdwallet](https://github.com/crypto-bundle/bc-wallet-tron-hdwallet/tree/develop/deploy/helm/hdwallet)
* [bc-wallet-ethereum-hdwallet-api/deploy/helm/hdwallet](https://github.com/crypto-bundle/bc-wallet-ethereum-hdwallet/tree/develop/deploy/helm/hdwallet)


## Licence

**bc-wallet-common-hdwallet-controller** is licensed under the [MIT](./LICENSE) License.