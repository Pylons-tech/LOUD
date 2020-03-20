
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


### Hunt recipe without sword (1 or 2)
```
{
    "ID": "LOUD-hunt-with-no-weapon-recipe-v0.0.0-1579053457",
    "CoinInputs":[],    
    "ItemInput": [],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Program": "randi(2)+1",
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

### Hunt recipe with a sword (attack *4 or attack *5 )
```
{
    "ID": "LOUD-hunt-with-a-sword-recipe-v0.0.0-1583631194",
    "CoinInputs":[],
    "ItemInputs":[{
        "Doubles": [{"Key": "attack", "MinValue": "1.0", "MaxValue": "1000.0"}],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1000"}],
        "Strings": []
    }],
    "Entries":{
        "CoinOutputs":[{
            "Coin":"loudcoin",
            "Program": "int(attack * double(randi(2)+4))",
            "Weight": 1
        }],
        "ItemOutputs":[]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's hunt with a sword recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to hunt with a sword.",
    "BlockInterval":"0"
}
```
### Buy Wooden sword lv1 price: 100 gold 
```
{
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

### Sword sell recipe, attack * (randi(2)+20) gold
{
    "ID": "LOUD-hunt-with-a-sword-recipe-v0.0.0-1583631194",
    "CoinInputs":[],
    "ItemInputs":[{
        "Doubles": [{"Key": "attack", "MinValue": "1.0", "MaxValue": "1000.0"}],
        "Longs": [{"Key": "level", "MinValue": "1", "MaxValue": "1000"}],
        "Strings": []
    }],
    "Entries": {
        "CoinOutputs": [{
            "Coin":"loudcoin",
            "Program": "int(attack * double(randi(2)+20))",
            "Weight":1
        }]
    },
    "ExtraInfo":"",
    "Sender":"eugen",
    "Name": "LOUD's sword sell recipe",
    "CookbookName": "Legend of Undead Dragon",
    "Description": "this recipe is used to sell a sword.",
    "BlockInterval":"0"
}
