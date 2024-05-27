package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"errors"
	"time"
)

type (
	Transaction struct {
	}
)

func NewTransaction() repository.ITransaction {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.TransactionLog)(nil)); err != nil {
		panic(err)
	}
	return &Transaction{}
}

func (repo *Transaction) InsertTransaction(ctx context.Context, transaction *model.Transaction) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(transaction)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert transaction failed")
	}
	return nil
}

func (repo *Transaction) InsertTransactionLog(ctx context.Context, log ...model.TransactionLog) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(&log)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert transaction_log failed")
	}
	return nil
}
