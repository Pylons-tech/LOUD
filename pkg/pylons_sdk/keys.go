package pylonssdk

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	cryptokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetFromFields(from string, homeDir string) (sdk.AccAddress, string, error) {
	if from == "" {
		return nil, "", nil
	}

	keybase, err := keys.NewKeyBaseFromDir(homeDir)
	if err != nil {
		return nil, "", err
	}

	var info cryptokeys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			return nil, "", err
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			return nil, "", err
		}
	}

	return info.GetAddress(), info.GetName(), nil
}
