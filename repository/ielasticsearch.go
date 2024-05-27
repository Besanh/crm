package repository

import "context"

type IElasticsearch interface {
	CheckAliasExist(ctx context.Context, alias string) (bool, error)
	CreateAlias(ctx context.Context, alias string) error
	CreateDocRabbitMQ(ctx context.Context, tenant, routing, uuid, vendor string, esDoc map[string]any) (bool, error)
	CreateAliasRabbitMQ(ctx context.Context, alias string) (bool, error)
	InsertLog(ctx context.Context, domainUuid, index, docId string, esDoc map[string]any) error
	UpdateDocById(ctx context.Context, domainUuid, index, docId string, esDoc map[string]any) error
}

var ESRepo IElasticsearch
