# Install plugins:
#  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#  go get -d github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

common_proto:
	protoc -I ./pkg/proto/common/ \
		--go_out=./pkg/grpc/common/ \
		--go_opt=paths=source_relative \
		--go-grpc_out=./pkg/grpc/common \
		--go-grpc_opt=paths=source_relative \
		--openapiv2_out=logtostderr=true:./docs/api/ \
		--grpc-gateway_out=./pkg/grpc/common \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=source_relative \
		--doc_out=./docs/api/ \
		--doc_opt=markdown,$@.md \
		./pkg/proto/common/*.proto

hdwallet_proto:
	protoc -I ./pkg/proto/hdwallet_api/ -I ./pkg/proto/common/ \
		--go_out=./pkg/grpc/hdwallet/ \
		--go_opt=paths=source_relative \
		--go-grpc_out=./pkg/grpc/hdwallet/ \
		--go-grpc_opt=paths=source_relative \
		--openapiv2_out=logtostderr=true:./docs/api/ \
		--grpc-gateway_out=./pkg/grpc/hdwallet/ \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=source_relative \
		--doc_out=./docs/api/ \
		--doc_opt=markdown,$@.md \
		./pkg/proto/hdwallet_api/*.proto

controller_proto:
	protoc -I ./pkg/proto/controller_api/ -I ./pkg/proto/common/ \
    		--go_out=./pkg/grpc/controller/ \
    		--go_opt=paths=source_relative \
    		--go-grpc_out=./pkg/grpc/controller/ \
    		--go-grpc_opt=paths=source_relative \
    		--openapiv2_out=logtostderr=true:./docs/api/ \
    		--grpc-gateway_out=./pkg/grpc/controller/ \
    		--grpc-gateway_opt=logtostderr=true \
    		--grpc-gateway_opt=paths=source_relative \
    		--doc_out=./docs/api/ \
    		--doc_opt=markdown,$@.md \
    		./pkg/proto/controller_api/*.proto

default: hdwallet

build_container:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval container_registry=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-controller)
	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	$(eval short_commit_id=$(shell git rev-parse --short HEAD))
	$(eval commit_id=$(shell git rev-parse HEAD))
	$(eval build_number=0)
	$(eval build_date=$(shell date +%s))
	$(eval release_tag=$(shell git describe --tags $(commit_id))-$(short_commit_id)-$(build_number))

	docker build \
		--ssh default=$(SSH_AUTH_SOCK) \
		--platform $(platform) \
		--build-arg RELEASE_TAG=$(release_tag) \
		--build-arg COMMIT_ID=$(commit_id) \
		--build-arg SHORT_COMMIT_ID=$(short_commit_id) \
		--build-arg BUILD_NUMBER=$(build_number) \
		--build-arg BUILD_DATE_TS=$(build_date) \
		--tag $(container_registry):$(release_tag) \
		--tag $(container_registry):latest .

	docker push $(container_registry):$(release_tag)
	docker push $(container_registry):latest

.PHONY: hdwallet_proto build_container
