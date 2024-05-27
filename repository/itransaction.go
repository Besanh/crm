package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITransaction interface {
	InsertTransaction(ctx context.Context, transaction *model.Transaction) error
	InsertTransactionLog(ctx context.Context, log ...model.TransactionLog) error
}

var TransactionRepo ITransaction
