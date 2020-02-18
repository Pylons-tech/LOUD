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
##### For Afti

Once you see "something went wrong" message on loud app, you can just see loud.log and search for "comparing afticli and pyloncli ;)" and you will be able to see relevant logs on before and after.
It looks like below.

```
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