
build-daemon:
	go build -ldflags="-w -s" -tags daemon -o bin/ ./cmd/...
	go build -ldflags="-w -s" -o bin/ ./cmd/dconn  # dconn is a special case

build-daemonless:
	go build -ldflags="-w -s" -o bin/ ./cmd/...

build: build-daemonless

install-daemon:
	go install -tags daemon ./cmd/...
	go install ./cmd/dconn

install-daemonless:
	go install ./cmd/...

install: install-daemonless
