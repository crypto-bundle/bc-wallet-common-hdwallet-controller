# Change Log

## [initial] - 06.03.2022 - 16.03.2023
* Created go module as bc-wallet-eth-hdwallet
* Added proto files for gRPC API
* Integrated common dependencies
* Moved to crypto-bundle namespace
* Added wallet persistent store
* Added functionality for gRPC handlers
  * AddNewWallet
  * GetDerivationAddress
  * GetEnabledWallets
  * GetDerivationAddressByRange
* Added mnemonic encryption via rsa-keys
* Added MIT licence
* Refactoring service for supporting TRON blockchain
* Created Helm chart

## [v0.0.4] 16.03.2023
### Changed
* Refactoring wallet pool service-component:
  * Added wallet pool unit
  * Added unit maker 
  * Added support of multiple and single mnemonic wallet
  * Added timer for mnemonic unloading flow

## [v0.0.5] 05.04.2023
### Added
* Encryption private data via hashicorp vault
* Added gRPC client config 
### Changed
* Cleaned up repository:
  * Removed ansible database deployment script
  * Removed vault polices
  * Removed private data from helm-chart
* Updated common-libs:
  * removed old bc-wallet-common dependency
  * integrated lib-common dependencies:
    * lib-postgres
    * lib-config
    * lib-grpc
    * lib-tracer
    * lib-logger
    * lib-vault
### Fixed
* Fixed bug in wallet init stage
* Fixed crash in wallet pool init stage
* Fixed bugs in flow in new wallet creation

## [v0.0.6 - v0.0.22] 05.04.2023 - 28.04.2023
### Added
* Added gRPC client wrapper 
* Small security improvements:
  * Filling private keys with zeroes - private key clearing
* Added data cache flow for storing wallet in redis and nats 
* Added new gRPC-handler - GetWalletInfo
### Changed
* Changed deployment flow
  * Added helm-chart option for docker container repository  
  * Fixed helm-chart template for VAULT_DATA_PATH variable
* Optimization in get addresses by range flow
### Fixed
* Fixed bug in sign transaction flow
* Fixed migrations - wrong rollback SQL-code, missing drop index and drop table