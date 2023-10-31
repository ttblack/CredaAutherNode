package credaContract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"
)

func SetMerkleRoot() abi.ABI {
	definition := "{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"setMerkleRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		log.Error("SetMerkleRoot failed", "error", err)
		return a
	}

	return a
}
