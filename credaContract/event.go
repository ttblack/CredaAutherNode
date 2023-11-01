package credaContract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"
)

func SetMerkleRoot() abi.ABI {
	definition := "[{\n      \"inputs\": [\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"_merkleRoot\",\n          \"type\": \"bytes32\"\n        }\n      ],\n      \"name\": \"setMerkleRoot\",\n      \"outputs\": [],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"function\"\n    }]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		log.Error("SetMerkleRoot failed", "error", err)
		return a
	}

	return a
}
