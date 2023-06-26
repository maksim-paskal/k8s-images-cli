test:
	./scripts/validate-license.sh
	go fmt ./cmd ./pkg/...
	go vet ./cmd ./pkg/...
	go mod tidy
	go test -race ./cmd ./pkg/...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v
run:
	go run -race ./cmd -logLevel=debug $(args)
test-release:
	go run github.com/goreleaser/goreleaser@latest release --snapshot --skip-publish --clean