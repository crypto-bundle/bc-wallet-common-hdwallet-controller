# bc-wallet-tron-hdwallet

## Description

Service of management owned hd-wallets - mnemonic wallets for Tron blockchain. 
Application purpose - create new mnemonic wallet, sign transaction, get derivation address

## Api

Service has a GRPC-api - api doc [here](./docs/hdwallet_api/hdwallet_proto.md)

## Mnemmonic wallets

Mnemonic wallets stored in Postgresql database as encrypted Hashicord **Vault** data

## Requirements

### k8s
Helm deploy describes in [./deploy/helm/api](./deploy/helm/api)

### PostgreSQL
* Database: bc-wallet-tron-hdwallet
* Users:
  * bc-wallet-tron-hdwallet-api - SELECT, INSERT, UPDATE privileges
  * bc-wallet-tron-hdwallet-migrator - CREATE, DELETE, DROP, SELECT, INSERT, UPDATE privileges
  * bc-wallet-tron-hdwallet-updater - CREATE, DELETE, DROP, SELECT, INSERT, UPDATE privileges

### Nats
* Users:
  * bc-wallet-tron-hdwallet-api
  * bc-wallet-tron-hdwallet-migrator
  * bc-wallet-tron-hdwallet-updater
* KV buckets:
  * <STAGE_PREFIX>__BC_WALLET_TRON_HDWALLET__MNEMONIC-WALLETS

### Vault
* Users:
  * bc-wallet-tron-hdwallet-api
  * bc-wallet-tron-hdwallet-migrator
  * bc-wallet-tron-hdwallet-updater
* Tokens:
  * vault token create -display-name bc-wallet-tron-hdwallet-api
  * vault token create -display-name bc-wallet-tron-hdwallet-migrator
  * vault token create -display-name bc-wallet-tron-hdwallet-updater
* Buckets:
    * kv/crypto-bundle/bc-wallet-tron-hdwallet/commo
    * kv/crypto-bundle/bc-wallet-tron-hdwallet/api
    * kv/crypto-bundle/bc-wallet-tron-hdwallet/migrator
    * kv/crypto-bundle/bc-wallet-tron-hdwallet/updater
* Transit
  * vault write -f transit/keys/crypto-bundle/bc-wallet-tron-hdwallet

### k8s secrets
* Vault tokens:
    * vault_bc_wallet_tron_hdwallet_api_user_token
    * vault_bc_wallet_tron_hdwallet_migrator_user_token
    * vault_bc_wallet_tron_hdwallet_updater_user_token

## Deployment

### K8s

#### Secrets

```
kubectl create secret generic bc-wallet-tron-hdwallet \
  --from-literal=vault_api_user_token='<insert_token_here>' \
  --from-literal=vault_migrator_user_token='<insert_token_here>' \
  --from-literal=vault_updater_user_token='<insert_token_here>' \
  --from-literal=vault_transit_secret_key='crypto-bundle-bc-wallet-tron-hdwallet'
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

vault write -f transit/keys/crypto-bundle-bc-wallet-tron-hdwallet

vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/common DB_DATABASE=bc-wallet-tron-hdwallet
vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/api \
  DB_USERNAME=bc-wallet-tron-hdwallet-api \
  DB_PASSWORD=password \
  
vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/migrator \
  DB_USERNAME=bc-wallet-tron-hdwallet-migrator \
  DB_PASSWORD=password \
  REDIS_USER=bc-wallet-tron-hdwallet-migrator \
  REDIS_PASSWORD=password \
  NATS_USER=bc-wallet-tron-hdwallet-migrator \
  NATS_PASSWORD=password
  
vault kv put kv/crypto-bundle/bc-wallet-tron-hdwallet/updater \
  DB_USERNAME=bc-wallet-tron-hdwallet-updater \
  DB_PASSWORD=password \
  REDIS_USER=bc-wallet-tron-hdwallet-updater \
  REDIS_PASSWORD=password \
  NATS_USER=bc-wallet-tron-hdwallet-updater \
  NATS_PASSWORD=password
```

## Licence

**bc-wallet-tron-hdwallet** has a proprietary license.

Switched to proprietary license from MIT - [CHANGELOG.MD - v0.0.24](./CHANGELOG.md)