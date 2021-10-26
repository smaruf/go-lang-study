package proto

import (
	. "common/proto"
	"context"
)

type WalletService struct {
	//todo: call appropriate internal and external service or push
}

func NewWalletService() *WalletService {
	return &WalletService{}
}

func (w WalletService) GetAccount(ctx context.Context, filter *AccountFilter) (*AccountInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletAccountClient()
	defer service.Conn.Close()
	return service.Client.GetAccount(ctx, filter)
}
func (w WalletService) CreateAccount(ctx context.Context, filter *AccountInfo) (*AccountInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletAccountClient()
	defer service.Conn.Close()
	return service.Client.CreateAccount(ctx, filter)
}
func (w WalletService) CloseAccount(ctx context.Context, filter *AccountInfo) (*AccountInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletAccountClient()
	defer service.Conn.Close()
	return service.Client.CloseAccount(ctx, filter)
}

func (w WalletService) CheckAccount(ctx context.Context, filter *AccountFilter) (*Message, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletAccountClient()
	defer service.Conn.Close()
	return service.Client.CheckAccount(ctx, filter)
}

func (w WalletService) GetAccountBalance(ctx context.Context, filter *AccountFilter) (*BalanceInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletAccountClient()
	defer service.Conn.Close()
	return service.Client.GetAccountBalance(ctx, filter)
}

func (w WalletService) GetTransaction(ctx context.Context, filter *TransactionFilter) (*TransactionInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletSearchClient()
	defer service.Conn.Close()
	return service.Client.GetTransaction(ctx, filter)
}

func (w WalletService) FindTransactions(filter *TransactionFilter, server WalletService_FindTransactionsServer) error {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletSearchClient()
	defer service.Conn.Close()
	serves, err := service.Client.FindTransactions(context.Background(), filter)
	if err != nil {
		return err
	}
	return serves.SendMsg(server)
}

func (w WalletService) InitiateTransfer(ctx context.Context, info *TransactionInfo) (*TransactionInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletTransferClient()
	defer service.Conn.Close()
	return service.Client.InitiateTransfer(ctx, info)
}

func (w WalletService) ConfirmTransfer(ctx context.Context, info *TransactionInfo) (*TransactionInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletTransferClient()
	defer service.Conn.Close()
	return service.Client.ConfirmTransfer(ctx, info)
}

func (w WalletService) RevertTransfer(ctx context.Context, info *TransactionInfo) (*TransactionInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletTransferClient()
	defer service.Conn.Close()
	return service.Client.RevertTransfer(ctx, info)
}

func (w WalletService) RequestTransfer(ctx context.Context, info *TransactionInfo) (*Message, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletTransferClient()
	defer service.Conn.Close()
	return service.Client.RequestTransfer(ctx, info)
}

func (w WalletService) ResponseTransferRequest(ctx context.Context, info *TransactionInfo) (*TransactionInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletTransferClient()
	defer service.Conn.Close()
	return service.Client.ResponseTransferRequest(ctx, info)
}

func (w WalletService) ManageAccount(ctx context.Context, info *AccountInfo) (*AccountInfo, error) {
	//todo: need to calling steps
	// log audit trail
	// verify request
	service := NewWalletAccountClient()
	defer service.Conn.Close()
	return service.Client.ManageAccount(ctx, info)
}

func (w WalletService) mustEmbedUnimplementedWalletServiceServer() {
	panic("implement me")
}
