dev:
	air

cli:
	go run cmd/cli/main.go

setup:
	go install go.uber.org/mock/mockgen@latest

tidy:
	go mod tidy -v

t: test
test:
	go test ./internal/feature/... ./internal/core/...

tc: test.cov
test.cov:
	$(ENV_LOCAL_TEST) \
	go test -covermode=count -coverprofile=covprofile.out ./internal/feature/... ./internal/core/...
	make test.cov.xml

tc.xml: test.cov.xml
test.cov.xml:
	gocov convert covprofile.out > covprofile.xml

f: fmt
fmt:
	go fmt ./...

w: wire
wire:
	wire ./...

g: generate
generate:
	go generate ./...

b: build
build:
	go build -o apiserver ./cmd

docker:
	docker build -t server:latest .

docker.build:
	docker build -t coin-game-be -f Dockerfile.dev .

docker.dev:
	docker build -t coin-game-be:dev -f Dockerfile.dev .
