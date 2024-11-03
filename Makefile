
GO=go
GOLDFLAGS=-v -s -w -X main.VersionString=$(PROGRAM_VERSION)
GOBUILDFLAGS+=-v -p $(shell nproc) -ldflags="$(GOLDFLAGS)"
PROGRAM=thalos-server
PROGRAM_VERSION ?= 1.1.7-rc2
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
CFGDIR=$(PREFIX)/etc/thalos
DOCKER_IMAGE_REPO ?= ghcr.io/eosswedenorg/thalos
DOCKER_IMAGE_TAG ?= $(PROGRAM_VERSION)

.PHONY: build build/$(PROGRAM) build/thalos-tools test docker-image docker-publish

build: build/$(PROGRAM)

build/$(PROGRAM) :
	$(GO) build $(GOBUILDFLAGS) -o $@ ./cmd/thalos/

tools : build/thalos-tools

build/thalos-tools :
	$(GO) build $(GOBUILDFLAGS) -o $@ ./cmd/tools/

docker-image:
	docker image build --build-arg VERSION=$(PROGRAM_VERSION) -t $(DOCKER_IMAGE_REPO):$(DOCKER_IMAGE_TAG) docker

docker-publish:
	docker image push $(DOCKER_IMAGE_REPO):$(DOCKER_IMAGE_TAG)

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
