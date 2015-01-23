.PHONY : test

CURDIR=$(shell pwd)
BUILD_PATH=$(CURDIR)/build
BUILD_GOPATH=$(BUILD_PATH)/src/pkg.deepin.io/
FIXGOPATH=$(BUILD_PATH):$(GOPATH)
BIN_PATH=$(CURDIR)/bin

DDAEMON_SRC=$(CURDIR)/daemon

ifndef USE_GCCGO
    GOBUILD = go build
else
    LDFLAGS = $(shell pkg-config --libs gio-2.0)
    GOBUILD = go build -compiler gccgo -gccgoflags "${LDFLAGS}"
endif

build:
	mkdir -p $(BUILD_GOPATH)
	ln -s $(CURDIR) $(BUILD_GOPATH)/sync
	cd $(DDAEMON_SRC)  && GOPATH=$(FIXGOPATH) ${GOBUILD} -o $(BIN_PATH)/deepin-sync-service
	rm -r build

test:
	cd $(CURDIR) && GOPATH=$(FIXGOPATH) go test -v

install:
	install -Dm755 $(BIN_PATH)/deepin-sync-service $(DESTDIR)/usr/lib/deepin-daemon/deepin-sync-service

clean:
	@-rm -rf build/*
