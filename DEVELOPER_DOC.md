
## Coins

loudcoin

## Items
- Lv1 wooden sword
    attack: 3
    level: 1
    Name: Wooden sword
- Lv2 wooden sword
    attack: 6
    level: 2
    Name: Wooden sword
- Lv1 copper sword
    attack: 10
    level: 1
    Name: Copper sword
- Lv2 copper sword
    attack: 20
    level: 2
    Name: Copper sword
- Lv1 silver sword
    attack: 30
    level: 1
    Name: Silver sword
- Lv1 bronze sword
    attack: 50
    level: 1
    Name: Bronze sword
- Lv1 iron sword
    attack: 100
    level: 1
    Name: Iron sword
## Recipes

### Buy Wooden sword lv1
Output: Lv1 wooden sword
Price: 100 loudcoin

### Wooden sword lv1 to lv2 upgrade
Output: Lv2 wooden sword
Price: 100 loudcoin
Input item: Lv1 wooden sword

### Buy Copper sword lv1
Output: Lv1 copper sword
Price: 250 loudcoin

### Copper sword lv1 to lv2 upgrade
Output: Lv2 copper sword
Price: 100 loudcoin
Input item: Lv1 copper sword

### Make silver sword

Input item: Goblin ear
Price: 250 loudcoin
Output: Lv1 silver sword

### Make bronze sword

Input item: Wolf tail
Price: 250 loudcoin
Output: Lv1 bronze sword

### Make iron sword

Input item: Troll toes
Price: 250 loudcoin
Output: Lv1 iron sword

### Sword sell recipe, attack * (randi(2)+20) gold
Output: `attack * (randi(2)+20)` loudcoin
Input: Any item which has `attack` and `level` attributes

### Hunt recipe without sword
Reward: 1 or 2
Character die percent: 10%
When character die, no gold is returned.

### Hunt recipe with a sword

Reward: 1 + `attack / 2`
Character dying percent: 5%
Sword lose percent: 5%
When character die, no gold is returned.

### Fight Goblin

Goblin HP: 10
Goblin attack: 1
Reward: 50 loudcoin

character should carry sword to fight goblin.
Total received damage you get from Goblin is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Goblin.

Character Dying percent: `(Total received damage - HP) / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Goblin ear bonus item percent: `HP / Total received damage * 0.1`

When character die, no gold is returned.

### Fight Wolf

Wolf HP: 15
Wolf attack: 3
Reward: 150 loudcoin

character should carry sword to fight wolf.
Total received damage you get from Wolf is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Wolf.

Character Dying percent: `(Total received damage - HP) / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Wolf tail bonus item percent: `HP / Total received damage * 0.1`

When character die, no gold is returned.

### Fight Troll

Troll HP: 20
Troll attack: 5
Reward: 300 loudcoin

character should carry sword to fight troll.
Total received damage you get from Troll is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Troll.

Character Dying percent: `(Total received damage - HP) / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`
Troll toes bonus item percent: `HP / Total received damage * 0.1`

When character die, no gold is returned.

### Fight Giant

Troll HP: 100
Troll attack: 10
Reward: 3000 loudcoin

character should carry iron sword to fight giant.
Total received damage you get from Giant is calculated by using `EnemyHP * EnemyAttack / SwordAttack`, and dying percentage is related to Total received damage you get from Giant.

Character Dying percent: `(Total received damage - HP) / Total received damage`
Character Alive percent: `HP / Total received damage`
Sword lose percent: `HP / Total received damage * 0.1`

When character die, no gold is returned.

### Restore health

HP: +20
Price: 10 loudcoin

If HP is bigger than MaxHP, it is set to MaxHP automatically.
