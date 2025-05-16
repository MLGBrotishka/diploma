LOCAL_BIN:=$(CURDIR)/bin

build:
	go build -o bin/migrator cmd/migrator/main.go
PHONY: build

run:
	go run cmd/migrator.go
PHONY: run

test:
	go test -coverprofile=cover.out -p=3 -count=2 -cover -tags=integration -race  ./... 
PHONY: test

mock:
	go generate ./...
PHONY: mock

# Proto
# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
	https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

# Вендоринг внешних proto файлов
.vendor-proto: .vendor-rm  vendor-proto/google/protobuf

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@latest \
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/godoc@latest

MIGRATOR_PROTO_PATH:=api

.PHONY: generate 
generate: # .bin-deps # .vendor-proto
	protoc \
	-I api \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out pkg/${MIGRATOR_PROTO_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out pkg/${MIGRATOR_PROTO_PATH} \
    --go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out pkg/${MIGRATOR_PROTO_PATH} \
    --grpc-gateway_opt paths=source_relative \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	--openapiv2_out ./api/	 \
	api/migrator/migrator.proto \
	api/auth/auth.proto
	go mod tidy



# Документация
swagger_port ?= 8081

.PHONY: swagger
swagger: # .bin-deps
	bin/swagger serve --no-open --port $(swagger_port) --flavor=swagger api/migrator/migrator.swagger.json

godoc_port ?= 8082

.PHONY: godoc
godoc: #.bin-deps
	godoc -http=:$(godoc_port)