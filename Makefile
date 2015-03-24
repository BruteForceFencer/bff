CC=gcc
DEST_DIR=out

all: bffctl core daemonize

bffctl: directories
	go build -o $(DEST_DIR)/bin/bffctl bffctl/*.go

core: directories
	go build -o $(DEST_DIR)/util/bffcore core/*.go
	cp core/config.json $(DEST_DIR)/config.json
	cp -r core/assets $(DEST_DIR)/assets

daemonize: directories
	$(CC) -o $(DEST_DIR)/util/daemonize daemonize/main.c

directories:
	mkdir -p $(DEST_DIR)
	mkdir -p $(DEST_DIR)/bin
	mkdir -p $(DEST_DIR)/util
	mkdir -p $(DEST_DIR)/var

clean:
	rm -rf $(DEST_DIR)
