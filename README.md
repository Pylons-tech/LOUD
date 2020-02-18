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
'F': go to forest
'S': go to shop
'M': go to market
'T': go to settings
'D': go to develop
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
##### Development channel

Development channel is available and to do automation process on development channel

To do automation process correctly with afti's java, copy artifacts_txutil.sh file which has below content
```
java -cp "walletcore_txutil_jar/walletcore_txutil.jar:walletcore_txutil_jar/*" com.pylons.fuzzer.Main $1 $2 $3 $4 $5
```
And also copy walletcore_txutil_jar folder in project root scope.

Run below command to run automation.
```
make ARGS="afti -locald -userest -automate" run
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
Development channel is available and to do automation process on remote node
```
make ARGS="afti -userest -automate" run
```