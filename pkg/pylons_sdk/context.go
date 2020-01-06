package pylonssdk

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func NewCLIContextFromArgs(nodeURI string, cdc, config *Config) (*context.CLIContext, error) {

	rpc := rpcclient.NewHTTP(nodeURI, "/websocket")

	fromAddress, fromName, err := GetFromFields(config.From, DefaultHome)
	if err != nil {
		return nil, err
	}

	return &context.CLIContext{
		NodeURI:     nodeURI,
		Client:      rpc,
		From:        config.From,
		TrustNode:   config.TrustNode,
		FromAddress: fromAddress,
		FromName:    fromName,
		SkipConfirm: true,
	}, nil

}
