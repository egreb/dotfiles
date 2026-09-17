.PHONY: build test check install
build:
	go build -o dist/sbx ./cmd/sbx

test:
	go test ./...

check: test build
	go vet ./...
	WORKFLOW="$(CURDIR)/dist/sbx" python3 tests/test-sbx.py

install: build
	mkdir -p "$(HOME)/.config/sbx/bin"
	install -m 755 dist/sbx "$(HOME)/.config/sbx/bin/sbx.next"
	mv "$(HOME)/.config/sbx/bin/sbx.next" "$(HOME)/.config/sbx/bin/sbx"
