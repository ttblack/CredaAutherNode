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
	listener          *javaListener.MerkleRootListener
	credaOracle       *credaContract.CredaOracle
	merkleRootChan    chan string
	currentMerkleRoot string
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
	current, err := a.credaOracle.GetMerkleRoot()
	if err != nil {
		log.Println("GetMerkleRoot err", err)
	}
	a.currentMerkleRoot = current.String()

	go a.listener.Start()

	for {
		select {
		case root := <-a.merkleRootChan:
			newRoot := common.HexToHash(root)
			log.Println("new root: ", newRoot, "currentMerkleRoot", a.currentMerkleRoot)
			if newRoot.String() != a.currentMerkleRoot {
				if err := a.credaOracle.SetMerkleRoot(newRoot); err != nil {
					log.Println("SetMerkleRoot err", err)
				} else {
					a.currentMerkleRoot = newRoot.String()
				}
			}

		case <-interceptor.ShutdownChannel():
			a.listener.Stop()
			wg.Done()
			return
		}
	}
}
