# Change Log

## [v0.0.24] - 28.04.2023 17:50 MSK

### Changed

#### Switching to a proprietary license.
License of **bc-wallet-tron-hdwallet** repository changed to proprietary.

Origin repository - https://github.com/crypto-bundle/bc-wallet-tron-hdwallet

The MIT license is replaced by me (_Kotelnikov Aleksei_) as an author and maintainer.

The license has been replaced with a proprietary one, with the condition of maintaining the authorship
and specifying in the README.md file in the section of authors and contributors.

[@gudron (Kotelnikov Aleksei)](https://github.com/gudron) - author and maintainer of [crypto-bundle project](https://github.com/crypto-bundle)

The commit is signed with the key -
gudron2s@gmail.com
E456BB23A18A9347E952DBC6655133DD561BF3EC

## [v0.0.26] - 14.05.2023

### Changed
* Docker container build
* Fixed Nats-kv cache bucket bugs

## [v0.0.27] - 09.06.2023
### Added 
* Logs in main.go

## [v0.0.28] - 13.06.2023
### Fixed
* Create wallet bug. Bug in restoration from cache MnemonicWalletItem entity

## [v0.0.29] - 06.07.2023
### Fixed
* Application can't process init stage with empty wallets table

## [v0.0.30] - 11.07.2023
### Added
* sync.Pool usage in GetDerivationAddress gRPC method
* AddressIdentitiesCount parameter to GetDerivationAddressByRange response
### Changed
* GetDerivationAddressByRange gRPC method - added support of multiple ranges per request
* Added iterator pattern to request form - DerivationAddressByRangeForm

## [v0.0.31] - 13.07.2023
### Added
* GetDerivationAddressByRange method to gRPC-client wrapper

## [v0.0.32] - 14.07.2023
### Fixed
* Calculation of requested addresses count by range with GetDerivationAddressByRange gRPC method
* Added case of range with one address item - AddressRangeFrom is equals to AddressRangeTo

## [v0.0.33] - 30.08.2023
### Fixed
* GetAddressesByRange method - problem with gRPC-request cancel context.

## [v0.0.34 - v0.0.35] - 03.09.2023
### Fixed
* AddNewWallet method - problem with infinite context.Done() call. Because of wrong usage gRPC-request context
* Local deployment Makefile changes. Now supports building a Docker image via Podman

## [v0.0.36 - v0.0.37] - 23.09.2023
### Added
* Added examples of env-file for api and migrator applications
* Deployment k8s cluster context in Makefile helm deployment 
### Changed
* Loading of local .env file from path in APP_LOCAL_ENV_FILE_PATH environment variable

## [v0.0.38] - 08.10.2023
### Added
* AddNewWallet method to gRPC-client wrapper
* HotWalletIndex to response in add_wallet gRPC method

## [v0.0.39] - 14.11.2023
### Changed
* New version of bc-connector-common library - v1.3.21
* Migrated to usage of internal proto files descriptions
    * "github.com/fbsobreira/gotron-sdk/pkg/proto/api" replaced by "gitlab.heronodes.io/bc-platform/bc-connector-common/pkg/grpc/bc_adapter_api/proto/vendored/tron/node/api"
    * "github.com/fbsobreira/gotron-sdk/pkg/proto/core" replaced by "gitlab.heronodes.io/bc-platform/bc-connector-common/pkg/grpc/bc_adapter_api/proto/vendored/tron/node/core"
