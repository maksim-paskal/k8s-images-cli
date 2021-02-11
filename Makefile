test:
	./scripts/validate-license.sh
	go fmt ./cmd
	go mod tidy
	go test -race ./cmd
	golangci-lint run --allow-parallel-runners -v --enable-all --disable testpackage --fix
run:
	export CGO_ENABLED=0
	export GOFLAGS="-trimpath"
	go build -o k8s-images-cli -v ./cmd
	./k8s-images-cli -logLevel=INFO $(args)
build:
	@./scripts/build-all.sh
	ls -lah _dist
	go mod tidy