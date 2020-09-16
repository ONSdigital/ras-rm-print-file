NAME := ras-rm-printfile
MAIN_SRC_FILE=cmd/ras-rm-print-file/main.go

.PHONY: test
test:
	go test  ./...

.PHONY: build
build:
	go build -o build/$(NAME) $(MAIN_SRC_FILE)
