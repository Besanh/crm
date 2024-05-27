package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	IClassifyTag interface {
		PostClassifyTag(ctx context.Context, domainUuid, userUuid string, classifyTag *model.ClassifyTag) (int, any)
		GetClassifyTags(ctx context.Context, domainUuid string, limit, offset int, filter model.ClassifyTagFilter) (int, any)
		GetClassifyTagById(ctx context.Context, domainUuid, id string) (int, any)
		PutClassifyTagById(ctx context.Context, domainUuid, userUuid, id string, ClassifyTag model.ClassifyTag) (int, any)
		// DeleteClassifyTag(ctx context.Context, domainUuid, id string) (int, any)
		// ExportClassifyTags(ctx context.Context, domainUuid, userUuid, fileType string, filter model.ClassifyTagFilter) (string, error)
	}
	ClassifyTag struct {
	}
)

func NewClassifyTag() IClassifyTag {
	return &ClassifyTag{}
}

func (s *ClassifyTag) PostClassifyTag(ctx context.Context, domainUuid, userUuid string, classifyTag *model.ClassifyTag) (int, any) {
	classifyTag.DomainUuid = domainUuid
	classifyTag.ClassifyTagUuid = uuid.NewString()
	classifyTag.CreatedBy = userUuid
	classifyTag.CreatedAt = time.Now()

	if err := repository.ClassifyTagRepo.InsertClassifyTag(ctx, *classifyTag); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": classifyTag.ClassifyTagUuid,
	})
}

func (s *ClassifyTag) GetClassifyTags(ctx context.Context, domainUuid string, limit, offset int, filter model.ClassifyTagFilter) (int, any) {
	total, ClassifyTags, err := repository.ClassifyTagRepo.GetClassifyTags(ctx, domainUuid, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(ClassifyTags, total, limit, offset)
}

func (s *ClassifyTag) GetClassifyTagById(ctx context.Context, domainUuid, id string) (int, any) {
	classifyTag, err := repository.ClassifyTagRepo.GetClassifyTagById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(classifyTag)
}

func (s *ClassifyTag) PutClassifyTagById(ctx context.Context, domainUuid, userUuid, id string, classifyTag model.ClassifyTag) (int, any) {
	classifyTagExist, err := repository.ClassifyTagRepo.GetClassifyTagById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	// classifyTagInfoExist, err := repository.ClassifyTagRepo.GetClassifyTagById(ctx, domainUuid, classifyTagExist.ClassifyTagUuid)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// } else if classifyTagInfoExist == nil {
	// 	return response.BadRequestMsg("Classify tag name not found")
	// } else if classifyTagInfoExist.TagName != classifyTag.TagName {
	// 	return response.BadRequestMsg("Classify tag name already exist")
	// }

	classifyTagExist.TagName = classifyTag.TagName
	classifyTagExist.LimitedFunction = classifyTag.LimitedFunction
	classifyTagExist.Status = classifyTag.Status
	classifyTagExist.Description = classifyTag.Description
	classifyTagExist.UpdatedBy = userUuid
	classifyTagExist.UpdatedAt = time.Now()

	// classifyTagUser := []model.classifyTagUser{}
	// if len(ClassifyTag.Member) > 0 {
	// 	for _, val := range ClassifyTag.Member {
	// 		classifyTagUser = append(classifyTagUser, model.classifyTagUser{
	// 			DomainUuid:         domainUuid,
	// 			classifyTagUserUuid: uuid.NewString(),
	// 			ClassifyTagUuid:     ClassifyTag.ClassifyTagUuid,
	// 			EntityType:         "member",
	// 			UserUuid:           val,
	// 		})
	// 	}
	// }

	// if len(ClassifyTag.Staff) > 0 {
	// 	for _, val := range ClassifyTag.Staff {
	// 		classifyTagUser = append(classifyTagUser, model.classifyTagUser{
	// 			DomainUuid:         domainUuid,
	// 			classifyTagUserUuid: uuid.NewString(),
	// 			EntityType:         "staff",
	// 			ClassifyTagUuid:     ClassifyTag.ClassifyTagUuid,
	// 			UserUuid:           val,
	// 		})
	// 	}
	// }

	if err := repository.ClassifyTagRepo.InsertClassifyTag(ctx, *classifyTagExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}

// func (s *ClassifyTag) DeleteClassifyTag(ctx context.Context, domainUuid, id string) (int, any) {
// 	ClassifyTagExist, err := repository.ClassifyTagRepo.GetClassifyTagById(ctx, domainUuid, id)
// 	if err != nil {
// 		log.Error(err)
// 		return response.ServiceUnavailableMsg(err.Error())
// 	} else if ClassifyTagExist == nil {
// 		return response.NotFoundMsg("Classify tag not found")
// 	}

// 	if err := repository.ClassifyTagRepo.DeleteClassifyTag(ctx, domainUuid, ClassifyTagExist); err != nil {
// 		log.Error(err)
// 		return response.ServiceUnavailableMsg(err.Error())
// 	}

// 	return response.OK(map[string]any{
// 		"id": id,
// 	})
// }
