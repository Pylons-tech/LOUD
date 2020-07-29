package loud

import sdk "github.com/cosmos/cosmos-sdk/types"

// LockedCoinDescribe describes the locked coin struct
type LockedCoinDescribe struct {
	ID     string
	Amount sdk.Coins
	Type   string // trade | recipe
}
