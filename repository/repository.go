package repository

import (
	"contactcenter-api/internal/elasticsearch"
	sqlclient "contactcenter-api/internal/sqlclient"
	"context"
	"strings"

	"github.com/uptrace/bun/schema"
)

var FusionSqlClient sqlclient.ISqlClientConn
var FreeswitchSqlClient sqlclient.ISqlClientConn
var ES elasticsearch.IESClient
var ESClient elasticsearch.IElasticsearchClient

func CreateTableCollate(client sqlclient.ISqlClientConn, ctx context.Context, table any) error {
	query := client.GetDB().NewCreateTable().Model(table).IfNotExists()
	value, _ := query.AppendQuery(schema.NewFormatter(query.Dialect()), nil)
	queryStr := string(value) + " COLLATE utf8mb4_general_ci"
	_, err := client.GetDB().QueryContext(ctx, queryStr)
	return err
}

func CreateTable(client sqlclient.ISqlClientConn, ctx context.Context, table any) error {
	query := client.GetDB().NewCreateTable().Model(table).IfNotExists()
	value, _ := query.AppendQuery(schema.NewFormatter(query.Dialect()), nil)
	queryStr := string(value)
	if client.GetDriver() == sqlclient.POSTGRESQL {
		queryStr = strings.ReplaceAll(queryStr, " char(36)", " uuid")
		queryStr = strings.ReplaceAll(queryStr, " timestamp", " timestamptz")
		queryStr = strings.ReplaceAll(queryStr, " timestamptz_only", " timestamp")
	}
	_, err := client.GetDB().QueryContext(ctx, queryStr)
	return err
}

func AddColumn(client sqlclient.ISqlClientConn, ctx context.Context, table any, column string) error {
	query := client.GetDB().NewAddColumn().Model(table).IfNotExists().ColumnExpr(column)
	value, _ := query.AppendQuery(schema.NewFormatter(query.Dialect()), nil)
	queryStr := string(value)
	if client.GetDriver() == sqlclient.POSTGRESQL {
		queryStr = strings.ReplaceAll(queryStr, " char(36)", " uuid")
		queryStr = strings.ReplaceAll(queryStr, " timestamp", " timestamptz")
	}
	_, err := client.GetDB().QueryContext(ctx, queryStr)
	return err
}

func Query(client sqlclient.ISqlClientConn, ctx context.Context, query string) error {
	if client.GetDriver() == sqlclient.POSTGRESQL {
		query = strings.ReplaceAll(query, " char(36)", " uuid")
		query = strings.ReplaceAll(query, " timestamp", " timestamptz")
	}
	_, err := client.GetDB().QueryContext(ctx, query)
	return err
}
