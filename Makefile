WORKDIR=./tmp
DB_DSN=postgres://postgres:tgbot@localhost/postgres?sslmode=disable

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
