package service

import (
	"bytes"
	"encoding/json"
	config "miborchestrator/configs"
	"miborchestrator/internal/entities"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const maxAttempts = 10

type MyProvider struct {
	Client *http.Client
	Config *config.Config
}

func NewProvider(cfg *config.Config) *MyProvider {
	client := &http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	client.Transport = transport
	client.Timeout = time.Second * 1
	return &MyProvider{
		Client: client,
		Config: cfg,
	}
}

func (p *MyProvider) shardNum(walletID int) int {
	return walletID % len(p.Config.Shards)
}

func (p *MyProvider) createWalletRequest(reqBody entities.WcTDO) {

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		logrus.Println(err)
	}

	bodyReader := bytes.NewReader(reqJSON)

	response, err := p.Client.Post(p.Config.Shards[p.shardNum(reqBody.WalletID)]+"api/createwallet", "application/json", bodyReader)
	if err != nil {
		logrus.Println(err)
		logrus.Println(response.StatusCode)
	}
	logrus.Println(response.StatusCode)

}

func (p *MyProvider) sendPreRequest(reqData entities.PrepareTransactionRequest) {
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		logrus.Printf("service_level:transaction_manager:prepareReader:Marshalerror:%v", err)
	}

	bodyReader := bytes.NewReader(reqJSON)

	var SuccessfulRequest bool
	var attempts int

	for !SuccessfulRequest || attempts > maxAttempts {
		var attempts int
		response, err := p.Client.Post(p.Config.Shards[p.shardNum(reqData.WalletID)]+"api/prepare", "application/json", bodyReader)
		attempts++
		if err != nil {
			logrus.Printf("servicelevel:service_level:transaction_manager: sending request error:%v", err)
		}
		if response.StatusCode == http.StatusAccepted {
			SuccessfulRequest = true
		}
	}
	reqData.ResponseChan <- SuccessfulRequest

}
