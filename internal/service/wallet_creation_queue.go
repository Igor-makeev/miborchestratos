package service

import (
	"context"
	"miborchestrator/internal/entities"
)

type WalletCreationQueue struct {
	WalletCreator
	buf chan *entities.WcTDO
}

func NewWalletCreationQueue(wk WalletCreator) *WalletCreationQueue {
	return &WalletCreationQueue{
		WalletCreator: wk,
		buf:           make(chan *entities.WcTDO, 100)}
}

func (wcq *WalletCreationQueue) Run(ctx context.Context) {
	go wcq.listen(ctx)
}

func (wcq *WalletCreationQueue) listen(ctx context.Context) {

	for {

		select {
		case req := <-wcq.buf:
			wcq.WalletCreator.createWalletRequest(*req)
		case <-ctx.Done():

			return
		}
	}
}
