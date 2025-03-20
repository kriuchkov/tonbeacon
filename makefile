# Description: Makefile for tonbeacon
vendor:
	go mod tidy && go mod vendor

compose-up:
	docker compose up --build --attach scanner consumer
	
compose-up-d:
	docker compose up -d --build --attach scanner consumer

compose-up-required:
	docker compose up postgres flyway kafka kafka-ui topic-creator


compose-down:
	docker compose down

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/grpc/v1/tonbeacon.proto


.PHONY: vendor compose-up compose-up-d compose-up-required compose-down gen-proto