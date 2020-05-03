
## Coins

loudcoin

## Items
- Lv1 wooden sword
```
    attack: 3
    level: 1
    Name: Wooden sword
```
- Lv2 wooden sword
```
    attack: 6
    level: 2
    Name: Wooden sword
```
- Lv1 copper sword
```
    attack: 10
    level: 1
    Name: Copper sword
```
- Lv2 copper sword
```
    attack: 20
    level: 2
    Name: Copper sword
```
- Lv1 silver sword
```
    attack: 30
    level: 1
    Name: Silver sword
```
- Lv1 bronze sword
```
    attack: 50
    level: 1
    Name: Bronze sword
```
- Lv1 iron sword
```
    attack: 100
    level: 1
    Name: Iron sword
```
- Lv1 angel sword
```
    attack: 1000
    level: 1
    Name: Angel sword
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
10% chance of character lose
When character die, no gold is returned.
```

### Hunt rabbits recipe with a sword

```
Reward: 1 + `attack / 2`
5% chance of character lose
5% chance of sword lose
When character die, no gold or sword is returned.
```

### Fight Goblin

```
Goblin HP: 10
Goblin attack: 1
Reward: 50 loudcoin
10% chance of sword lose
10% chance of "Goblin ear"
```

character should carry sword to fight goblin.
Total received damage you get from Goblin is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Goblin.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Goblin ear bonus item percent: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight Wolf

```
Wolf HP: 15
Wolf attack: 3
Reward: 150 loudcoin
10% chance of sword lose
10% chance of “Wolf tail”
```
character should carry sword to fight wolf.
Total received damage you get from Wolf is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Wolf.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Wolf tail bonus item percent: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight Troll

```
Troll HP: 20
Troll attack: 5
Reward: 300 loudcoin
10% chance of sword lose
10% chance of “Troll toes”
```

character should carry sword to fight troll.
Total received damage you get from Troll is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Troll.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Troll toes bonus item percent: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight Giant

Warn. Character with bonus skill can't fight Giant.

```
Giant HP: 100
Giant attack: 10
Reward: 3000 loudcoin
10% chance of sword lose
GiantKiller badget on character
4% chance of fire bonus skill
3% chance of ice bonus skill
2% chance of acid bonus skill
```

character should carry iron sword to fight giant.
Total received damage you get from Giant is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Giant.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight fire dragon

```
Fire Dragon HP: 300
Fire Dragon attack: 30
Reward: 10000 loudcoin
10% chance of sword lose
10% chance of “Fire scale” - drop from fire dragon
FireDragonKiller badget on character
```

character should carry iron sword to fight Fire dragon.
Total received damage you get from Fire dragon is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Fire dragon.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Drop from fire dragon: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight ice dragon

```
Ice Dragon HP: 300
Ice Dragon attack: 30
Reward: 10000 loudcoin
10% chance of sword lose
10% chance of “Icy shards” - drop from ice dragon
IceDragonKiller badget on character
```

character should carry iron sword to fight Ice dragon.
Total received damage you get from Ice dragon is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Ice dragon.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Drop from ice dragon: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight acid dragon

```
Acid Dragon HP: 300
Acid Dragon attack: 30
Reward: 10000 loudcoin
10% chance of sword lose
10% chance of “poison claws” - drop from acid dragon
AcidDragonKiller badget on character
```

character should carry iron sword to fight Acid dragon.
Total received damage you get from Acid dragon is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Acid dragon.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Drop from acid dragon: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Fight undead dragon

```
Undead Dragon HP: 1000
Undead Dragon attack: 100
Reward: 50000 loudcoin
10% chance of sword lose
UndeadDragonKiller badget on character
```

character should carry iron sword to fight Undead dragon.
Total received damage you get from Undead dragon is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Undead dragon.

```
Character Dying percent: `1 - HP / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Drop from Undead dragon: `HP / Total received damage * 0.1`
```

When character die, no gold is returned.

### Restore health

```
HP: +20
Price: 10 loudcoin
```

If HP is bigger than MaxHP, it is set to MaxHP automatically.

Warn: When a character fight or hunt rabbits, automatic health increaser just work to restore the health gained during relax.
For every block, the HP is increased by +1.