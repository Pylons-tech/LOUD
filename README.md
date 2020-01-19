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

Since Pylons repos are private below command will be needed to run test correctly.  

```
export GOPRIVATE="github.com/MikeSofaer"
```

Fresh Local test

```
pylonsd unsafe-reset-all
pylonsd start
pylonscli rest-server --chain-id pylonschain
make fixture_tests
make
rm nonce.json
./bin/loud
```

Fresh Remote node test
Search for below code
```
func SetupScreenAndEvents(world World) {
	user := world.GetUser(
```
And change the username to "eugen"
```
make
make run
```

Check eugen account created correctly on node by running
```
pylonscli query account $(pylonscli keys show -a eugen) --node 35.223.7.2:26657
```

make fixture_tests

Check if all recipes are created by using
```
pylonscli query pylons list_recipe --node 35.223.7.2:26657
```

If something went wrong, just create the remaining recipes.

rm world.db
rm nonce.json
change username to "afti" or your name
make
make ARGS="eugen" run