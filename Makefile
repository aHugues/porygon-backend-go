all: clean build run

build:
	mkdir -p build
	go build -o build

clean:
	go clean
	rm -r build

run:
	./build/porygon-backend-go

pretest-psql:
	psql -U azure porygontest -f ./ci/empty-db.sql
	psql -U azure porygontest -f ./ci/prepare-db-psql.sql

test-psql: pretest-psql
	mkdir -p build/
	go test -v -coverprofile=build/c-no-services.out ./models
	go test -v -coverprofile=build/c-services.out ./services
	sed -i 1d build/c-services.out
	cat build/c-no-services.out build/c-services.out > build/c.out
	go tool cover -html=build/c.out -o build/coverage.html

test_ci:
	go test -v -race -coverprofile=coverage.txt ./...

.PHONY: all test clean build