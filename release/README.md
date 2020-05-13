# Test result

LOUD v0.1.0 app is tested on Mojave (10.14.5), High Sierra(10.13.6), and Sierra(10.12.6).

# Upgrades from v0.0.5
1. Game items and characters modification
- Added special bonus characters like fire bonus, ice bonus, acid bonus, special dragons and undead dragon, angel sword
- Removed HP from recipes and set fixed % of dying per monsters
- Updated level upgrade mechanism
- Giant gives bonus skill by chance
- Bonus character can't fight giant

2. UI modifications
- Tab UI navigation
- Layout modification
- Better game text
- Add more visualizations using unicodes and colorful messages like green, red, orange
- Inventory overflow issue fix

3. Fix unexpected issues
- Fixed issue reported by @Gopher JK, non-rabbit sword lose & no active weapon case result description issue
- Fixed other issues that was found during development

4. Enhancement of app reporting system
- Better logging system to get specific error logs from users
- Make users to be able to change configuration at Package content (e.g. rest and cli endpoint and maximum wait block)


# How to run

- download from github release or trusted source
- Unzip the file
- copy to non-download folder where terminal can do read/write
- when opening first time, right click on app and select open - and it will show a popup and you can click open
- from the second time, it will be okay to double click the app to open
- It will open a terminal, and you can just enter your username to start game!
- That's it, enjoy LOUD game

# How to run on terminal

- Install and run
```
cd $LOUD_APP_DIR/Contents/Resources/
./install_and_run.sh ${player_name}
```

- Run
```
cd $LOUD_APP_DIR/Contents/Resources/
./bin/loud ${player_name}
```

# How to test

To test things easily

- Go to Develop (D)
- Get goblin's ear, troll's toes, wolf's tail (B)

- Go to Pylons Central(C)
- Go to buy character (1)
- Buy Tiger character using (1)
- Go back to Pylons Central Entry point by pressing Backspace
- Buy gold with pylons (2)

- Go to Shop by pressing (S)
- Go to Buy Items(1)
- Buy all the swords by pressing (1-5)

- Go to Home (H)
- Go to room for updating character's name (4)
- Select character to rename (arrow keys | number)
- Enter new name and press Enter key

To sell item (wooden sword and copper sword) and upgrade item (wooden sword and copper sword) can be tested under shop.

After that can go to Forest (F)
In forest, you can try to Hunt Rabbits(1), Fight Goblin, Wolf, Troll or Giant 🗿 (2-5)
To fight Giant, you need to carry iron sword.
To fight other monsters, need to carry at least 1 sword.

After fight goblin, wolf and troll, you can get bonus item as a chance.
And they can be used when creating sword like Silver, bronze and iron swords.

If you want to trade with items and coins, you can just place order or just fulfill the order under Pylons Central(M).