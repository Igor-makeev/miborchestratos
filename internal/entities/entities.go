package entities

type eventType int

const (
	InitTx eventType = iota
	PreparationsConfirmed
	CommitConfitmed
	TxSuccessfull
)

type TransactionRequest struct {
	WalletIDFrom int `json:"wallet_id_from" binding:"required"`
	WalletIDTo   int `json:"wallet_id_to" binding:"required"`
	Amount       int `json:"amount" binding:"required"`
}
type CreateWalletRequest struct {
	WalletID int `json:"wallet_id" binding:"required"`
}
type WcTDO struct {
	WalletID int `json:"wallet_id" binding:"required"`
	UserID   int `json:"user_id" binding:"required"`
}

type PrepareTransactionRequest struct {
	TxID         string    `json:"transaction_id" binding:"required"`
	WalletID     int       `json:"wallet_id" binding:"required"`
	Amount       int       `json:"amount" binding:"required"`
	ResponseChan chan bool `json:"-"`
}

type TxDTO struct {
	WalletID int       `json:"wallet_id" binding:"required"`
	TxID     string    `json:"transaction_id" binding:"required"`
	DTOEvent eventType `json:"event" binding:"required"`
}

type Request struct {
	TransactionRequest
	TxDTO
	Event eventType
}
