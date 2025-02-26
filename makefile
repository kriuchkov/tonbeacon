

vendor:
	go mod tidy && go mod vendor

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/grpc/v1/tonbeacon.proto