
GO=go
PROGRAM=build/eosio-ship-trace-reader

.PHONY: $(PROGRAM)

$(PROGRAM) :
	$(GO) build -o $@

clean :
	$(RM) -fr build
