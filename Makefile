OUTFILES := $(patsubst cmd/%.go,bin/%,$(wildcard cmd/*.go))

bin/%: cmd/%.go
	go build -o $@ $<

all: clean mod $(OUTFILES)

mod:
	go mod download

clean:
	rm bin/* || true

run:
	./bin/loud

fixture_tests:
	rm ./test/nonce.json || true
	go test -v ./test/