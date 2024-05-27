package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type (
	IClassifyCareerRepo interface {
		InsertClassifyCareer(ctx context.Context, classifyTag model.ClassifyCareer) error
		GetClassifyCareers(ctx context.Context, domainUuid string, limit, offset int) (int, []model.ClassifyCareer, error)
		GetClassifyCareerById(ctx context.Context, domainUuid, id string) (*model.ClassifyCareer, error)
		PutClassifyCareerById(ctx context.Context, domainUuid string, classifyCareer model.ClassifyCareer) error
	}
	IClassifyGroupRepo interface {
		InsertClassifyGroup(ctx context.Context, classifyTag model.ClassifyGroup) error
		GetClassifyGroups(ctx context.Context, domainUuid string, limit, offset int) (int, []model.ClassifyGroup, error)
		GetClassifyGroupById(ctx context.Context, domainUuid, id string) (*model.ClassifyGroup, error)
		PutClassifyGroupById(ctx context.Context, domainUuid string, classifyCareer model.ClassifyGroup) error
	}
	IClassifyTagRepo interface {
		InsertClassifyTag(ctx context.Context, classifyTag model.ClassifyTag) error
		GetClassifyTags(ctx context.Context, domainUuid string, limit, offset int) (int, []model.ClassifyTag, error)
		GetClassifyTagById(ctx context.Context, domainUuid, id string) (*model.ClassifyTag, error)
		PutClassifyTagById(ctx context.Context, domainUuid string, classifyCareer model.ClassifyTag) error
	}
)

var ClassifyCareerRepo IClassifyCareerRepo
var ClassifyGroupRepo IClassifyGroupRepo
var ClassifyTagRepo IClassifyTagRepo
