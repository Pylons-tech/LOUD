
## Coins

loudcoin

## Items
Lv1 wooden sword
```
{
    "Doubles": [],
    "Longs": [{ "Key": "level", "Value": "1" }],
    "Strings": [{ "Key": "Name","Value": "Wooden sword" }],
    "CookbookName": "Legend of Undead Dragon",
    "Sender": "michael"
}
```
Lv2 wooden sword
```
{
    "Doubles": [],
    "Longs": [{ "Key": "level", "Value": "2" }],
    "Strings": [{ "Key": "Name","Value": "Wooden sword" }],
    "CookbookName": "Legend of Undead Dragon",
    "Sender": "michael"
}
```
Lv1 copper sword
```
{
    "Doubles": [],
    "Longs": [{ "Key": "level", "Value": "1" }],
    "Strings": [{ "Key": "Name","Value": "Copper sword" }],
    "CookbookName": "Legend of Undead Dragon",
    "Sender": "michael"
}
```

Lv2 copper sword
```
{
    "Doubles": [],
    "Longs": [{ "Key": "level", "Value": "2" }],
    "Strings": [{ "Key": "Name","Value": "Copper sword" }],
    "CookbookName": "Legend of Undead Dragon",
    "Sender": "michael"
}
```


## Recipes

### Hunt recipe without sword (1, 5, 10)
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Count": 1,
            "Weight":1
        }, {
            "Coin":"loudcoin",
            "Count": 5,
            "Weight": 2
        }, {
            "Coin":"loudcoin",
            "Count": 10,
            "Weight": 1
        }],
        "ItemOutputs":[]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's hunt without sword recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to hunt without sword.",
    "BlockInterval":"0"
}
```

### Hunt recipe with lv1 Wooden sword  (10, 15, 20)
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1"}],
        "Strings": [{"Key": "Name", "Value": "Wooden sword"}]
    }],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Count": 10,
            "Weight":1
        }, {
            "Coin":"loudcoin",
            "Count": 15,
            "Weight": 2
        }, {
            "Coin":"loudcoin",
            "Count": 20,
            "Weight": 1
        }],
        "ItemOutputs":[]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's hunt with lv1 wooden sword recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to hunt with lv1 wooden sword.",
    "BlockInterval":"0"
}
```

### Hunt recipe with Wooden sword lv2 (20-30)
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "2", "MaxValue": "2"}],
        "Strings": [{"Key": "Name", "Value": "Wooden sword"}]
    }],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Count": 20,
            "Weight":1
        }, {
            "Coin":"loudcoin",
            "Count": 25,
            "Weight": 2
        }, {
            "Coin":"loudcoin",
            "Count": 30,
            "Weight": 1
        }],
        "ItemOutputs":[]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's hunt with lv2 wooden sword recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to hunt with lv2 wooden sword.",
    "BlockInterval":"0"
}
```

### Hunt recipe with Copper  sword lv1 (50-80)
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1"}],
        "Strings": [{"Key": "Name", "Value": "Copper sword"}]
    }],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Count": 50,
            "Weight":1
        }, {
            "Coin":"loudcoin",
            "Count": 65,
            "Weight": 2
        }, {
            "Coin":"loudcoin",
            "Count": 80,
            "Weight": 1
        }],
        "ItemOutputs":[]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's hunt with lv1 copper sword recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to hunt with lv1 copper sword.",
    "BlockInterval":"0"
}
```

