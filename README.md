## How it look like

![Pylons Central](https://github.com/Pylons-tech/LOUD/blob/master/screenshots/pylons_central.png)

## Setup development environment

```sh
git clone https://github.com/Pylons-tech/LOUD
brew install pre-commit
brew install golangci/tap/golangci-lint
pre-commit install
```

## Building

You need a correctly set up Golang environment over 1.11; this proejct uses `go mod`.

Run make.
```sh
    make
```
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
'M': go to pylons central
'T': go to settings
'D': go to develop
```

## To run test

Install pylonscli at `GOPATH/bin/pylonscli`.
Here GOPATH is configured as /Users/{username}/go

Since Pylons repos are private below command will be needed to run test correctly.  

```sh
export GOPRIVATE="github.com/Pylons-tech"
go get "github.com/Pylons-tech/pylons_sdk"
```

### Fresh Local test
Build game
```sh
make
```
Run daemon and rest-server
```sh
pylonsd unsafe-reset-all
pylonsd start
pylonscli rest-server --chain-id pylonschain
```
Create eugen account by running
```sh
make ARGS="eugen -locald" run
```
Create cookbook and recipes with eugen account in test keyring-backend by running
```sh
make fixture_tests ARGS="-locald --accounts=eugen"
```
Run game with name "michael"
```sh
make ARGS="michael -locald" run
```
Run game with name "michael" and use rest endpoint for tx send
```sh
make ARGS="michael -locald -userest" run
```
Deployment of recipes with custom account, this replace account1 to eugen at runtime.
##### Development channel

Development channel is available and to do automation process on development channel

To do automation process correctly with afti's java, copy artifacts_txutil.sh file which has below content
```sh
java -cp "jar/txutil.jar:jar/*" com.pylons.txutil.Main $1 $2 $3 $4 $5
```
And also copy `jar` folder in project root scope.

Run below command to run automation.
```sh
make ARGS="afti -locald -userest -automate" run
```
##### For Afti

Once you see "something went wrong" message on loud app, you can just see loud.log and search for "comparing afticli and pyloncli ;)" and you will be able to see relevant logs on before and after.
It looks like below.

```log
2020/02/18 comparing afticli and pyloncli ;) OTuhvokIB+vfqNxx5AAHXbD+xmlq/l7HiOtoJ+805YlkHIr+bnNWlCTYtxMf06w6isk+OGMgLL9MjIx64EVprA== 
and
 ouhd/DfYAgsXycksZz+bRXFsToWqOPe6XTC6ph7smEp/+CjummjBKzVQecIEJSMkBAvu+5kbmroMXqw51Qb73w==
2020/02/18 where
2020/02/18 19:47:08 msg= {"type":"pylons/CreateCookbook","value":{"CookbookID":"1582019224","Name":"tst_cookbook_name","Description":"addghjkllsdfdggdgjkkk","Version":"1.0.0","Developer":"asdfasdfasdf","SupportEmail":"a@example.com","Level":"0","Sender":"cosmos1tdfk4ec383nftjavzdtr5mg5uxhjgzzmkvp4jv","CostPerBlock":"5"}}
2020/02/18 username= afti112
2020/02/18 Bech32Addr= cosmos1tdfk4ec383nftjavzdtr5mg5uxhjgzzmkvp4jv
2020/02/18 privKey= a1e1247936bced4d713fa9ff18edde4dc381ef031c77aa049352a6e1ec1abb72
2020/02/18 account-number= 69
2020/02/18 sequence 13
```
You can get message to be signed, sequence, account number, private key, bech32addr, username on pyloncli, and tx sign result.

### Fresh Remote node test
Build game
```sh
make
```

Create eugen account by running
```sh
make ARGS="eugen" run
```

After creating eugen account on remote, check eugen account created correctly on node by running
```sh
pylonscli query account $(pylonscli keys show -a eugen) --node tcp://35.223.7.2:26657
```
Update cookbook name, id and recipe ID for correct version and timestamp.
**Warn:** If cookbook name does not change, it's refering to old version of cookbook since now it's finding cookbooks by name.

Deployment of recipes and cookbooks witu eugen account in test keyring-backend by running
```sh
make fixture_tests ARGS="-runserial --accounts=eugen"
```
Deployment of recipes using REST endpoint.
```sh
make fixture_tests ARGS="-userest --accounts=eugen"
```
Deployment of recipes with using existing cookbook (known cookbook)
```sh
make fixture_tests ARGS="-use-known-cookbook --accounts=eugen"
```
Verification stuff for recipes deployment
```sh
make fixture_tests ARGS="-verify-only --accounts=eugen"
```
Deployment of recipes with custom account, this replace account1 to eugen at runtime.

Check if cookbook and all recipes are created by using
```sh
pylonscli query pylons list_recipe --node tcp://35.223.7.2:26657
```

If something went wrong, just create the remaining recipes by editing `scenario/loud.json`.

Run game with name "michael"
```sh
make ARGS="michael" run
```

Run game with rest endpoint with name "michael"
```sh
make ARGS="michael -userest" run
```

Development channel is available and to do automation process on remote node
```sh
make ARGS="afti -userest -automate" run
```

### To run multiple instances on same computer

Just clone repos two times in different places something like below
```
/Users/admin/go/src/github.com/MikeSofaer/LOUD
/Users/admin/go/src/github.com/MikeSofaer/LOUD_C
```

And then run `make` and `make ARGS="XXX XXX" run` on each folder.

### How to create release version

- Download Platypus app from https://sveinbjorn.org/platypus 
- Run Platypus app and set sh script and icon for the game.
- After generating game, right click on app and select "Open package contents"
- Build pyloncli and loud game using go on a mac machine.
- Go to "Contents/Resources/" and paste resource files there like "bin", "locale".
- Publish new version

### How to debug log file

To find the errors happened in that log, you need to search for `level=warning`, `level=error`, `level=fatal`, `level=panic`.
Once that's found you can debug the `debug` and `info` logs near that log and it will make sense for you what happened.

To set log level directly for log file, you can update something like this, `testing.NewT(nil)` to `NewLogLevelT(nil, log.DebugLevel)`.

Set custom logging level at `config.yml`.
```
# sdk configurations
sdk:
  # trace, debug, info, warn, error, fatal, panic
  log_level: trace
```
For local work, use `config_local.yml`. I used to have the configuration issue when I configured on `config.yml` and did run app on local node.

### Configuration

```yaml
# sdk configurations
sdk:
  max_wait_block: 3 # if transaction does not exist on chain until after 3 blocks, game think that transaction failed
  rest_endpoint: http://35.238.123.59:80 # rest endpoint is not being used on loud game 
  cli_endpoint: tcp://35.223.7.2:26657 # we broadcast to remote node for loud game
  log_level: trace # log level configuration: trace, debug, info, warn, error, fatal, panic

# app configurations
app:
  daemon_timeout_commit: 2 # this should be same as timeout_commit configured at pylonsd configuration
```
