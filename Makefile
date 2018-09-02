PKGS := $(shell go list ./... | grep -v /vendor)
BIN_DIR := $(GOPATH)/bin
DEP_BIN := $(BIN_DIR)/dep
GOMETALINTER := $(BIN_DIR)/gometalinter

dep-ensure: $(DEP_BIN)
	$(DEP_BIN) ensure -v

$(DEP_BIN):
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

.PHONY: test
test: dep-ensure
	go test $(PKGS)

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor
