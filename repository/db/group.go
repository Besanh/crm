package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type Group struct {
}

func NewGroup() repository.IGroup {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	repository.FusionSqlClient.GetDB().RegisterModel((*model.GroupUser)(nil))
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Group)(nil), "department_uuid char(36) NULL DEFAULT NULL"); err != nil {
		log.Fatal(err)
	}
	return &Group{}
}

func (repo *Group) SelectGroupsOfDomain(ctx context.Context, domainUuid string, filter model.GroupFilter, limit, offset int) ([]model.Group, int, error) {
	groups := make([]model.Group, 0)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&groups).
		Where("domain_uuid = ?", domainUuid)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if len(filter.GroupName) > 0 {
		query = query.Where("group_name LIKE ?", "%"+filter.GroupName+"%")
	}
	if len(filter.ManageGroupUuids) > 0 {
		query = query.Where("group_uuid IN (?)", bun.In(filter.ManageGroupUuids))
	}
	count, err := query.ScanAndCount(ctx, &groups)
	if err == sql.ErrNoRows {
		return groups, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return groups, count, nil
	}
}

func (repo *Group) SelectGroupsWithTotalUser(ctx context.Context, domainUuid string, filter model.GroupFilter, limit, offset int) (*[]model.GroupViewWithTotalUser, int, error) {
	groups := new([]model.GroupViewWithTotalUser)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(groups).
		ColumnExpr("g.*").
		ColumnExpr("d.domain_name").
		ColumnExpr("dp.department_name").
		ColumnExpr(`
		CASE
			WHEN g.domain_uuid IS NULL THEN false
			ELSE true
		END AS is_allow_delete
		`).
		ColumnExpr("count(u.user_uuid) as total_users").
		Join("INNER JOIN v_domains d on d.domain_uuid = g.domain_uuid").
		Join("LEFT JOIN v_group_users gu on gu.group_uuid = g.group_uuid").
		Join("LEFT JOIN v_users u on u.user_uuid = gu.user_uuid").
		Join("LEFT JOIN department dp ON dp.department_uuid = g.department_uuid").
		Where("g.domain_uuid = ?", domainUuid).
		Group("g.group_uuid", "d.domain_uuid", "dp.department_uuid")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if len(filter.GroupName) > 0 {
		query = query.Where("g.group_name LIKE ?", "%"+filter.GroupName+"%")
	}
	if len(filter.ManageGroupUuids) > 0 {
		query = query.Where("g.group_uuid IN (?)", bun.In(filter.ManageGroupUuids))
	}
	count, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return groups, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return groups, count, nil
	}
}

func (repo *Group) SelectGroupById(ctx context.Context, domainUuid string, groupUuid string) (*model.Group, error) {
	group := new(model.Group)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(group).
		Where("domain_uuid = ?", domainUuid).
		Where("group_uuid = ?", groupUuid)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return group, nil
	}
}

func (repo *Group) SelectGroupByName(ctx context.Context, domainUuid string, groupName string) (*model.Group, error) {
	group := new(model.Group)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(group).
		Where("domain_uuid = ?", domainUuid).
		Where("group_name = ?", groupName)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return group, nil
	}
}

func (repo *Group) InsertGroup(ctx context.Context, group *model.Group) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(group)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert group failed")
	}
	return nil
}

func (repo *Group) InsertGroupUser(ctx context.Context, groupUser *model.GroupUser) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(groupUser)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert group_user failed")
	}
	return nil
}

func (repo *Group) SelectGroupInfoById(ctx context.Context, domainUuid string, groupUuid string) (*model.GroupView, error) {
	group := new(model.GroupView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(group).
		ColumnExpr("g.*").
		ColumnExpr("d.department_name").
		Relation("Users", func(sq *bun.SelectQuery) *bun.SelectQuery {
			sq = sq.
				ColumnExpr("u.*").
				ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
				Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
				Where("u.domain_uuid = ?", domainUuid)
			return sq
		}).
		Join("LEFT JOIN department d ON d.department_uuid = g.department_uuid").
		Relation("Campaigns").
		Relation("GroupKPI").
		Where("g.domain_uuid = ?", domainUuid).
		Where("g.group_uuid = ?", groupUuid)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return group, nil
	}
}

func (repo *Group) UpdateGroup(ctx context.Context, group *model.Group) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(group).WherePK()
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 0 {
		return errors.New("update group failed")
	}
	return nil
}

// func (repo *Group) UpdateGroupFullTransaction(ctx context.Context, group *model.Group, campaignGroups *[]model.CampaignGroup, groupUsers *[]model.GroupUser, groupKpi *model.GroupKPI) error {
// 	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
// 		if _, err := tx.NewUpdate().Model(group).WherePK().Exec(ctx); err != nil {
// 			return err
// 		}
// 		if _, err := tx.NewDelete().Model((*model.CampaignGroup)(nil)).
// 			Where("group_uuid = ?", group.GroupUuid).
// 			Where("domain_uuid = ?", group.DomainUuid).
// 			Exec(ctx); err != nil {
// 			return err
// 		}
// 		if len(*campaignGroups) > 0 {
// 			if _, err := tx.NewInsert().Model(campaignGroups).Exec(ctx); err != nil {
// 				return err
// 			}
// 		}
// 		if _, err := tx.NewDelete().Model((*model.GroupUser)(nil)).
// 			Where("group_uuid = ?", group.GroupUuid).
// 			Where("domain_uuid = ?", group.DomainUuid).
// 			Exec(ctx); err != nil {
// 			return err
// 		}
// 		if len(*groupUsers) > 0 {
// 			if _, err := tx.NewInsert().Model(groupUsers).Exec(ctx); err != nil {
// 				return err
// 			}
// 		}

// 		if groupKpi != nil {
// 			if _, err := tx.NewUpdate().Model(groupKpi).WherePK().Exec(ctx); err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})
// 	return err
// }

// func (repo *Group) DeleteGroupTransaction(ctx context.Context, domainUuid string, groupUuid string) error {
// 	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
// 		if _, err := tx.NewDelete().Model((*model.Group)(nil)).
// 			Where("group_uuid = ?", groupUuid).
// 			Where("domain_uuid = ?", domainUuid).
// 			Exec(ctx); err != nil {
// 			return err
// 		}
// 		if _, err := tx.NewDelete().Model((*model.GroupUser)(nil)).
// 			Where("group_uuid = ?", groupUuid).
// 			Exec(ctx); err != nil {
// 			return err
// 		}
// 		if _, err := tx.NewDelete().Model((*model.CampaignGroup)(nil)).
// 			Where("group_uuid = ?", groupUuid).
// 			Exec(ctx); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
