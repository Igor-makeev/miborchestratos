package entities

type Transaction struct {
	TXID string

	Amount int

	Substract Wallet
	Add       Wallet

	ProcessingStatus int
}

type Wallet struct {
	ID          int
	PrepareFlag bool
	CommitFlag  bool
}

func (t *Transaction) ToAddRequest() PrepareTransactionRequest {
	return PrepareTransactionRequest{TxID: t.TXID, WalletID: t.Add.ID, Amount: t.Amount, ResponseChan: make(chan bool)}
}
func (t *Transaction) ToSubstractRequest() PrepareTransactionRequest {
	return PrepareTransactionRequest{TxID: t.TXID, WalletID: t.Substract.ID, Amount: -t.Amount, ResponseChan: make(chan bool)}
}
