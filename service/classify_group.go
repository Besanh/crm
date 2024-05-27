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
	IClassifyGroup interface {
		PostClassifyGroup(ctx context.Context, domainUuid, userUuid string, classifyGroup *model.ClassifyGroup) (int, any)
		GetClassifyGroups(ctx context.Context, domainUuid string, limit, offset int, filter model.ClassifyGroupFilter) (int, any)
		GetClassifyGroupById(ctx context.Context, domainUuid, id string) (int, any)
		PutClassifyGroupById(ctx context.Context, domainUuid, userUuid, id string, ClassifyGroup model.ClassifyGroup) (int, any)
		// DeleteClassifyGroup(ctx context.Context, domainUuid, id string) (int, any)
		// ExportClassifyGroups(ctx context.Context, domainUuid, userUuid, fileType string, filter model.ClassifyGroupFilter) (string, error)
	}
	ClassifyGroup struct {
	}
)

func NewClassifyGroup() IClassifyGroup {
	return &ClassifyGroup{}
}

func (s *ClassifyGroup) PostClassifyGroup(ctx context.Context, domainUuid, userUuid string, classifyGroup *model.ClassifyGroup) (int, any) {
	classifyGroup.DomainUuid = domainUuid
	classifyGroup.ClassifyGroupUuid = uuid.NewString()
	classifyGroup.CreatedBy = userUuid
	classifyGroup.CreatedAt = time.Now()

	if err := repository.ClassifyGroupRepo.InsertClassifyGroup(ctx, *classifyGroup); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(map[string]any{
		"id": classifyGroup.ClassifyGroupUuid,
	})
}

func (s *ClassifyGroup) GetClassifyGroups(ctx context.Context, domainUuid string, limit, offset int, filter model.ClassifyGroupFilter) (int, any) {
	total, ClassifyGroups, err := repository.ClassifyGroupRepo.GetClassifyGroups(ctx, domainUuid, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Pagination(ClassifyGroups, total, limit, offset)
}

func (s *ClassifyGroup) GetClassifyGroupById(ctx context.Context, domainUuid, id string) (int, any) {
	classifyGroup, err := repository.ClassifyGroupRepo.GetClassifyGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(classifyGroup)
}

func (s *ClassifyGroup) PutClassifyGroupById(ctx context.Context, domainUuid, userUuid, id string, classifyGroup model.ClassifyGroup) (int, any) {
	classifyGroupExist, err := repository.ClassifyGroupRepo.GetClassifyGroupById(ctx, domainUuid, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	// classifyGroupInfoExist, err := repository.ClassifyGroupRepo.GetClassifyGroupById(ctx, domainUuid, classifyGroupExist.ClassifyGroupUuid)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// } else if classifyGroupInfoExist == nil {
	// 	return response.BadRequestMsg("Classify Group name not found")
	// } else if classifyGroupInfoExist.GroupName != classifyGroup.GroupName {
	// 	return response.BadRequestMsg("Classify Group name already exist")
	// }

	classifyGroupExist.GroupName = classifyGroup.GroupName
	// classifyGroupExist.LimitedFunction = classifyGroup.LimitedFunction
	classifyGroupExist.Status = classifyGroup.Status
	classifyGroupExist.Description = classifyGroup.Description
	classifyGroupExist.UpdatedBy = userUuid
	classifyGroupExist.UpdatedAt = time.Now()

	// classifyGroupUser := []model.classifyGroupUser{}
	// if len(ClassifyGroup.Member) > 0 {
	// 	for _, val := range ClassifyGroup.Member {
	// 		classifyGroupUser = append(classifyGroupUser, model.classifyGroupUser{
	// 			DomainUuid:         domainUuid,
	// 			classifyGroupUserUuid: uuid.NewString(),
	// 			ClassifyGroupUuid:     ClassifyGroup.ClassifyGroupUuid,
	// 			EntityType:         "member",
	// 			UserUuid:           val,
	// 		})
	// 	}
	// }

	// if len(ClassifyGroup.Staff) > 0 {
	// 	for _, val := range ClassifyGroup.Staff {
	// 		classifyGroupUser = append(classifyGroupUser, model.classifyGroupUser{
	// 			DomainUuid:         domainUuid,
	// 			classifyGroupUserUuid: uuid.NewString(),
	// 			EntityType:         "staff",
	// 			ClassifyGroupUuid:     ClassifyGroup.ClassifyGroupUuid,
	// 			UserUuid:           val,
	// 		})
	// 	}
	// }

	if err := repository.ClassifyGroupRepo.InsertClassifyGroup(ctx, *classifyGroupExist); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"id": id,
	})
}

// func (s *ClassifyGroup) DeleteClassifyGroup(ctx context.Context, domainUuid, id string) (int, any) {
// 	ClassifyGroupExist, err := repository.ClassifyGroupRepo.GetClassifyGroupById(ctx, domainUuid, id)
// 	if err != nil {
// 		log.Error(err)
// 		return response.ServiceUnavailableMsg(err.Error())
// 	} else if ClassifyGroupExist == nil {
// 		return response.NotFoundMsg("Classify Group not found")
// 	}

// 	if err := repository.ClassifyGroupRepo.DeleteClassifyGroup(ctx, domainUuid, ClassifyGroupExist); err != nil {
// 		log.Error(err)
// 		return response.ServiceUnavailableMsg(err.Error())
// 	}

// 	return response.OK(map[string]any{
// 		"id": id,
// 	})
// }
