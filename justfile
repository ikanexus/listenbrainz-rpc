build:
    go build -o listenbrainz-rpc main.go

run *args:
    go run main.go {{args}}

test:
    go test ./...

tidy:
    go mod tidy

lint:
    golangci-lint run

release-snapshot:
    goreleaser release --clean --snapshot

release:
    goreleaser release --clean

upgrade-deps *args:
  go-mod-upgrade {{args}}
