package credaContract

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"syscall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ttblack/CredaAutherNode/client"
	"golang.org/x/term"
)

type CredaOracle struct {
	rpcClient *client.Client
	key       *keystore.Key
}

func New(rpc, contractAddress, keystorePath string) (*CredaOracle, error) {
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

	fmt.Print("Enter Keystore Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("read password for keystore err: %v", err))
	}
	password := string(bytePassword)

	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, err
	}
	oracle := &CredaOracle{
		rpcClient: cli,
		key:       key,
	}

	fmt.Println("")

	return oracle, nil
}

func (c *CredaOracle) SetMerkleRoot(merkleRoot common.Hash) error {
	a := SetMerkleRoot()
	input, err := a.Pack("setMerkleRoot", merkleRoot)
	if err != nil {
		return errors.New(fmt.Sprintf("pack param err: %v", err))
	}

	hash, err := c.makeAndSendContractTransaction(input, big.NewInt(0))
	if err != nil {
		return errors.New(fmt.Sprintf("make and send tx err: %v", err))
	}

	log.Println("SetMerkleRoot OK, tx hash", hash.String())
	return nil
}

func (c *CredaOracle) GetMerkleRoot() (common.Hash, error) {
	a := GetMerkleRoot()
	input, err := a.Pack("getRoot")
	if err != nil {
		return common.Hash{}, errors.New(fmt.Sprintf("pack param err: %v", err))
	}

	msg := ethereum.CallMsg{From: common.Address{}, To: c.rpcClient.GetContractAddress(), Data: input}

	out, err := c.rpcClient.CallContract(context.TODO(), msg, nil)
	if err != nil || len(out) <= 0 {
		return common.Hash{}, errors.New(fmt.Sprintf("CallContract output empty or err: %v", err))
	}

	hash := common.BytesToHash(out)
	return hash, nil
}

func (c *CredaOracle) makeAndSendContractTransaction(data []byte, value *big.Int) (common.Hash, error) {
	var hash common.Hash
	from := c.key.Address

	ctx := context.Background()
	gasPrice, err := c.rpcClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Printf("SuggestGasPrice err: %v", err)
		return hash, err
	}
	msg := ethereum.CallMsg{From: from, To: c.rpcClient.GetContractAddress(), Data: data, GasPrice: gasPrice}
	gasLimit, err := c.rpcClient.EstimateGas(ctx, msg)
	if err != nil {
		log.Printf("EstimateGas err: %v", err)
		return hash, err
	}
	if gasLimit == 0 {
		return hash, errors.New("EstimateGasLimit is 0")
	}
	nonce, err := c.rpcClient.PendingNonceAt(ctx, from)
	if err != nil {
		log.Printf("PendingNonceAt err: %v", err)
		return hash, err
	}

	tx := NewTransaction(nonce, c.rpcClient.GetContractAddress(), value, gasLimit, gasPrice, data)
	log.Println("makeAndSendContractTransaction", "gasLimit", gasLimit, "gasPrice", gasPrice, "nonce", nonce)
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
