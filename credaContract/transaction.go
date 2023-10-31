package credaContract

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func RawWithSignature(key *ecdsa.PrivateKey, chainID *big.Int, transaction *types.Transaction) ([]byte, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}
	tx, err := opts.Signer(crypto.PubkeyToAddress(key.PublicKey), transaction)
	if err != nil {
		return nil, err
	}
	rawTX, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return rawTX, nil
}

func NewTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	return tx
}
