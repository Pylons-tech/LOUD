OUTFILES := $(patsubst cmd/%.go,bin/%,$(wildcard cmd/*.go))

bin/%: cmd/%.go
	go build -o $@ $<

all: clean mod $(OUTFILES)

mod:
	go mod download

clean:
	rm bin/* || true

run:
	pylonscli config chain-id pylonschain
	pylonscli config output json
	pylonscli config indent true
	pylonscli config trust-node true
	./bin/loud ${ARGS}

fixture_tests:
	rm ./test/nonce.json || true
	go test -v ./test/ ${ARGS}