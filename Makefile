.PHONY: goose-up goose-status gqlgen

# Run goose migrations inside the bff container
GOOSE_SERVICE := bff
GOOSE_WORKDIR := /go/src/app
GOOSE_DIR := db/migrations
GOOSE_CMD := cd $(GOOSE_WORKDIR) && goose -dir $(GOOSE_DIR) mysql "$$DB_USERNAME:$$DB_PASSWORD@tcp($$DB_HOST:3306)/$$DB_DATABASE?parseTime=true"

goose-up:
	docker compose run --rm $(GOOSE_SERVICE) sh -c '$(GOOSE_CMD) up'

goose-status:
	docker compose run --rm $(GOOSE_SERVICE) sh -c '$(GOOSE_CMD) status'

# Run gqlgen code generation for the BFF GraphQL schema
gqlgen:
	docker compose run --rm bff sh -c 'gqlgen generate'

PROTO_DIR := grpc/proto
PROTO_SERVICE := grpc
PROTOC_OUTPUT_DIRS := bff/pkg/pb backend/pkg/pb
ROOT_PROTO_FILES := $(wildcard *.proto)
PROTO_FILES_FROM_DIR := $(shell if [ -d $(PROTO_DIR) ]; then find $(PROTO_DIR) -name '*.proto'; fi)
STRIPPED_PROTO_FILES := $(patsubst $(PROTO_DIR)/%,%,$(PROTO_FILES_FROM_DIR))
PROTO_FILES := $(strip $(ROOT_PROTO_FILES) $(STRIPPED_PROTO_FILES))
PROTO_INCLUDE_PATHS := --proto_path=. $(if $(PROTO_FILES_FROM_DIR),--proto_path=$(PROTO_DIR))

proto: _require_proto_files
	@set -e; \
	if [ -z "$(strip $(PROTOC_OUTPUT_DIRS))" ]; then \
		echo "No output directories specified for protobuf generation"; \
		exit 1; \
	fi; \
	for out_dir in $(PROTOC_OUTPUT_DIRS); do \
		echo "Generating protobuf files into $$out_dir"; \
		docker compose run --rm --build $(PROTO_SERVICE) protoc $(PROTO_INCLUDE_PATHS) \
			--go_out=$$out_dir --go_opt=paths=source_relative \
			--go-grpc_out=$$out_dir --go-grpc_opt=paths=source_relative \
			$(PROTO_FILES); \
	done

_require_proto_files:
	@if [ -z "$(strip $(PROTO_FILES))" ]; then \
		echo "No .proto files found in $(PROTO_DIR) or current directory"; \
		exit 1; \
	fi

docker-shell:
	docker compose run --rm $(PROTO_SERVICE) bash
