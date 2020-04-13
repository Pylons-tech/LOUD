#!/usr/bin/env bash

if [[ -z "${GOPATH}" ]]; then
  echo "export GOPATH"
  export GOPATH=$HOME/go
  mkdir $HOME/go
  mkdir $HOME/go/bin
else
  echo "GOPATH already exist"
fi

# copy pylonscli to $GOPATH/bin directory
echo "saving pylonscli to $GOPATH/bin/pylonscli"
cp ./bin/pylonscli $GOPATH/bin/pylonscli

# configuration for pylonscli
./bin/pylonscli config chain-id pylonschain
./bin/pylonscli config output json
./bin/pylonscli config indent true
./bin/pylonscli config trust-node true

# Run loud
./bin/loud $1 $2 $3 $4 $5
