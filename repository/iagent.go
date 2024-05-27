package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type IAgent interface {
	InsertAgentTransaction(ctx context.Context, user *model.User, contact *model.VContact, groupUser *[]model.GroupUser, extension *model.Extension, extensionUser *model.ExtensionUser, callCenterAgent *model.CallCenterAgent, contactEmail *model.VContactEmail, roleGroup *model.RoleGroup) error
	UpdateAgentTransaction(ctx context.Context, user *model.User, contact *model.VContact, extension *model.Extension, extensionUser *model.ExtensionUser, callCenterAgent *model.CallCenterAgent, contactEmail *model.VContactEmail) error
}

var AgentRepo IAgent
