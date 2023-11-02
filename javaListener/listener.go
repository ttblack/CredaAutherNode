package javaListener

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MerkleRootListener struct {
	url               string
	merkleRoot_ch     chan string
	currentMerkleRoot string
	delayTimer        *time.Ticker
}

func Create(api string, merkleRootCh chan string, intervalSeconds int64) (*MerkleRootListener, error) {
	l := &MerkleRootListener{
		url:               api,
		merkleRoot_ch:     merkleRootCh,
		currentMerkleRoot: "",
		delayTimer:        time.NewTicker(time.Duration(intervalSeconds) * time.Second),
	}
	return l, nil
}

type MerkleRootData struct {
	DateRef string `json:"dateRef"`
	Root    string `json:"root"`
}

type GetMerkleRootResult struct {
	Code    int32           `json:"code"`
	Message string          `json:"message"`
	Data    *MerkleRootData `json:"data"`
}

func (l *MerkleRootListener) httpGetMerkleRoot() (string, error) {
	client := &http.Client{}
	log.Println("url", l.url)
	// 创建一个GET请求
	req, err := http.NewRequest("GET", l.url, nil)
	if err != nil {
		return "", errors.New(fmt.Sprintf("new http request err: %v", err))
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("send http request err: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("read http respond err: %v", err))
	}

	// 打印响应内容
	log.Println("http respond body", string(body))

	var result GetMerkleRootResult
	err = json.Unmarshal(body, &result)
	if err != nil || result.Data == nil {
		return "", errors.New(fmt.Sprintf("unmarshal http respond err: %v", err))
	}

	log.Println("root", result.Data.Root)

	return result.Data.Root, nil
}

func (l *MerkleRootListener) Start() {
	go func() {
		for {
			select {
			case <-l.delayTimer.C:
				root, err := l.httpGetMerkleRoot()
				if err != nil {
					log.Println("httpGetMerkleRoot err", err)
				} else {
					l.merkleRoot_ch <- root
				}
			}
		}
	}()
}

func (l *MerkleRootListener) Stop() {
	fmt.Println("MerkleRootListener stop")
	l.delayTimer.Stop()
}
