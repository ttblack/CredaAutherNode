package node

import (
	"fmt"
	"sync"

	"github.com/ttblack/CredaAutherNode/config"
	"github.com/ttblack/CredaAutherNode/javaListener"
	"github.com/ttblack/CredaAutherNode/signal"
)

type AutherNode struct {
	listener       *javaListener.MerkleRootListener
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
	return node, nil
}

func (a *AutherNode) Start(wg *sync.WaitGroup, interceptor *signal.Interceptor) {
	wg.Add(1)
	go a.listener.Start()

	for {
		select {
		case root := <-a.merkleRootChan:
			fmt.Println("root", root)
		case <-interceptor.ShutdownChannel():
			a.listener.Stop()
			wg.Done()
			return
		}
	}
}
