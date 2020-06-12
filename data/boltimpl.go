package loud

import (
	"strings"

	"github.com/Pylons-tech/LOUD/log"
	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
	bolt "github.com/coreos/bbolt"
)

type dbWorld struct {
	filename string
	database *bolt.DB
}

func (w *dbWorld) GetUser(username string) User {
	return getUserFromDB(w, username)
}

func (w *dbWorld) newUser(username string) UserData {
	userData := UserData{
		Username: username,
		Location: Home,
		Gold:     0,
		Items:    []Item{},
	}

	return userData
}

func (w *dbWorld) Close() {
	if w.database != nil {
		w.database.Close()
	}
}

func (w *dbWorld) load() {
	log.WithFields(log.Fields{
		"func_start": "dbWorld.load",
		"file_name":  w.filename,
	}).Debugf("function start")

	db, err := bolt.Open(w.filename, 0600, nil)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warnln("error reading database file")
	} else {
		// Make default tables
		err = db.Update(func(tx *bolt.Tx) error {
			buckets := []string{"users"}

			for _, bucket := range buckets {
				_, err := tx.CreateBucketIfNotExists([]byte(bucket))

				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warnln("update database error")
		}
	}
	w.database = db

	log.WithFields(log.Fields{
		"func_end":  "dbWorld.load",
		"file_name": w.filename,
	}).Debugln("function end")
}

// LoadWorldFromDB will set up an on-disk based world
func LoadWorldFromDB(filename string) World {
	newWorld := dbWorld{filename: filename}
	newWorld.load()
	return &newWorld
}

// UserData is a JSON-serializable set of information about a User.
type UserData struct {
	Gold            int
	PylonAmount     int
	Username        string `json:""`
	Address         string
	Location        UserLocation
	Items           []Item
	Characters      []Character
	ActiveCharacter Character
	DeadCharacter   Character
	PrivKey         string
	TargetMonster   string
	UsingWeapon     Item
	lastTransaction string
	lastTxMetaData  string
	lastUpdate      int64
}

type dbUser struct {
	UserData
	world *dbWorld
}

func (user *dbUser) GetPrivKey() string {
	return user.UserData.PrivKey
}

func (user *dbUser) GetLocation() UserLocation {
	return user.UserData.Location
}

func (user *dbUser) SetLocation(loc UserLocation) {
	user.UserData.Location = loc
}

func (user *dbUser) Reload() {
	log.WithFields(log.Fields{
		"func_start": "user.Reload",
	}).Debugln("function start")
	var record []byte
	if user.world.database != nil {
		err := user.world.database.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("users"))
			record = bucket.Get([]byte(user.UserData.Username))

			return nil
		})
		if err != nil {
			log.WithFields(log.Fields{
				"db_error": "db.View",
				"error":    err,
			}).Warnln("db error")
		}
	}

	if record == nil {
		log.WithFields(log.Fields{
			"user_not_exist": user.UserData.Username,
		}).Info("user does not exist, creating a new one ...")
		user.UserData = user.world.newUser(user.UserData.Username)
		user.Save()
	} else {
		err := MSGUnpack(record, &(user.UserData))
		if err != nil {
			log.Warnln("Was not able to unpack data record into user data")
		}
		log.WithFields(log.Fields{
			"loaded_user": pylonSDK.JSONFormatter(user.UserData),
		}).Infof("loaded user")
		user.FixLoadedData()
	}
	user.UserData.PrivKey = InitPylonAccount(user.UserData.Username)

	// Initial Sync
	SyncFromNode(user)

	log.WithFields(log.Fields{
		"func_end": "user.Reload",
	}).Debugln("function end")
}

func (user *dbUser) Save() {
	bytes, err := MSGPack(user.UserData)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warnf("userdata MSGPack error")
		return
	}

	if user.world.database != nil {
		err = user.world.database.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("users"))

			err = bucket.Put([]byte(user.UserData.Username), bytes)

			return err
		})
		if err != nil {
			log.WithFields(log.Fields{
				"db_error": "db.Update",
			}).Warnln("update world database error")
		}
	}
}

