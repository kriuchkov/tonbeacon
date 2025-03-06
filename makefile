# Description: Makefile for tonbeacon

vendor:
	go mod tidy && go mod vendor

compose-up:
	docker compose up --build --attach scanner

compose-down:
	docker compose down
	
compose-up-d:
	docker compose up -d --build --attach scanner

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/grpc/v1/tonbeacon.proto