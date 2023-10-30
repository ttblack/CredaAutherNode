package javaListener

import (
	"fmt"
	"time"
)

type MerkleRootListener struct {
	url               string
	merkleRoot_ch     chan string
	currentMerkleRoot string
	delayTimer        *time.Ticker
}

func Create(api string, merkleRootCh chan string, intervalSeconds int64) (*MerkleRootListener, error) {
	fmt.Println("interval", intervalSeconds)
	l := &MerkleRootListener{
		url:               api,
		merkleRoot_ch:     merkleRootCh,
		currentMerkleRoot: "",
		delayTimer:        time.NewTicker(time.Duration(intervalSeconds) * time.Second),
	}
	return l, nil
}

func (l *MerkleRootListener) Start() {
	go func() {
		for {
			select {
			case c := <-l.delayTimer.C:
				l.merkleRoot_ch <- c.String()
			}
		}
	}()
}

func (l *MerkleRootListener) Stop() {
	fmt.Println("MerkleRootListener stop")
	l.delayTimer.Stop()
}
