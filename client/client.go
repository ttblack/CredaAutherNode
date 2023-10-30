package client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	rpcClient       *rpc.Client
	contractAddress common.Address
}

func Dial(rawurl string, contractAddress string) (*Client, error) {
	rpcClient, err := rpc.DialContext(context.TODO(), rawurl)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(contractAddress) {
		return nil, errors.New("contract address is not correct")
	}
	address := common.HexToAddress(contractAddress)
	return &Client{
		rpcClient:       rpcClient,
		contractAddress: address,
	}, nil
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := c.rpcClient.CallContext(ctx, &result, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}
