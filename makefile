build-all:
	cd migrator && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker-compose up --force-recreate --build -d

stop-all:
	docker-compose down