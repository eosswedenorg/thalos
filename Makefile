
GO=go
PROGRAM=build/eosio-ship-trace-reader

.PHONY: $(PROGRAM) test

$(PROGRAM) :
	$(GO) build -o $@

test:
	$(GO) test -v ./...

clean :
	$(RM) -fr build
