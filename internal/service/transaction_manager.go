package service

import (
	"context"
	"math/rand"
	"miborchestrator/internal/entities"
	"sync"
	"time"
)

type Provider interface {
	sendPreRequest(entities.PrepareTransactionRequest)
	createWalletRequest(entities.WcTDO)
}
type TxStatus int

const (
	StatusProcessing TxStatus = iota
	StatusCommited
	StatusSuccessful
	StatusFailed
)

type TransactionManager struct {
	provider Provider
	stream   chan entities.Request
	txLog    map[string]*entities.Transaction
	wg       sync.WaitGroup
}

func NewTransactionManager(ctx context.Context, provider Provider) *TransactionManager {
	tm := &TransactionManager{
		stream:   make(chan entities.Request, 100),
		txLog:    make(map[string]*entities.Transaction),
		provider: provider,
	}
	tm.Run(ctx)
	return tm
}

func (tm *TransactionManager) Run(ctx context.Context) {
	go tm.listen(ctx)
}

func (tm *TransactionManager) Close() {
	close(tm.stream)
}

func (tm *TransactionManager) listen(ctx context.Context) {

	for {

		select {
		case req := <-tm.stream:
			tm.distributeRequest(ctx, req)

		case <-ctx.Done():
			tm.wg.Wait()
			return
		}
	}
}

func (tm *TransactionManager) distributeRequest(ctx context.Context, req entities.Request) {
	switch req.Event {
	case entities.InitTx:
		tx := tm.createTransaction(req.WalletIDFrom, req.WalletIDTo, req.Amount)
		tm.txLog[tx.TXID] = tx
		go tm.Transfer(ctx, tx)
	case entities.PreparationsConfirmed:
		//..
	case entities.CommitConfitmed:
		//..
	case entities.TxSuccessfull:
		//..
	}
}
func (tm *TransactionManager) createTransaction(fromID, toID, amount int) *entities.Transaction {
	return &entities.Transaction{
		TXID:             generateID(),
		Add:              entities.Wallet{ID: fromID},
		Substract:        entities.Wallet{ID: toID},
		Amount:           amount,
		ProcessingStatus: int(StatusProcessing),
	}

}

func (tm *TransactionManager) Transfer(ctx context.Context, tx *entities.Transaction) {

	PrepareShardRequests := getPreShardRequests(tx)
	for _, req := range PrepareShardRequests {
		go tm.provider.sendPreRequest(req)
	}

	tx.Substract.PrepareFlag = <-PrepareShardRequests["from"].ResponseChan
	tx.Add.PrepareFlag = <-PrepareShardRequests["to"].ResponseChan

}

func getPreShardRequests(tx *entities.Transaction) map[string]entities.PrepareTransactionRequest {
	return map[string]entities.PrepareTransactionRequest{
		"to":   tx.ToAddRequest(),
		"from": tx.ToSubstractRequest(),
	}
}

func generateID() string {
	rand.Seed(time.Now().UTC().UnixNano())

	bytes := make([]byte, 30)

	for i := 0; i < 30; i++ {
		bytes[i] = byte(randInt(97, 122))
	}

	return string(bytes)
}
func randInt(min int, max int) int {

	return min + rand.Intn(max-min)
}
