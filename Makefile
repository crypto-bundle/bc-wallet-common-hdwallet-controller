# Install plugins:
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#  go get -d github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

hdwallet:
	protoc -I ./pkg/proto/ -I . -I ./pkg/proto/ \
    		--go_out=./pkg/grpc/hd_wallet_api/proto/ \
    		--go_opt=paths=source_relative \
    		--go-grpc_out=./pkg/grpc/hd_wallet_api/proto/ \
    		--go-grpc_opt=paths=source_relative \
    		--openapiv2_out=logtostderr=true:./docs/hd_wallet_api/ \
    		--grpc-gateway_out=./pkg/grpc/hd_wallet_api/proto/ \
    		--grpc-gateway_opt=logtostderr=true \
    		--grpc-gateway_opt=paths=source_relative \
    		--doc_out=./docs/hd_wallet_api/ \
    		--doc_opt=markdown,$@.md \
    		./pkg/proto/*.proto

migrate:
	 goose -dir ./migrations postgres "host=192.168.1.103 port=5434 user=bc-wallet-eth password=password dbname=bc-wallet-eth sslmode=disable" up


default: hdwallet

.PHONY: hdwallet
