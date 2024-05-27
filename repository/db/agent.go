package db

import (
	"contactcenter-api/common/model"
	sqlclient "contactcenter-api/internal/sqlclient"
	"contactcenter-api/repository"
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type Agent struct {
}

func NewAgent() repository.IAgent {
	return &Agent{}
}

func (repo *Agent) InsertAgentTransaction(ctx context.Context, user *model.User, contact *model.VContact, groupUser *[]model.GroupUser, extension *model.Extension, extensionUser *model.ExtensionUser, callCenterAgent *model.CallCenterAgent, contactEmail *model.VContactEmail, roleGroup *model.RoleGroup) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(contact).Exec(ctx); err != nil {
			return err
		}
		if groupUser != nil && len(*groupUser) > 0 {
			if _, err := tx.NewInsert().Model(groupUser).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewInsert().Model(extension).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(extensionUser).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(callCenterAgent).Exec(ctx); err != nil {
			return err
		}

		// Contact crm
		if contactEmail != nil {
			if _, err := tx.NewInsert().Model(contactEmail).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Agent) UpdateAgentTransaction(ctx context.Context, user *model.User, contact *model.VContact, extension *model.Extension, extensionUser *model.ExtensionUser, callCenterAgent *model.CallCenterAgent, contactEmail *model.VContactEmail) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(user).WherePK().
			Column("username", "password", "salt", "level", "role_uuid", "enable_webrtc").
			Exec(ctx); err != nil {
			return err
		}
		if repository.FusionSqlClient.GetDriver() == sqlclient.MYSQL {
			if _, err := tx.NewInsert().Model(contact).
				Where("c.contact_uuid = ?", contact.ContactUuid).
				On("DUPLICATE KEY UPDATE").
				Set("contact_name_given = ?", contact.ContactNameGiven).
				Set("contact_name_middle = ?", contact.ContactNameMiddle).
				Set("contact_name_family = ?", contact.ContactNameFamily).
				Set("contact_nickname = ?", contact.ContactNickname).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			if _, err := tx.NewInsert().Model(contact).
				Where("c.contact_uuid = ?", contact.ContactUuid).
				On("CONFLICT (contact_uuid) DO UPDATE").
				Set("contact_name_given = ?", contact.ContactNameGiven).
				Set("contact_name_middle = ?", contact.ContactNameMiddle).
				Set("contact_name_family = ?", contact.ContactNameFamily).
				Set("contact_nickname = ?", contact.ContactNickname).
				Exec(ctx); err != nil {
				return err
			}
		}
		if extension != nil {
			if _, err := tx.NewInsert().Model(extension).Exec(ctx); err != nil {
				return err
			}
		}
		if extensionUser != nil {
			if _, err := tx.NewDelete().Model((*model.ExtensionUser)(nil)).
				Where("user_uuid = ?", extensionUser.UserUuid).
				WhereOr("extension_uuid = ?", extensionUser.ExtensionUuid).
				Exec(ctx); err != nil {
				return err
			}
			if _, err := tx.NewInsert().Model(extensionUser).Exec(ctx); err != nil {
				return err
			}
		}
		// if _, err := tx.NewInsert().Model(callCenterAgent).Exec(ctx); err != nil {
		// 	return err
		// }
		if repository.FusionSqlClient.GetDriver() == sqlclient.MYSQL {
			if _, err := tx.NewInsert().Model(callCenterAgent).
				Where("cca.user_uuid = ?", callCenterAgent.UserUuid).
				On("DUPLICATE KEY UPDATE").
				Set("agent_name = ?", callCenterAgent.AgentName).
				Set("agent_id = ?", callCenterAgent.AgentId).
				Set("agent_contact = ?", callCenterAgent.AgentContact).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			query := tx.NewInsert().Model(callCenterAgent).
				Where("cca.user_uuid = ?", callCenterAgent.UserUuid).
				On("CONFLICT (user_uuid) DO UPDATE").
				Set("agent_name = ?", callCenterAgent.AgentName).
				Set("agent_id = ?", callCenterAgent.AgentId).
				Set("agent_contact = ?", callCenterAgent.AgentContact)
			if _, err := query.
				Exec(ctx); err != nil {
				return err
			}
		}

		if repository.FusionSqlClient.GetDriver() == sqlclient.MYSQL {
			if _, err := tx.NewInsert().Model(contactEmail).
				Where("ce.contact_uuid = ?", contactEmail.ContactUuid).
				On("DUPLICATE KEY UPDATE").
				Set("email_address = ?", contactEmail.EmailAddress).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			if _, err := tx.NewInsert().Model(contactEmail).
				Where("ce.contact_uuid = ?", contactEmail.ContactUuid).
				On("CONFLICT (contact_uuid) DO UPDATE").
				Set("email_address = ?", contactEmail.EmailAddress).
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}
