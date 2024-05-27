package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type (
	ILogstash interface {
		InsertLogstash(ctx context.Context, domainUuid string, logstash model.Logstash) error
	}
	ILogstashES interface {
		GetDocById(ctx context.Context, domainUuid, id string) (*model.Logstash, error)
	}
)

var LogstashRepo ILogstash
var LogstashRepoES ILogstashES
