# Install plugins:
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
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

migrate:
	 goose -dir ./migrations postgres "host=data4.gdrn.me port=5434 user=bc-wallet-tron-hdwallet password=password dbname=bc-wallet-tron-hdwallet sslmode=disable" up

build:
	docker build -t cr.selcloud.ru/crypto-bundle/bc-wallet-tron-hdwallet:latest .
	docker push cr.selcloud.ru/crypto-bundle/bc-wallet-tron-hdwallet:latest

default: hdwallet

deploy:
	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))

	docker buildx build --platform linux/amd64,linux/arm64 --push -t cr.selcloud.ru/crypto-bundle/bc-wallet-tron-hdwallet:$(build_tag) .

	helm --kubeconfig ~/.kube/kubenet.config --kube-context microk8s upgrade \
		--install bc-wallet-tron-hdwallet-api \
		--set "global.build_tag=$(build_tag)" --set "global.env=$(env)"= ./deploy/helm/api \
		--values=./deploy/helm/api/values.yaml --values=./deploy/helm/api/values_$(env).yaml

.PHONY: migrate hdwallet build deploy
