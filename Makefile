WORKDIR=./tmp
DB_DSN=postgres://postgres:tgbot@localhost/postgres?sslmode=disable
VERSION=2

wire-gen:
	wire ./cmd/api

docker-build:
	rm -rf $(WORKDIR)
	mkdir $(WORKDIR)
	cp go.mod $(WORKDIR)
	cp go.sum $(WORKDIR)
	cp -r cmd $(WORKDIR)
	cp -r internal $(WORKDIR)
	docker build -f ./Dockerfile . --tag ghcr.io/eretic431/datingtelegrambot
	rm -rf $(WORKDIR)

docker-push: docker-build
	docker push ghcr.io/eretic431/datingtelegrambot:latest

migrate_up:
	migrate -path=./migrations -database=${DB_DSN} up ${VERSION}

migrate_down:
	migrate -path=./migrations -database=${DB_DSN} down ${VERSION}

migrate_force:
	migrate -path=./migrations -database=${DB_DSN} force ${VERSION}

migrate_create:
	migrate create -seq -ext=.sql -dir=./migrations ${NAME}

test-coverage:
	mkdir -p "coverage"
	go test ./internal/usecase ./internal/data/postgres -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	rm coverage/coverage.out
	# See http://inglorion.net/software/detach/
	detach xdg-open coverage/coverage.html

test:
	go test ./internal/usecase ./internal/data/postgres -v