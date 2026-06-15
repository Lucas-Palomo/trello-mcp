BIN_NAME := trello-mcp
OUT_DIR ?= dist
DEST ?=

.PHONY: build install clean

build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BIN_NAME) .

install: build
	test -n "$(DEST)"
	mkdir -p $(DEST)
	cp $(OUT_DIR)/$(BIN_NAME) $(DEST)/$(BIN_NAME)

clean:
	rm -rf $(OUT_DIR)
