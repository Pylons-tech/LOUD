# Commands list

```log
pylonscli tx broadcast automation/tx_pyloncli_signed.json
pylonscli keys show -a afti3135
./artifacts_txutil.sh AUTO_CREATE_COOKBOOK cee82ada86c96b51b402e8f67469178ba241d054874d029ccd255aa3b62fb5cf 98 5 automation/msg.json 
./artifacts_txutil.sh SIGNED_TX d18a973d6a8c0cb9d145778fe9a2dd73bf8f58a2bac872e92df45736cac931bc 110 3 automation/msg.json
pylonscli query account cosmos1r8eazsulhl7g6zg6czz9378l6dg5kug8v4zea0
pylonscli tx sign automation/tx.json --account-number 110 --sequence 3 --offline --from artipsi5


2020/02/20 12:39:46 comparing afticli and pyloncli ;) nJNr+c+Ln1waqluQgDvZ8QDVvpFfEPJWd/Nc+yD3noUGEaJp+twEzx9P6TEmw7uhZcAxU7NM6HRPh9ToKIXTyA== 
and
 8O5jkXt/zE8MHSzybKW6Hof6Fy9Jfp9gGM5AahVHvacPdzeEevjQwlcHSRIOgrZjvvxGCYOr0d+eK4QTKxlJtw==
2020/02/20 12:39:46 where
2020/02/20 12:39:46 msg= {"type":"pylons/CreateCookbook","value":{"CookbookID":"1582166381","Name":"tst_cookbook_name","Description":"addghjkllsdfdggdgjkkk","Version":"1.0.0","Developer":"asdfasdfasdf","SupportEmail":"a@example.com","Level":"0","Sender":"cosmos1n57x7ej944sccavqsu8eyvqysc8ymjuzdszjh9","CostPerBlock":"5"}}
2020/02/20 12:39:46 username= artipsi5
2020/02/20 12:39:46 Bech32Addr= cosmos1n57x7ej944sccavqsu8eyvqysc8ymjuzdszjh9
2020/02/20 12:39:46 privKey= d18a973d6a8c0cb9d145778fe9a2dd73bf8f58a2bac872e92df45736cac931bc
2020/02/20 12:39:46 account-number= 110
2020/02/20 12:39:46 sequence 3
```
