# bc-wallet-tron-hdwallet

## Nats

Pattern to create buckets or streams is
```
%s__BC_WALLET_TRON_HDWALLET__CRYPTO_BUNDLE__ETHEREUM
```
%s - is ENVIRONMENT name, for example - dev, prod, bc_team1 (personal test stand for bc team)

### Create Streams
```
nats stream add --config ./deploy/nats/create_wallet_stream_cb_eth.json
```

## K8s

### Secrets

```
kubectl create secret generic bc-wallet-tron-hdwallet \
  --from-file=bc-wallet-tron-hdwallet-rsa-key=./build/secrets/rsa/private.pem \
  --from-literal=redis_username= --from-literal=redis_password='password' \
  --from-literal=nats_username='user' --from-literal=nats_password='password' \
  --from-literal=db_name='bc-wallet-tron-hdwallet' --from-literal=db_username='bc-wallet-tron-hdwallet' \
  --from-literal=db_password='password' --from-literal=vault_auth_path='<insert_token_here>' \
  --from-literal=vault_transit_secret_key='bc-wallet-tron-hdwallet'
```

## DB

### Local migrations

```
make migrate
```

## Enryption

```
openssl genrsa -out ./build/secrets/rsa/private.pem 4096
```


## Vault
```bash
vault secrets enable -address=http://data.tr.gdrn.me:8200 transit
vault token create -address=http://data.tr.gdrn.me:8200 -display-name bc-wallet-tron-hdwallet
vault write -address=http://data.tr.gdrn.me:8200 -f transit/keys/bc-wallet-tron-hdwallet
vault secrets enable -address=http://data.tr.gdrn.me:8200 -path=kv kv
vault kv put -address=http://data.tr.gdrn.me:8200 kv/bc-wallet-tron-hdwallet DB_USERNAME=bc-wallet-tron-hdwallet
vault kv put -address=http://data.tr.gdrn.me:8200 kv/bc-wallet-tron-hdwallet DB_PASSWORD=password
vault kv put -address=http://data.tr.gdrn.me:8200 kv/bc-wallet-tron-hdwallet REDIS_USER=bc-wallet-tron-hdwallet
vault kv put -address=http://data.tr.gdrn.me:8200 kv/bc-wallet-tron-hdwallet REDIS_PASSWORD=password
vault kv put -address=http://data.tr.gdrn.me:8200 kv/bc-wallet-tron-hdwallet NATS_USER=nats-user
vault kv put -address=http://data.tr.gdrn.me:8200 kv/bc-wallet-tron-hdwallet NATS_USER=password
```