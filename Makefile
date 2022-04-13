# Server control
start:
	docker-compose up -d --build
	# Open logs
	detach xdg-open http://localhost:8888

stop:
	docker-compose down

rm-volumes:
	docker volume rm $(docker volume ls -q)

check-build:
	go build ./cmd/api
	rm api

# Test
test-coverage:
	mkdir -p "coverage"
	go test ./internal/usecase ./internal/data/postgres -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	rm coverage/coverage.out
	# See http://inglorion.net/software/detach/
	detach xdg-open coverage/coverage.html

test:
	go test ./internal/usecase ./internal/data/postgres -v

# Migrations
migrate-create:
	migrate create -ext sql -dir ./migrations -seq ${NAME}
