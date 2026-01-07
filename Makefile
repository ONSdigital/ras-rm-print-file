NAME := ras-rm-printfile
MAIN_SRC_FILE=cmd/ras-rm-print-file/main.go

.PHONY: test
test:
	@if [ "$$(uname -m)" = "aarch64" ] || [ "$$(uname -m)" = "arm64" ]; then \
		GOARCH=arm64 \
		GOOS=darwin; \
	elif [ "$$(uname -m)" = "x86_64" ]; then \
		GOARCH=amd64 \
		GOOS=linux; \
	else \
		echo "Unsupported architecture: $$(uname -m)"; exit 1; \
	fi; \
	GOOS=$$GOOS CGO_ENABLED=0 GOARCH=$$GOARCH go test  ./...

.PHONY: build
build:
	@if [ "$$(uname -m)" = "aarch64" ] || [ "$$(uname -m)" = "arm64" ]; then \
		GOARCH=arm64 \
		GOOS=darwin; \
	elif [ "$$(uname -m)" = "x86_64" ]; then \
		GOARCH=amd64 \
		GOOS=linux; \
	else \
		echo "Unsupported architecture: $$(uname -m)"; exit 1; \
	fi; \
	GOOS=$$GOOS CGO_ENABLED=0 GOARCH=$$GOARCH go build -o build/$(NAME) $(MAIN_SRC_FILE)

.PHONY: fmt
fmt:
	go fmt  ./...

