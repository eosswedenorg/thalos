
GO=go
PROGRAM=build/thalos-server

.PHONY: $(PROGRAM) test

$(PROGRAM) :
	$(GO) build -o $@ cmd/main/main.go

test:
	$(GO) test -v ./...

clean :
	$(RM) -fr build
