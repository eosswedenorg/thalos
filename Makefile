
GO=go
GOBUILDFLAGS=-v --buildmode=pie -ldflags="-v -s -w -X main.VersionString=$(PROGRAM_VERSION)"
PROGRAM=thalos-server
PROGRAM_VERSION=0.1.0
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
CFGDIR=$(PREFIX)/etc/thalos

.PHONY: build build/$(PROGRAM) test

build: build/$(PROGRAM)

build/$(PROGRAM) :
	$(GO) build $(GOBUILDFLAGS) -o $@ cmd/thalos/main.go

install: build
	install -D build/$(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)
	install -m 644 -D config.example.yml $(DESTDIR)$(CFGDIR)/config.yml

test:
	$(GO) test -v ./...

clean :
	$(RM) -fr build
