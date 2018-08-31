PKGS := $(shell go list ./... | grep -v /vendor)
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: test
test:
	go test $(PKGS)

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor
