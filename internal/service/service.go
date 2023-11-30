package service

import (
	"context"
	config "miborchestrator/configs"
	"miborchestrator/internal/entities"
	"miborchestrator/internal/repository"
)

type WalletCreator interface {
	createWalletRequest(entities.WcTDO)
}
type OrchestratorService interface {
	Createwallet(walletID int)
	InitTransfer(fromID, toID, amount int)
}

type Authorization interface {
	CreateUser(ctx context.Context, user entities.User) error
	GenerateToken(ctx context.Context, login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
	TxManager *TransactionManager
	WcQueue   *WalletCreationQueue
	OrSrv     OrchestratorService
	Config    *config.Config
}

func NewService(ctx context.Context, Config *config.Config, repo *repository.Repository) *Service {
	provider := NewProvider(Config)
	srv := &Service{
		Authorization: NewAuthService(repo, Config),
		TxManager:     NewTransactionManager(ctx, provider),
		WcQueue:       NewWalletCreationQueue(provider),

		Config: Config,
	}

	return srv
}
func (s *Service) AddtoCreateWalletQueue(req *entities.WcTDO) {
	s.WcQueue.buf <- req
}

func (s *Service) AddtoTxManagerQueue(req entities.Request) {
	s.TxManager.stream <- req
}

func (s *Service) Close(ctx context.Context) error {
	return nil
}
