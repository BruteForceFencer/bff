CC=gcc
DEST_DIR=build

all: core

core: directories
	go build -o $(DEST_DIR)/bff *.go
	cp config.json $(DEST_DIR)/config.json
	cp -r assets $(DEST_DIR)/assets

directories:
	mkdir -p $(DEST_DIR)

clean:
	rm -rf $(DEST_DIR)

.PHONY: clean directories
