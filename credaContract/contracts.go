package credaContract

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ttblack/CredaAutherNode/client"
)

type CredaOracle struct {
	rpcClient *client.Client
	key       *keystore.Key
}

func New(rpc, contractAddress, keystorePath, password string) (*CredaOracle, error) {
	cli, err := client.Dial(rpc, contractAddress)
	if err != nil {
		return nil, err
	}
	if keystorePath == "" {
		return nil, errors.New("keystorePath is nil")
	}
	file, err := os.OpenFile(keystorePath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(file)
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, err
	}
	oracle := &CredaOracle{
		rpcClient: cli,
		key:       key,
	}
	return oracle, nil
}

func (c *CredaOracle) SetMerkleRoot(merkleRoot common.Hash) (common.Hash, error) {
	a := SetMerkleRoot()
	input, err := a.Pack("setMerkleRoot", merkleRoot)
	if err != nil {
		return common.Hash{}, errors.New(fmt.Sprintf("pack param err: %v", err))
	}

	hash, err := c.makeAndSendContractTransaction(input, big.NewInt(0))
	fmt.Println("hash", hash.String(), "err", err)
	if err != nil {
		return common.Hash{}, errors.New(fmt.Sprintf("make and send tx err: %v", err))
	}

	log.Info("SetMerkleRoot hash: %v", hash)
	return hash, nil
}

func (c *CredaOracle) makeAndSendContractTransaction(data []byte, value *big.Int) (common.Hash, error) {
	var hash common.Hash
	from := c.key.Address

	ctx := context.Background()
	gasPrice, err := c.rpcClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Error("SuggestGasPrice", "error", err)
		return hash, err
	}
	msg := ethereum.CallMsg{From: from, To: c.rpcClient.GetContractAddress(), Data: data, GasPrice: gasPrice}
	gasLimit, err := c.rpcClient.EstimateGas(ctx, msg)
	if err != nil {
		log.Error("EstimateGas", "error", err)
		return hash, err
	}
	if gasLimit == 0 {
		return hash, errors.New("EstimateGasLimit is 0")
	}
	nonce, err := c.rpcClient.PendingNonceAt(ctx, from)
	if err != nil {
		log.Error("PendingNonceAt", "error", err)
		return hash, err
	}

	tx := NewTransaction(nonce, c.rpcClient.GetContractAddress(), value, gasLimit, gasPrice, data)
	log.Info("makeAndSendContractTransaction", "gasLimit", gasLimit, "gasPrice", gasPrice, "nonce", nonce)
	return c.SignAndSendTransaction(ctx, tx)
}

func (c *CredaOracle) SignAndSendTransaction(ctx context.Context, tx *types.Transaction) (common.Hash, error) {
	id, err := c.rpcClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	rawTX, err := RawWithSignature(c.key.PrivateKey, id, tx)
	if err != nil {
		return common.Hash{}, err
	}

	return c.rpcClient.SendRawTransaction(ctx, rawTX)
}
