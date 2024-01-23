CC         = go
BUILD_PATH = ./bin
SRC        = main.go
TARGET     = wpe
BINS       = $(BUILD_PATH)/$(TARGET)
INST       = /usr/local/bin
CONF_DEST  = ~/.config/wallpaper_engine/
CONF       = ./configs/config.json

.PHONY: all clean build run install

all: run

clean:
	rm -rf $(BUILD_PATH)

build: clean
	mkdir -p $(CONF_DEST)
	mkdir -p $(BUILD_PATH)
	$(CC) build -o $(BINS) $(SRC)

run: build
	$(BINS) --help

install: build 
	sudo cp -v $(BINS) $(INST)
	cp -vn $(CONF) $(CONF_DEST)
