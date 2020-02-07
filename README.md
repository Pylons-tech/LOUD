## Building

You need a correctly set up Golang environment over 1.11; this proejct uses `go mod`.

Run make.

    make

## Running

Then run `bin/loud` from this folder.

## Closing

Pressing `Esc` key on keyboard makes you to end your game.

## Playing

Key action descriptions

```
'H': go to home
'S': go to shop
'F': go to forest
```

## To run test

Install pylonscli at `GOPATH/bin/pylonscli`.
Here GOPATH is configured as /Users/{username}/go

Since Pylons repos are private below command will be needed to run test correctly.  

```
export GOPRIVATE="github.com/Pylons-tech"
```

### Fresh Local test
Build game
```
make
```
Run daemon and rest-server
```
pylonsd unsafe-reset-all
pylonsd start
pylonscli rest-server --chain-id pylonschain
```
Create eugen account by running
```
make ARGS="eugen -locald" run
```
Create cookbook and recipes by running
```
make fixture_tests ARGS="-locald"
```
Run game with name "michael"
```
make ARGS="michael -locald" run
```
Run game with name "michael" and use rest endpoint for tx send
```
make ARGS="michael -locald -userest" run
```
### Fresh Remote node test
Build game
```
make
```
Create eugen account by running
```
make ARGS="eugen" run
```
After creating eugen account on remote, check eugen account created correctly on node by running
```
pylonscli query account $(pylonscli keys show -a eugen) --node 35.223.7.2:26657
```
Create recipes and cookbooks by running
```
make fixture_tests ARGS="-runserial"
```
Check if cookbook and all recipes are created by using
```
pylonscli query pylons list_recipe --node 35.223.7.2:26657
```
If something went wrong, just create the remaining recipes by editing `scenario/loud.json`.

Run game with name "michael"
```
make ARGS="michael" run
```

Run game with rest endpoint with name "michael"
```
make ARGS="michael -userest" run
```