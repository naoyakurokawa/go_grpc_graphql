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
