
GO=go
GOBUILDFLAGS=-v -ldflags="-v -s -w -X main.VersionString=$(PROGRAM_VERSION)"
PROGRAM=thalos-server
PROGRAM_VERSION=0.1.1
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
CFGDIR=$(PREFIX)/etc/thalos

.PHONY: build build/$(PROGRAM) build/benchmark test

build: build/$(PROGRAM)

build/$(PROGRAM) :
	$(GO) build $(GOBUILDFLAGS) -o $@ cmd/thalos/main.go

build-benchmark : build/benchmark

build/benchmark :
	$(GO) build $(GOBUILDFLAGS) -o $@ cmd/bench/main.go

install: build
	install -D build/$(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)
	install -m 644 -D config.example.yml $(DESTDIR)$(CFGDIR)/config.yml

install-scripts:
	install -m 755 -t $(DESTDIR) scripts/start.sh scripts/stop.sh

test:
	$(GO) test -v ./...

clean :
	$(RM) -fr build
