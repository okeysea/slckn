#####################################################
# Build
#####################################################
SRCS := $(shell find . -type f -name '*.go')

# Alpineで動かすのでStatic化
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""
TARGETS := $(shell find ./cmd/ -type d | sed 's!^.*/!!')

all: $(TARGETS)
$(TARGETS):
	go build -a -buildvcs=false -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$@ ./cmd/$@

# bin/$(NAME): $(SRCS)
# 	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

#####################################################
# Compose Utils
#####################################################
COMMON_COMPOSE=nerdctl compose -f docker-compose.yml
.PHONY: up upd down down-all
build:
	${COMMON_COMPOSE} build

up:
	${COMMON_COMPOSE} up

upd:
	${COMMON_COMPOSE} up -d

down:
	${COMMON_COMPOSE} down

down-all:
	${COMMON_COMPOSE} down --rmi all --volumes --remove-orphans

#####################################################
# Golang
#####################################################
.PHONY: go-sh go-tidy go-fmt
go-sh:
	${COMMON_COMPOSE} run --rm app ash

go-tidy:
	${COMMON_COMPOSE} run --rm app sh -c 'go mod tidy'

go-fmt:
	${COMMON_COMPOSE} run --rm app sh -c 'go fmt ./...'

#####################################################
# Misc
#####################################################
.PHONY: fix-perm
# Fix Permissions
fix-perm:
	sudo chown -R $(shell whoami):$(shell whoami) .