### Hunt recipe with Copper sword lv2 (80-120)
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "2", "MaxValue": "2"}],
        "Strings": [{"Key": "Name", "Value": "Copper sword"}]
    }],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Count": 80,
            "Weight":1
        }, {
            "Coin":"loudcoin",
            "Count": 100,
            "Weight": 2
        }, {
            "Coin":"loudcoin",
            "Count": 120,
            "Weight": 1
        }],
        "ItemOutputs":[]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's hunt with lv2 copper sword recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to hunt with lv2 copper sword.",
    "BlockInterval":"0"
}
```

### Buy Wooden sword lv1 price: 100 gold 
```
{
    "RType": "0",
    "CoinInputs":[{
        "Coin": "loudcoin",
        "Count": "100"
    }],
    "ItemInput": [],
    "Entries":{
        "CoinOutputs":[],
        "ItemOutputs":[
            {
                "Doubles":[],
                "Longs":[
                    {
                        "Rate":"1.0",
                        "Key":"level",
                        "WeightRanges":[{ "Lower": 1, "Upper":1,"Weight":1 }]
                    }
                ],
                "Strings":[{ "Key":"Name", "Value":"Wooden sword", "Rate":"1.0" }],
                "Weight":1
            }
        ]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Wooden sword lv1 buy recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to buy wooden sword lv1.",
    "BlockInterval":"0"
}
```

### Wooden sword lv1 to lv2 upgrade price: 100 gold
```
{
    "RType": "1",
    "CoinInputs":[{
        "Coin": "loudcoin",
        "Count": "100"
    }],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1"}],
        "Strings": [{"Key": "Name", "Value": "Wooden sword"}]
    }],
    "ToUpgrade": {
        "Doubles": [],
        "Longs": [{
            "Key": "level", 
            "WeightRanges":[{ "Lower": 1, "Upper":1,"Weight":1 }]
        }],
        "Strings": []
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Wooden sword lv1 to lv2 upgrade recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to upgrade wooden sword level.",
    "BlockInterval":"0"
}
```


### Buy Copper sword lv1 price: 250 gold
```
{
    "RType": "0",
    "CoinInputs":[{
        "Coin": "loudcoin",
        "Count": "250"
    }],
    "ItemInput": [],
    "Entries":{
        "CoinOutputs":[],
        "ItemOutputs":[
            {
                "Doubles":[],
                "Longs":[
                    {
                        "Rate":"1.0",
                        "Key":"level",
                        "WeightRanges":[{ "Lower": 1, "Upper":1,"Weight":1 }]
                    }
                ],
                "Strings":[{ "Key":"Name", "Value":"Copper sword", "Rate":"1.0" }],
                "Weight":1
            }
        ]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Copper sword lv1 buy recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to buy copper sword lv1.",
    "BlockInterval":"0"
}
```

### Copper sword lv1 to lv2 upgrade price: 250 gold
```
{
    "RType": "1",
    "CoinInputs":[{
        "Coin": "loudcoin",
        "Count": "250"
    }],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1"}],
        "Strings": [{"Key": "Name", "Value": "Copper sword"}]
    }],
    "ToUpgrade": {
        "Doubles": [],
        "Longs": [{
            "Key": "level", 
            "WeightRanges":[{ "Lower": 1, "Upper":1,"Weight":1 }]
        }],
        "Strings": []
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Copper sword lv1 to lv2 upgrade recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to upgrade copper sword level.",
    "BlockInterval":"0"
}
```

### Wooden sword lv1 sell recipe, 80 gold
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1"}],
        "Strings": [{"Key": "Name", "Value": "Wooden sword"}]
    }],
    "Entries": {
        "CoinOutputs": [{
            "Coin":"loudcoin",
            "Count": 80,
            "Weight":1
        }]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Lv1 wooden sword sell recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to sell lv1 wooden sword.",
    "BlockInterval":"0"
}
```

### Wooden sword lv2 sell recipe, 160 gold
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "2", "MaxValue": "2"}],
        "Strings": [{"Key": "Name", "Value": "Wooden sword"}]
    }],
    "Entries": {
        "CoinOutputs": [{
            "Coin":"loudcoin",
            "Count": 160,
            "Weight":1
        }]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Lv2 wooden sword sell recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to sell lv1 wooden sword.",
    "BlockInterval":"0"
}
```
### Copper sword lv1 sell recipe, 200 gold
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1"}],
        "Strings": [{"Key": "Name", "Value": "Copper sword"}]
    }],
    "Entries": {
        "CoinOutputs": [{
            "Coin":"loudcoin",
            "Count": 200,
            "Weight":1
        }]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Lv1 copper sword sell recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to sell lv1 copper sword.",
    "BlockInterval":"0"
}
```
### Copper sword lv2 sell recipe, 400 gold
```
{
    "RType": "0",
    "CoinInputs":[],
    "ItemInput": [{
        "Doubles": [],
        "Longs": [{"Key": "level", "MinValue": "2", "MaxValue": "2"}],
        "Strings": [{"Key": "Name", "Value": "Copper sword"}]
    }],
    "Entries": {
        "CoinOutputs": [{
            "Coin":"loudcoin",
            "Count": 400,
            "Weight":1
        }]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's Lv2 copper sword sell recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to sell lv2 copper sword.",
    "BlockInterval":"0"
}
```