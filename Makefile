# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	"$(CURDIR)/scripts/tidy.sh"

## audit: run quality control checks
.PHONY: audit
audit:
	"$(CURDIR)/scripts/audit.sh"

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	$(CURDIR)/scripts/test_cover.sh

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	$(CURDIR)/scripts/test.sh

## build/go: build the go application
.PHONY: build/app
build/app:
	$(CURDIR)/scripts/build_app.sh

## build/go: build the go application
.PHONY: build/web
build/web:
	$(CURDIR)/scripts/build_web.sh

## build/docker: build the application as a Docker image
.PHONY: build/docker
build/docker:
	$(CURDIR)/scripts/build_docker.sh

## run: run the  application
.PHONY: run
run: build/app
	bin/fineasy

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## schema/up: apply database schema migrations
.PHONY: schema/up
schema/up:
	$(CURDIR)/scripts/schema.sh up

## schema/down: rollback database schema migrations
.PHONY: schema/down
schema/down:
	$(CURDIR)/scripts/schema.sh down

## pg/up: start the PostgreSQL database container for local development
.PHONY: pg/up
pg/up:
	$(CURDIR)/scripts/pg_up.sh