
## Coins

Pylon: 🔷
Gold: 💰
It is described as loudcoin denom.

## Items
1. Wooden sword
- Lv1
```
    attack: 3
```
- Lv2
```
    attack: 6
```
2. Copper sword
- Lv1
```
    attack: 10
```
- Lv2
```
    attack: 20
```
3. Silver sword
- Lv1
```
    attack: 30
```
4. Bronze sword
- Lv1
```
    attack: 50
```
5. Iron sword
- Lv1
```
    attack: 100
```
6. Angel sword
- Lv1
```
    attack: 1000
```

## Recipes

### Buy Wooden sword lv1
```
Output: Lv1 wooden sword
Price: 100 loudcoin
```

### Wooden sword lv1 to lv2 upgrade
```
Output: Lv2 wooden sword
Price: 100 loudcoin
Input item: Lv1 wooden sword
```

### Buy Copper sword lv1
```
Output: Lv1 copper sword
Price: 250 loudcoin
```
### Copper sword lv1 to lv2 upgrade
```
Output: Lv2 copper sword
Price: 100 loudcoin
Input item: Lv1 copper sword
```

### Make silver sword
```
Input item: Goblin ear
Price: 250 loudcoin
Output: Lv1 silver sword
```

### Make bronze sword

```
Input item: Wolf tail
Price: 250 loudcoin
Output: Lv1 bronze sword
```

### Make iron sword

```
Input item: Troll toes
Price: 250 loudcoin
Output: Lv1 iron sword
```

### Make angel sword
```
Input item: Drops from 3 special dragons; fire dragon, acid dragon, ice dragon
Price: 20000 loudcoin
Output: Lv1 angel sword
```

### Sword sell recipe, attack * (randi(2)+20) gold
```
Output: `attack * (randi(2)+20)` loudcoin
Input: Any item which has `attack` and `level` attributes
```

### Hunt rabbits recipe without sword
```
Reward: 1 or 2
5% chance of character dying
When character die, no gold is returned.
```

### Hunt rabbits recipe with a sword

```
Reward: 1 + `attack / 2`
2% chance of character dying
3% chance of sword lose
When character die, no gold or sword is returned.
```

### Fight Goblin 👺

```
Goblin HP: 10
Goblin attack: 1
Reward: 50 loudcoin
2% chance of character dying
3% chance of sword lose
10% chance of "Goblin ear"
```

Character should carry sword to fight goblin.
When character die, no gold is returned.

### Fight Wolf 🐺

```
Wolf HP: 15
Wolf attack: 3
Reward: 150 loudcoin
3% chance of character dying
3% chance of sword lose
10% chance of “Wolf tail”
```
Character should carry sword to fight wolf.
When character die, no gold is returned.

### Fight Troll 👻

```
Troll HP: 20
Troll attack: 5
Reward: 300 loudcoin
4% chance of character dying
3% chance of sword lose
10% chance of “Troll toes”
```

Character should carry sword to fight troll.
When character die, no gold is returned.

### Fight Giant 🗿

Warn. Character with bonus skill can't fight Giant.

```
Giant HP: 100
Giant attack: 10
Reward: 3000 loudcoin
GiantKiller badget on character
5% chance of character dying
3% chance of sword lose
4% chance of fire bonus skill
3% chance of ice bonus skill
3% chance of acid bonus skill
```

Character should carry iron sword to fight giant.
When character die, no gold is returned.

### Fight fire dragon 🦐

```
Fire Dragon HP: 300
Fire Dragon attack: 30
Reward: 10000 loudcoin
2% chance for character dying
3% chance of sword lose
10% chance of “Fire scale” - drop from fire dragon
FireDragonKiller badget on character
```

Character should carry iron sword to fight Fire dragon.
When character die, no gold is returned.

### Fight ice dragon 🦈

```
Ice Dragon HP: 300
Ice Dragon attack: 30
Reward: 10000 loudcoin
2% chance for character dying
3% chance of sword lose
10% chance of “Icy shards” - drop from ice dragon
IceDragonKiller badget on character
```

Character should carry iron sword to fight Ice dragon.
When character die, no gold is returned.

### Fight acid dragon 🐊

```
Acid Dragon HP: 300
Acid Dragon attack: 30
Reward: 10000 loudcoin
2% chance for character dying
3% chance of sword lose
10% chance of “poison claws” - drop from acid dragon
AcidDragonKiller badget on character
```

Character should carry iron sword to fight Acid dragon.
When character die, no gold is returned.

### Fight undead dragon 🐉

```
Undead Dragon HP: 1000
Undead Dragon attack: 100
Reward: 50000 loudcoin
2% chance of character dying
3% chance of sword lose
UndeadDragonKiller badget on character
```

Character should carry iron sword to fight Undead dragon.
When character die, no gold is returned.

### Level and XP

Character gets 1 XP when hunting rabbit.
Character gets `Enemy HP * Enemy attack` when fighting monsters.

Level is upgraded in this mechanism on each fight.
```
Level = Level + XP/(level^3 + 5)
```
