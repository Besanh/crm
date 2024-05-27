package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type UserCrmRepo struct{}

func NewUserCrm() repository.IUserCrm {
	repo := &UserCrmRepo{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *UserCrmRepo) InitTable() {
}

func (repo *UserCrmRepo) InitColumn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := repository.FusionSqlClient.GetDB().NewAddColumn().Model((*model.User)(nil)).IfNotExists().ColumnExpr("unit_uuid uuid").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *UserCrmRepo) InitIndex() {
}

func (repo *UserCrmRepo) GetUserCrms(ctx context.Context, domainUuid string, limit, offset int, filter model.UserFilter) (int, []model.UserView, error) {
	users := new([]model.UserView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(users).
		ColumnExpr("u.user_uuid, e.extension_uuid, u.domain_uuid, u.user_enabled, u.user_status").
		ColumnExpr("u.username, e.extension, u.level").
		ColumnExpr("u.role_uuid").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		ColumnExpr("units.unit_uuid, units.unit_name").
		Relation("RoleGroups").
		Relation("Units").
		Join("LEFT JOIN v_extension_users as eu ON eu.user_uuid = u.user_uuid").
		Join("LEFT JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid")
	if len(domainUuid) > 0 && domainUuid != "all" {
		query.Where("u.domain_uuid = ?", domainUuid)
	}

	if len(filter.ManageUserUuids) > 0 {
		query.Where("u.user_uuid IN (?)", bun.In(filter.ManageUserUuids))
	}

	if len(filter.UserUuid) > 0 {
		query.Where("u.user_uuid IN (?)", bun.In(filter.UserUuid))
	}

	if len(filter.StartTime) > 0 {
		query.Where("u.add_date >= ? ", filter.StartTime)
	}

	if len(filter.EndTime) > 0 {
		query.Where("u.add_date <= ? ", filter.EndTime)
	}

	if len(filter.Common) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("u.username ILIKE ?", "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("c.contact_name_given"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("c.contact_name_middle"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("c.contact_name_family"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("ce.email_address"), "%"+filter.Common+"%").
				WhereOr("? LIKE ?", bun.Ident("e.extension"), "%"+filter.Common+"%")
		})
	} else {
		if len(filter.Fullname) > 0 {
			query.Where("c.contact_name_given ILIKE ?", "%"+filter.Fullname+"%").
				WhereOr("c.contact_name_middle ILIKE ?", "%"+filter.Fullname+"%").
				WhereOr("c.contact_name_family ILIKE ?", "%"+filter.Fullname+"%")
		}
		if len(filter.Email) > 0 {
			query.Where("ce.email_address ILIKE ?", "%"+filter.Email+"%")
		}
		if len(filter.Extension) > 0 {
			query.Where("e.extension LIKE ?", "%"+filter.Extension+"%")
		}
		if len(filter.Name) > 0 {
			query.Where("u.username ILIKE ?", "%"+filter.Name+"%")
		}
	}

	if len(filter.Level) > 0 {
		query.
			Where("u.level = ?", filter.Level)
	}
	if len(filter.Levels) > 0 {
		query.
			Where("u.level IN (?)", bun.In(filter.Levels))
	}
	if len(filter.Enabled) > 0 {
		query.Where("u.user_enabled = ?", filter.Enabled)
	}

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	// if !filter.IsAll {
	// 	if len(filter.ManageUserUuids) > 0 {
	// 		query.
	// 			Where("u.user_uuid IN (?)", bun.In(filter.ManageUserUuids))
	// 	}
	// }
	if filter.IsMapExtension == "true" {
		query.Where("e.extension_uuid IS NOT NULL")
	} else if filter.IsMapExtension == "false" {
		query.Where("e.extension_uuid IS NULL")
	}
	if len(filter.RoleUuid) > 0 {
		query.Where("u.role_uuid = ?", filter.RoleUuid)
	}
	if len(filter.UnitUuid) > 0 {
		query.Where("u.unit_uuid = ?", filter.UnitUuid)
	}
	if len(filter.ManageExcludeUnitUuids) > 0 {
		query.Where("u.unit_uuid NOT IN (?)", bun.In(filter.ManageExcludeUnitUuids))
	}
	order := ""
	switch filter.Order {
	case "username":
		order = "u.username"
	case "email":
		order = "ce.email_address"
	case "extension":
		order = "e.extension"
	case "role":
		order = "role.role_name"
	case "user_enabled":
		order = "u.user_enabled"
	default:
		order = "u.username"
	}
	switch filter.OrderDirection {
	case "desc":
		order = order + " DESC"
	default:
		order = order + " ASC"
	}
	if len(order) > 0 {
		query.Order(order)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, errors.New("no data")
	} else if err != nil {
		return 0, nil, err
	}
	return total, *users, nil
}

func (repo *UserCrmRepo) GetUserCrmById(ctx context.Context, id string) (*model.UserView, error) {
	result := new(model.UserView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(result).
		Relation("RoleGroups").
		Relation("Units").
		ColumnExpr("u.*, e.extension_uuid, e.extension").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		ColumnExpr("units.unit_uuid, units.unit_name").
		Join("LEFT JOIN v_extension_users as eu ON eu.user_uuid = u.user_uuid").
		Join("LEFT JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid").
		Where("u.user_uuid = ?", id)
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, nil
}

func (repo *UserCrmRepo) GetUsersInfoOfUnit(ctx context.Context, domainUuid string, userUuid, userlevel string, unitUuids []string, level []string, isExtension bool) (*[]model.UserInfoData, error) {
	users := new([]model.UserInfoData)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(users).
		ColumnExpr("DISTINCT u.user_uuid, u.username").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		Join("LEFT JOIN v_contacts c ON c.contact_uuid = u.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid ").
		Join("LEFT JOIN unit ON unit.unit_uuid = u.unit_uuid").
		Where("u.domain_uuid = ?", domainUuid)
	if len(unitUuids) > 0 {
		query.Where("unit.unit_uuid IN (?)", bun.In(unitUuids))
	}
	if len(level) > 0 {
		query.Where("u.level IN (?)", bun.In(level))
	}
	if len(userlevel) > 0 {
		switch userlevel {
		case "leader":
			query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("u.level IN (?)", bun.In([]string{"agent", "user"})).
					WhereOr("u.level = ? AND u.user_uuid = ?", "leader", userUuid)
			})
		case "manager":
			query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("u.level IN (?)", bun.In([]string{"leader", "agent", "user"})).
					WhereOr("u.level = ? AND u.user_uuid = ?", "manager", userUuid)
			})
		}
	}
	if isExtension {
		subQuery := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("v_extensions as e").
			Join("LEFT JOIN v_extension_users as eu ON eu.extension_uuid = e.extension_uuid").
			Where("eu.user_uuid = u.user_uuid").
			Column("e.extension_uuid", "e.extension").
			Limit(1)
		query.ColumnExpr("tmpe.*").Join("LEFT JOIN lateral (?) tmpe ON true", subQuery)
	}

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return users, nil
	} else if err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

func (repo *UserCrmRepo) PatchUserCrm(ctx context.Context, domainUuid, id, unitUuid, roleGroupUuid string) error {
	var model model.User
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&model).
		Where("user_uuid = ?", id).
		Set("unit_uuid = ?", unitUuid).
		Set("role_uuid = ?", roleGroupUuid)
	_, err := query.Exec(ctx)
	return err
}
