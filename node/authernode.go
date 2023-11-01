package node

import (
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ttblack/CredaAutherNode/config"
	"github.com/ttblack/CredaAutherNode/credaContract"
	"github.com/ttblack/CredaAutherNode/javaListener"
	"github.com/ttblack/CredaAutherNode/signal"
)

type AutherNode struct {
	listener       *javaListener.MerkleRootListener
	credaOracle    *credaContract.CredaOracle
	merkleRootChan chan string
}

func New(cfg *config.Config) (*AutherNode, error) {
	node := &AutherNode{
		merkleRootChan: make(chan string),
	}
	listener, err := javaListener.Create(cfg.MerkleRootAPI, node.merkleRootChan, cfg.MerkleRootListenerInterval)
	if err != nil {
		return nil, err
	}
	node.listener = listener

	oracle, err := credaContract.New(cfg.ChainRPC, cfg.CredaOracle, cfg.AutherKeystore, cfg.AutherPassword)
	if err != nil {
		return nil, err
	}
	node.credaOracle = oracle

	return node, nil
}

func (a *AutherNode) Start(wg *sync.WaitGroup, interceptor *signal.Interceptor) {
	wg.Add(1)
	go a.listener.Start()

	for {
		select {
		case root := <-a.merkleRootChan:
			hash := common.HexToHash(root)
			if err := a.credaOracle.SetMerkleRoot(hash); err != nil {
				log.Printf("SetMerkleRoot err: %v", err)
			}
		case <-interceptor.ShutdownChannel():
			a.listener.Stop()
			wg.Done()
			return
		}
	}
}
