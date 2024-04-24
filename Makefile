# Install plugins:
#  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#  go get -d github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

hdwallet_proto:
	protoc -I ./pkg/proto/ -I . -I ./pkg/proto/ \
    		--go_out=./pkg/grpc/hdwallet_api/proto/ \
    		--go_opt=paths=source_relative \
    		--go-grpc_out=./pkg/grpc/hdwallet_api/proto/ \
    		--go-grpc_opt=paths=source_relative \
    		--openapiv2_out=logtostderr=true:./docs/hdwallet_api/ \
    		--grpc-gateway_out=./pkg/grpc/hdwallet_api/proto/ \
    		--grpc-gateway_opt=logtostderr=true \
    		--grpc-gateway_opt=paths=source_relative \
    		--doc_out=./docs/hdwallet_api/ \
    		--doc_opt=markdown,$@.md \
    		./pkg/proto/*.proto

default: hdwallet

deploy:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval container_registry=$(repository)/crypto-bundle/bc-wallet-tron-hdwallet)
	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	docker build --no-cache --platform $(platform) --tag $(container_registry):$(build_tag) .
	docker push $(container_registry):$(build_tag)

	helm --kube-context $(context) upgrade \
		--install bc-wallet-tron-hdwallet-api \
		--set "global.container_registry=$(container_registry)" \
		--set "global.build_tag=$(build_tag)" \
		--set "global.env=$(env)" \
		--values=./deploy/helm/api/values.yaml \
		--values=./deploy/helm/api/values_$(env).yaml \
		./deploy/helm/api

.PHONY: hdwallet_proto deploy
