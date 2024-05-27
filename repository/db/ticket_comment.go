package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"
)

type TicketComment struct{}

func NewTicketCommentRepo() repository.ITicketComment {
	repo := &TicketComment{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *TicketComment) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	repository.FusionSqlClient.GetDB().RegisterModel((*model.TicketComment)(nil))
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.TicketComment)(nil)); err != nil {
		panic(err)
	}
}

func (repo *TicketComment) InitColumn() {
}

func (repo *TicketComment) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.TicketComment)(nil)).IfNotExists().Index("idx_ticket_comment_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *TicketComment) InsertTicketComment(ctx context.Context, domainUuid string, ticketComment model.TicketComment) error {
	resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&ticketComment).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("insert ticket_comment failed")
	}
	return nil
}

func (repo *TicketComment) GetCommentByTicketId(ctx context.Context, domainUuid, ticketUuid string) (*[]model.TicketComment, error) {
	comment := new([]model.TicketComment)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(comment).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return comment, nil
}

func (repo *TicketComment) DeleteTicketComment(ctx context.Context, domainUuid, ticketUuid string) error {
	TicketComment := new([]model.TicketComment)
	_, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model(TicketComment).
		Where("domain_uuid = ?", domainUuid).
		Where("ticket_uuid = ?", ticketUuid).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
