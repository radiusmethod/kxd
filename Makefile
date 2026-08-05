BINDIR = /usr/local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -s -w -X github.com/radiusmethod/kxd/src/cmd.version=$(VERSION)

help:          ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

install:       ## Install Target
	GOOS= GOARCH= GOARM= GOFLAGS= go build -ldflags="$(LDFLAGS)" -o ${BINDIR}/kxd
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=-=-=- "
	@echo "        __ __   _  __    ____       "
	@echo "       / //_/  | |/ /   / __ \      "
	@echo "      / ,<     |   /   / / / /      "
	@echo "     / /| |   /   |   / /_/ /       "
	@echo "    /_/ |_|  /_/|_|  /_____/        "
	@echo "                                    "
	@echo "    To Finish Installation add      "
	@echo "                                    "
	@echo "    eval \"\$$(kxd init zsh)\"          "
	@echo "                                    "
	@echo "  to your zshrc (or bash profile,   "
	@echo "  with 'init bash') then open a new "
	@echo "  terminal or source that file      "
	@echo "                                    "
	@echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=-=-=- "

uninstall:     ## Uninstall Target
	rm -f ${BINDIR}/kxd

.PHONY: test test-coverage docs
test:          ## Run tests
	go test ./...

test-coverage: ## Run tests with coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

docs:          ## Regenerate docs/*.md from the cobra command tree
	cd tools/gendocs && go run . ../../docs
