# bc-wallet-eth-hdwallet

## Nats

Pattern to create buckets or streams is
```
%s__BC_WALLET_ETH_HDWALLET_LIBER__ETHEREUM
```
%s - is ENVIRONMENT name, for example - dev, prod, bc_team1 (personal test stand for bc team)

### Create Streams
```
nats stream add --config ./deploy/nats/create_wallet_stream_liber_eth.json
```

## K8s

### Secrets

```
kubectl create secret generic bc-wallet-eth-hdwallet \
  --from-file=bc-wallet-eth-hdwallet-rsa-key=./build/secrets/rsa/private.pem \
  --from-literal=redis_username= --from-literal=redis_password='password' \
  --from-literal=nats_username='user' --from-literal=nats_password='password' \
  --from-literal=db_name='bc-wallet-eth-hdwallet' --from-literal=db_username='bc-wallet-eth-hdwallet' --from-literal=db_password='password'
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