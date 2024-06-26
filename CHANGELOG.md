# Change Log

## [v0.0.30] 21.06.2024
### Changed
* Changed config init flow
* Bump common-lib-vault version - bc-wallet-common-lib-vault v0.0.17
  * Added new environment variable VAULT_AUTH_TOKEN_FILE_PATH to *-example.env files
  * Implemented new vault client init flow
* Bump common-lib-postgres version - bc-wallet-common-lib-postgres v0.0.9

## [v0.0.29] 08.06.2024
### Changed
Changed MIT License to NON-AI MIT

## [v0.0.28] 07.06.2024
### Added
* Added PoW-shield to wallet-api request
  * Added pow_proofs database table
* Added Access-token shield
  * Added access_tokens database table
  * Added wallet_sessions_access_tokens database table
### Changed
* Controller-API divided into two parts - ManagerApi, WalletAPI
  * ManagerAPI require common system JWT-token for access
  * ManagerAPI now contains next methods:
    * AddNewWallet
    * ImportWallet
    * EnableWallet
    * GetWalletInfo
    * GetEnabledWallets
    * DisableWallet
    * DisableWallets
    * EnableWallets
    * GetAccount
  * WalletAPI require wallet-level JWT access token, which will be created in AddNewWallet method
  * WalletAPI now contains next methods:
    * GetWalletInfo
    * StartWalletSession
    * GetWalletSession
    * GetAllWalletSessions
    * CloseWalletSession
    * GetAccount
    * GetMultipleAccounts
    * PrepareSignRequest
    * ExecuteSignRequest

## [v0.0.27] 12.05.2024
### Added
* Added usage of lib-profiler
* Added migration container dockerfile
### Fixed
* Fixed bugs in hdwallet-client dial flow - search in socket file
* Fixed bug in shutdown flow of controller-api gRPC server
### Changed
* Changed vault transit key config - now names of transit keys loaded from vault
* Added task in Makefile for building migration container

## [v0.0.26] 07.05.2024
### Changed
* Bump version of bc-wallet-common-lib-vault v0.0.14
* Some changes in build container flow
* Added info to README.md file

## [v0.0.25] 05.05.2024
### Added
* Added indexes to sign_request and mnemonic_wallet_sessions tables
* Added integration tests for hdwallet-controller gRPC api
### Changed
* Changed GetDerivationAddress and GetDerivationAddressByRange gRPC methods:
  * Replaced by GetAccount and GetMultipleAccounts
  * Removed derivation_path field from sign_requests table
  * Added account_data field to sign_requests table
### Fixed
* Fixed bug in PrepareSignRequest flow - usage of old sign_request table fields

## [v0.0.24] 24.04.2024
### Added
* Added license banner to all *.go files
* New hdwallet-controller gRPC-method:
  * ImportWallet
  * EnableWallet
  * DisableWallet
  * DisableWallets
  * EnableWallets
  * GetAllWalletSessions
  * CloseWalletSession
  * PrepareSignRequest
  * ExecuteSignRequest
* New methods in hdwallet-api proto description:
  * EncryptMnemonic
  * ValidateMnemonic
  * UnLoadMultipleMnemonics
  * LoadDerivationAddress
  * SignData
### Changed
* Proto description separated by 3 files:
  * common
  * controller_api
  * hdwallet_api
* Removed helm-chart and helm deployment - move to *-hdwallet repository
* Removed nats cache
* Bump go version 1.19 -> 1.22
* Bump common-lib version:
  * bc-wallet-common-lib-config v0.0.5
  * bc-wallet-common-lib-grpc v0.0.4
  * bc-wallet-common-lib-healthcheck v0.0.4
  * bc-wallet-common-lib-logger v0.0.4
  * bc-wallet-common-lib-nats-queue v0.1.12
  * bc-wallet-common-lib-postgres v0.0.8
  * bc-wallet-common-lib-redis v0.0.7
  * bc-wallet-common-lib-tracer v0.0.4
  * bc-wallet-common-lib-vault v0.0.13

## [v0.0.23] 14.02.2024
### Info
Start of big application refactoring
### Added
* Added wallet sessions entities for storing in persistent and cache stores
### Changed
* Separated application on two parts
  * bc-wallet-common-hdwallet-controller
  * bc-wallet-tron-hdwallet
* Changed GetDerivationAddressByRange gRPC method - now support get addresses by multiple ranges
* Added HdWallet API proto description
  * new gRPC method - GenerateMnemonic
  * new gRPC method - LoadMnemonic
  * new gRPC method - UnLoadMnemonic
* Added Controller API proto description
  * new gRPC method - StartWalletSession
  * new gRPC method - GetWalletSession
* Removed go-tron-sdk dependency

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

## [v0.0.4] 16.03.2023
### Changed
* Refactoring wallet pool service-component:
  * Added wallet pool unit
  * Added unit maker
  * Added support of multiple and single mnemonic wallet
  * Added timer for mnemonic unloading flow

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