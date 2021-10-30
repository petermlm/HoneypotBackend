GO=go1.17

# =============================================================================
# Main Go Services
# =============================================================================

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

# =============================================================================
# Other Go Services
# =============================================================================

cexp:
	$(GO) run cmd/cexp/cexp.go
.PHONY: cexp

pexp:
	$(GO) run cmd/pexp/pexp.go
.PHONY: pexp

# =============================================================================
# Tests
# =============================================================================

test:
	$(GO) test ./...
.PHONY: test

# =============================================================================
# Clean
# =============================================================================

clean:
	$(GO) clean -testcache
.PHONY: clean
