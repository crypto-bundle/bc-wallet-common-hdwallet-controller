ARG PARENT_MIGRATION_CONTAINER_IMAGE_NAME="/crypto-bundle/bc-wallet-common-migrator:latest"

FROM $PARENT_MIGRATION_CONTAINER_IMAGE_NAME

COPY ./migrations /opt/bc-wallet-common-migrator/migrations