func (user *dbUser) GetAddress() string {
	return user.UserData.Address
}

func (user *dbUser) SetAddress(addr string) {
	user.UserData.Address = addr
}

func (user *dbUser) GetUserName() string {
	return user.UserData.Username
}

func (user *dbUser) SetGold(amount int) {
	user.UserData.Gold = amount
}
func (user *dbUser) GetGold() int {
	return user.UserData.Gold
}

func (user *dbUser) GetPylonAmount() int {
	return user.UserData.PylonAmount
}

func (user *dbUser) SetPylonAmount(amount int) {
	user.UserData.PylonAmount = amount
}

func (user *dbUser) SetItems(items []Item) {
	user.UserData.Items = items
}

func (user *dbUser) SelectFightWeapon() {
	switch user.TargetMonster {
	case TextRabbit: // no weapon is needed
		user.UsingWeapon = Item{}
	case TextGoblin, TextWolf, TextTroll: // any sword is ok
		weapons := user.InventorySwords()
		if len := len(weapons); len > 0 {
			user.UsingWeapon = weapons[len-1]
		}
	case TextGiant, TextDragonFire, TextDragonIce, TextDragonAcid: // iron sword is needed
		weapons := user.InventoryIronSwords()
		if len := len(weapons); len > 0 {
			user.UsingWeapon = weapons[len-1]
		}
	case TextDragonUndead: // angel sword is needed
		weapons := user.InventoryAngelSwords()
		if len := len(weapons); len > 0 {
			user.UsingWeapon = weapons[len-1]
		}
	default:
		user.UsingWeapon = Item{}
	}
}

func (user *dbUser) SetFightMonster(monster string) {
	user.TargetMonster = monster
	user.SelectFightWeapon()
}

func (user *dbUser) GetTargetMonster() string {
	return user.TargetMonster
}

func (user *dbUser) GetItemByID(ID string) *Item {
	iis := user.InventoryItems()
	for _, ii := range iis {
		if ii.ID == ID {
			return &ii
		}
	}
	return nil
}

func (user *dbUser) GetFightWeapon() *Item {
	if user.UsingWeapon.Name == "" {
		return nil
	}
	return &user.UsingWeapon
}

func (user *dbUser) SetCharacters(items []Character) {
	user.UserData.Characters = items
}

func (user *dbUser) SetActiveCharacterIndex(idx int) {
	length := len(user.UserData.Characters)
	if idx >= 0 && idx < length {
		user.UserData.ActiveCharacter = user.UserData.Characters[idx]
	} else {
		if len(user.ActiveCharacter.ID) > 0 {
			user.UserData.DeadCharacter = user.UserData.ActiveCharacter
		}
		user.UserData.ActiveCharacter = Character{}
	}
}

func (user *dbUser) FixLoadedData() {
	if len(user.ActiveCharacter.ID) > 0 {
		length := len(user.UserData.Characters)
		idx := user.GetActiveCharacterIndex()
		if idx >= 0 && idx < length {
			// update to latest character status
			user.UserData.ActiveCharacter = user.UserData.Characters[idx]
		} else {
			// it means old active character is dead
			user.UserData.DeadCharacter = user.UserData.ActiveCharacter
			user.UserData.ActiveCharacter = Character{}
		}
	}
}

func (user *dbUser) GetActiveCharacterIndex() int {
	// order can be changed on every sync so we should take care of the index by finding by id
	for idx, ic := range user.InventoryCharacters() {
		if ic.ID == user.ActiveCharacter.ID {
			return idx
		}
	}
	return -1
}

func (user *dbUser) GetActiveCharacter() *Character {
	if user.UserData.ActiveCharacter.Name == "" {
		return nil
	}
	return &user.UserData.ActiveCharacter
}

func (user *dbUser) GetDeadCharacter() *Character {
	if user.UserData.DeadCharacter.Name == "" {
		return nil
	}
	return &user.UserData.DeadCharacter
}

