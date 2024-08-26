include .env
LOCAL_BIN:=$(CURDIR)/bin
CHAT_API:=chat_v1
CERT_FOLDER:=cert

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/$(CHAT_API)
	protoc --proto_path api/$(CHAT_API) \
	--go_out=pkg/$(CHAT_API) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/$(CHAT_API) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/$(CHAT_API)/chat.proto

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATIONS_DIR} postgres ${PG_DSN} down -v

gen-cert:
	mkdir -p $(CERT_FOLDER)
	openssl genrsa -out $(CERT_FOLDER)/ca.key 4096 && \
    openssl req -new -x509 -key $(CERT_FOLDER)/ca.key -sha256 -subj "/C=RU/ST=Moscow/O=Test, Inc." -days 365 -out $(CERT_FOLDER)/ca.cert && \
    openssl genrsa -out $(CERT_FOLDER)/service.key 4096 && \
    openssl req -new -key $(CERT_FOLDER)/service.key -out $(CERT_FOLDER)/service.csr -config $(CERT_FOLDER)/certificate.conf && \
    openssl x509 -req -in $(CERT_FOLDER)/service.csr -CA $(CERT_FOLDER)/ca.cert -CAkey $(CERT_FOLDER)/ca.key -CAcreateserial \
        -out $(CERT_FOLDER)/service.pem -days 365 -sha256 -extfile $(CERT_FOLDER)/certificate.conf -extensions req_ext