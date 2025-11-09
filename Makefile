# ========= Variables =========
GOOSE_SERVICE    := bff
GOOSE_WORKDIR    := /go/src/app
GOOSE_DIR        := db/migrations
GOOSE_CMD        := cd $(GOOSE_WORKDIR) && goose -dir $(GOOSE_DIR) mysql "$$DB_USERNAME:$$DB_PASSWORD@tcp($$DB_HOST:3306)/$$DB_DATABASE?parseTime=true"

BACKEND_SERVICE  := backend
BACKEND_WORKDIR  := /go/src/app

PROTO_DIR        := grpc/proto
PROTO_SERVICE    := grpc
PROTOC_OUTPUT_DIRS := bff/pkg/pb backend/pkg/pb
ROOT_PROTO_FILES := $(wildcard *.proto)
PROTO_FILES_FROM_DIR := $(shell if [ -d $(PROTO_DIR) ]; then find $(PROTO_DIR) -name '*.proto'; fi)
STRIPPED_PROTO_FILES := $(patsubst $(PROTO_DIR)/%,%,$(PROTO_FILES_FROM_DIR))
PROTO_FILES      := $(strip $(ROOT_PROTO_FILES) $(STRIPPED_PROTO_FILES))
PROTO_INCLUDE_PATHS := --proto_path=. $(if $(PROTO_FILES_FROM_DIR),--proto_path=$(PROTO_DIR))

# ========= PHONY =========
.PHONY: \
  goose-up goose-status goose-down \
  backend-mock-category backend-mock backend-test backend-go-test \
  gqlgen proto _require_proto_files \
  docker-shell grpc-shell \
  up down restart logs

# ========= Docker (optional) =========
up:
	docker compose up -d

down:
	docker compose down

restart: down up

logs:
	docker compose logs -f

# ========= goose =========
goose-up:
	docker compose run --rm $(GOOSE_SERVICE) sh -c '$(GOOSE_CMD) up'

goose-down:
	docker compose run --rm $(GOOSE_SERVICE) sh -c '$(GOOSE_CMD) down'

goose-status:
	docker compose run --rm $(GOOSE_SERVICE) sh -c '$(GOOSE_CMD) status'

# ========= gqlgen =========
gqlgen:
	docker compose run --rm bff sh -c 'gqlgen generate'

# ========= backend helpers =========
backend-mock:
	cd backend && set -e; \
		go run github.com/golang/mock/mockgen@v1.6.0 -destination=domain/repository/mock/category_repository_mock.go -package=mock backend/domain/repository CategoryRepository; \
		go run github.com/golang/mock/mockgen@v1.6.0 -destination=domain/repository/mock/task_repository_mock.go -package=mock backend/domain/repository TaskRepository; \
		go run github.com/golang/mock/mockgen@v1.6.0 -destination=domain/repository/mock/subtask_repository_mock.go -package=mock backend/domain/repository SubTaskRepository; \
		go run github.com/golang/mock/mockgen@v1.6.0 -destination=domain/repository/mock/user_repository_mock.go -package=mock backend/domain/repository UserRepository

backend-mock-category: backend-mock
	@echo "Generated repository mocks."

backend-test:
	docker compose run --rm $(BACKEND_SERVICE) sh -c 'cd $(BACKEND_WORKDIR) && go test ./...'

backend-go-test:
	cd backend && tmp_dir=$$(mktemp -d); \
		GOCACHE=$$tmp_dir go test ./...; \
		rm -rf $$tmp_dir

# ========= protobuf =========
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

# ========= Shell =========
docker-shell:
	docker compose run --rm $(PROTO_SERVICE) bash

grpc-shell:
	docker compose run --rm $(PROTO_SERVICE) bash