func (user *dbUser) InventoryItems() []Item {
	return user.UserData.Items
}

func (user *dbUser) HasPreItemForAnItem(item Item) bool {
	if len(item.PreItems) == 0 {
		return true
	}
	for _, pi := range item.PreItems {
		if len(user.InventoryItemIDByName(pi)) == 0 {
			return false
		}
	}
	return true
}

func (user *dbUser) InventoryItemIDByName(name string) string {
	iis := user.InventoryItems()
	for _, ii := range iis {
		if strings.EqualFold(ii.Name, name) {
			return ii.ID
		}
	}
	return ""
}

func (user *dbUser) InventoryAngelSwords() []Item {
	iis := user.InventoryItems()
	uis := []Item{}
	for _, ii := range iis {
		if ii.Name == AngelSword {
			uis = append(uis, ii)
		}
	}
	return uis
}

func (user *dbUser) InventoryIronSwords() []Item {
	iis := user.InventoryItems()
	uis := []Item{}
	for _, ii := range iis {
		if ii.Name == IronSword {
			uis = append(uis, ii)
		}
	}
	return uis
}

func (user *dbUser) InventorySwords() []Item {
	iis := user.InventoryItems()
	uis := []Item{}
	for _, ii := range iis {
		if ii.IsSword() {
			uis = append(uis, ii)
		}
	}
	return uis
}

func (user *dbUser) InventoryCharacters() []Character {
	return user.UserData.Characters
}

func (user *dbUser) InventoryUpgradableItems() []Item {
	iis := user.InventoryItems()
	uis := []Item{}
	for _, ii := range iis {
		if ii.Level == 1 && (ii.Name == CopperSword || ii.Name == WoodenSword) {
			uis = append(uis, ii)
		}
	}
	return uis
}

func (user *dbUser) InventorySellableItems() []Item {
	return user.InventoryItems()
}

func (user *dbUser) GetLastTxHash() string {
	return user.UserData.lastTransaction
}

func (user *dbUser) GetLastTxMetaData() string {
	return user.UserData.lastTxMetaData
}

func (user *dbUser) SetLastTransaction(tx, metadata string) {
	user.UserData.lastTransaction = tx
	user.UserData.lastTxMetaData = metadata
}

func (user *dbUser) SetLatestBlockHeight(h int64) {
	user.UserData.lastUpdate = h
}

func (user dbUser) GetLatestBlockHeight() int64 {
	return user.UserData.lastUpdate
}

func (user dbUser) GetMatchedItems(itspec ItemSpec) []Item {
	mitems := []Item{}
	for _, item := range user.InventoryItems() {
		if len(itspec.Name) != 0 && item.Name != itspec.Name {
			continue
		}
		if itspec.Attack[1] > 0 && (item.Attack < itspec.Attack[0] || item.Attack > itspec.Attack[1]) {
			continue
		}
		if itspec.Level[1] > 0 && (item.Level < itspec.Level[0] || item.Level > itspec.Level[1]) {
			continue
		}
		mitems = append(mitems, item)
	}
	return mitems
}

func (user dbUser) GetMatchedCharacters(chspec CharacterSpec) []Character {
	mchars := []Character{}
	for _, char := range user.InventoryCharacters() {
		if chspec.Special > 0 && char.Special != chspec.Special {
			continue
		}
		if len(chspec.Name) != 0 && char.Name != chspec.Name {
			continue
		}
		if chspec.XP[1] > 0 && (char.XP < chspec.XP[0] || char.XP > chspec.XP[1]) {
			continue
		}
		if chspec.Level[1] > 0 && (char.Level < chspec.Level[0] || char.Level > chspec.Level[1]) {
			continue
		}
		mchars = append(mchars, char)
	}
	return mchars
}

func getUserFromDB(world *dbWorld, username string) User {
	user := dbUser{
		UserData: UserData{
			Username: username,
		},
		world: world,
	}

	user.Reload()

	return &user
}
