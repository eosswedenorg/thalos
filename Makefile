
GO=go
GOLDFLAGS=-v -s -w -X main.VersionString=$(PROGRAM_VERSION)
GOBUILDFLAGS=-v -p $(shell nproc) -ldflags="$(GOLDFLAGS)"
PROGRAM=thalos-server
PROGRAM_VERSION=1.1.1
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
CFGDIR=$(PREFIX)/etc/thalos

.PHONY: build build/$(PROGRAM) build/thalos-tools test

build: build/$(PROGRAM)

build/$(PROGRAM) :
	$(GO) build $(GOBUILDFLAGS) -o $@ cmd/thalos/main.go cmd/thalos/server.go

tools : build/thalos-tools

build/thalos-tools :
	$(GO) build $(GOBUILDFLAGS) -o $@ $(shell find cmd/tools -type f -name *.go)

install: build tools
	install -D build/$(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)
	install -D build/thalos-tools $(DESTDIR)$(BINDIR)/thalos-tools
	install -m 644 -D config.example.yml $(DESTDIR)$(CFGDIR)/config.yml

install-scripts:
	install -m 755 -t $(DESTDIR) scripts/start.sh scripts/stop.sh

build-deb:
	dpkg-buildpackage -b -us -uc

test:
	$(GO) test -v ./...
	cd api; $(GO) test -v ./...
clean :
	$(RM) -fr build
