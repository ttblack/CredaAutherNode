package credaContract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ttblack/CredaAutherNode/client"
)

type CredaOracle struct {
	rpcClient *client.Client
}

func New(rpc, contractAddress string) (*CredaOracle, error) {
	cli, err := client.Dial(rpc, contractAddress)
	if err != nil {
		return nil, err
	}
	oracle := &CredaOracle{
		rpcClient: cli,
	}
	return oracle, nil
}

func (c *CredaOracle) SetMerkleRoot(merkleRoot common.Hash) {

}
