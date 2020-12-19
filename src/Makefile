listener:
	go1.15 run cmd/listener/listener.go
.PHONY: listener

processor:
	go1.15 run cmd/processor/processor.go
.PHONY: processor

storer:
	go1.15 run cmd/storer/storer.go
.PHONY: storer

webserver:
	go1.15 run cmd/webserver/webserver.go
.PHONY: webserver

cexp:
	go1.15 run cmd/cexp/cexp.go
.PHONY: cexp

pexp:
	go1.15 run cmd/pexp/pexp.go
.PHONY: pexp

test:
	go1.15 test ./...
.PHONY: test

clean:
	go1.15 clean -testcache
.PHONY: clean
