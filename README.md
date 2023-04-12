# bc-wallet-tron-hdwallet

Service of management owned hd-wallets - mnemonic wallets for Tron blockchain. 
Application purpose - create new mnemonic wallet, sign transaction, get derivation address

## Api

Service has a GRPC-api - api doc [here](./docs/hdwallet_api/hdwallet_proto.md)

## Mnemmonic wallets

Mnemonic wallets stored in Postgresql database as encrypted Hashicord **Vault** data

## Deployment

### K8s

#### Secrets

```
kubectl create secret generic bc-wallet-tron-hdwallet \
  --from-literal=vault_api_user_token='<insert_token_here>' \
  --from-literal=vault_migrator_user_token='<insert_token_here>' \
  --from-literal=vault_updater_user_token='<insert_token_here>' \
  --from-literal=vault_transit_secret_key='bc-wallet-tron-hdwallet'
```

### DB

In production environment application require 3 database users, for example:
* bc-wallet-tron-hdwallet-api - regular worker user for api
* bc-wallet-tron-hdwallet-migrator - user for run database migration
* bc-wallet-tron-hdwallet-updater - user for run infrastructure migration

#### Local migrations

```
make migrate
```

## Nats

Pattern to create buckets or streams is
```
%s__BC-WALLET-TRON-HDWALLET__MNEMONIC-WALLETS 
```
%s - is STAGE name, for example - dev, prod, bc_team1 (personal test stand for bc team)

### Create Buckets


```
nats kv add DEV__BC-WALLET-TRON-HDWALLET__MNEMONIC-WALLETS --replicas 1 --history 3 --storage=memory --description="mnemonic wallets cache storage for TRON hdwallet service"
```
### Clear Buckets
```
nats kv purge DEV__BC-WALLET-TRON-HDWALLET__MNEMONIC-WALLETS
```

### Vault

In production environment application require 3 Vault tokens:
* bc-wallet-tron-hdwallet-api 
* bc-wallet-tron-hdwallet-migrator
* bc-wallet-tron-hdwallet-updater

#### Vault local prepare
```bash
vault token create -display-name bc-wallet-tron-hdwallet-api
vault token create -display-name bc-wallet-tron-hdwallet-migrator
vault token create -display-name bc-wallet-tron-hdwallet-updater

vault write -f transit/keys/crypto-bundle/bc-wallet-tron-hdwallet

vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/common DB_DATABASE=bc-wallet-tron-hdwallet
vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/api \
  DB_USERNAME=bc-wallet-tron-hdwallet-api \
  DB_PASSWORD=password \
  
vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/migrator \
  DB_USERNAME=bc-wallet-tron-hdwallet-migrator \
  DB_PASSWORD=password \
  REDIS_USER=bc-wallet-tron-hdwallet-migrator \
  REDIS_PASSWORD=password
  NATS_USER=bc-wallet-tron-hdwallet-migrator \
  NATS_PASSWORD=password
  
vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/updater \
  DB_USERNAME=bc-wallet-tron-hdwallet-updater \
  DB_PASSWORD=password
  REDIS_USER=bc-wallet-tron-hdwallet-updater \
  REDIS_PASSWORD=password
  NATS_USER=bc-wallet-tron-hdwallet-updater \
  NATS_PASSWORD=password
```

## Licence

bc-wallet-tron-hdwallet is licensed under the [MIT](./LICENSE) License.