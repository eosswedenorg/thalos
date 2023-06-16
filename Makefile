
GO=go
GOBUILDFLAGS=-v -p $(shell nproc) -ldflags="-v -s -w -X main.VersionString=$(PROGRAM_VERSION)"
PROGRAM=thalos-server
PROGRAM_VERSION=0.1.2
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
CFGDIR=$(PREFIX)/etc/thalos

.PHONY: build build/$(PROGRAM) build/thalos-tools test

build: build/$(PROGRAM)

build/$(PROGRAM) :
	$(GO) build $(GOBUILDFLAGS) -o $@ cmd/thalos/main.go

tools : build/thalos-tools

build/thalos-tools :
	$(GO) build $(GOBUILDFLAGS) -o $@ cmd/tools/main.go cmd/tools/bench.go cmd/tools/validate.go

install: build tools
	install -D build/$(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)
	install -D build/thalos-tools $(DESTDIR)$(BINDIR)/thalos-tools
	install -m 644 -D config.example.yml $(DESTDIR)$(CFGDIR)/config.yml

install-scripts:
	install -m 755 -t $(DESTDIR) scripts/start.sh scripts/stop.sh

test:
	$(GO) test -v ./...

clean :
	$(RM) -fr build
