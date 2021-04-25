GO=go1.16

listener:
	$(GO) run cmd/listener/listener.go
.PHONY: listener

processor:
	$(GO) run cmd/processor/processor.go
.PHONY: processor

storer:
	$(GO) run cmd/storer/storer.go
.PHONY: storer

webserver:
	$(GO) run cmd/webserver/webserver.go
.PHONY: webserver

cexp:
	$(GO) run cmd/cexp/cexp.go
.PHONY: cexp

pexp:
	$(GO) run cmd/pexp/pexp.go
.PHONY: pexp

test:
	$(GO) test ./...
.PHONY: test

clean:
	$(GO) clean -testcache
.PHONY: clean
