package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type Extension struct {
}

func NewExtension() repository.IExtension {
	return &Extension{}
}

func (repo *Extension) GetExtensionByExten(ctx context.Context, domainUuid, exten string) (*model.ExtensionView, error) {
	extension := new(model.ExtensionView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(extension).
		ColumnExpr("e.extension_uuid, e.domain_uuid, e.extension, e.enabled").
		ColumnExpr("d.domain_name").
		Join("INNER JOIN v_domains d ON d.domain_uuid = e.domain_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		Where("e.extension = ?", exten).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return extension, nil
	}
}
func (repo *Extension) GetExtensionByUserUuid(ctx context.Context, domainUuid, userUuid string) (*model.ExtensionView, error) {
	extension := new(model.ExtensionView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(extension).
		ColumnExpr("d.domain_name").
		ColumnExpr("e.extension, e.extension_uuid, e.domain_uuid").
		Join("INNER JOIN v_domains d ON d.domain_uuid = e.domain_uuid").
		Join("JOIN v_extension_users as eu ON eu.extension_uuid = e.extension_uuid").
		Join("JOIN v_users as u ON u.user_uuid = eu.user_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		Where("u.user_uuid = ?", userUuid).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return extension, nil
	}
}

func (repo *Extension) GetExtensionInfoByUserUuid(ctx context.Context, domainUuid, userUuid string) (*model.ExtensionInfoWithPassword, error) {
	extension := new(model.ExtensionInfoWithPassword)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(extension).
		ColumnExpr("d.domain_name").
		ColumnExpr("e.extension, e.extension_uuid, e.password").
		Join("INNER JOIN v_domains d ON d.domain_uuid = e.domain_uuid").
		Join("JOIN v_extension_users as eu ON eu.extension_uuid = e.extension_uuid").
		Join("JOIN v_users as u ON u.user_uuid = eu.user_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		Where("u.user_uuid = ?", userUuid).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return extension, nil
	}
}

func (repo *Extension) UpdateExtensionEnabled(ctx context.Context, domainUuid, extension string, isEnabled bool) error {
	value := "true"
	if !isEnabled {
		value = "false"
	}
	query := repository.FusionSqlClient.GetDB().NewUpdate().
		Model((*model.Extension)(nil)).
		Set("enabled = ?", value).
		Where("extension = ?", extension).
		Where("domain_uuid = ?", domainUuid)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("update extension failed")
	}
	return nil
}

func (repo *Extension) GetExtensionsOfUserUuids(ctx context.Context, domainUuid string, userUuids []string) (*[]model.ExtensionData, error) {
	extensions := new([]model.ExtensionData)
	err := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(extensions).
		ColumnExpr("DISTINCT e.extension_uuid").
		Column("e.extension", "e.domain_uuid", "eu.user_uuid").
		Join("INNER JOIN v_extension_users eu ON eu.extension_uuid = e.extension_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		Where("eu.user_uuid IN (?)", bun.In(userUuids)).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return extensions, nil
	}
}

func (repo *Extension) GetExtensionsInfo(ctx context.Context, domainUuid string, filter model.ExtensionFilter, limit, offset int) (*[]model.ExtensionInfo, int, error) {
	extensions := new([]model.ExtensionInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(extensions).
		Join("LEFT JOIN v_extension_users eu ON e.extension_uuid = eu.extension_uuid").
		Join("LEFT JOIN v_users u ON u.user_uuid = eu.user_uuid").
		Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Join("INNER JOIN v_domains d ON d.domain_uuid = e.domain_uuid").
		Join("LEFT JOIN v_call_center_agents a ON a.user_uuid = u.user_uuid").
		Join("LEFT JOIN unit ON unit.unit_uuid = u.unit_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		Order("e.extension ASC").
		Limit(limit).Offset(offset)
	query = query.Column("u.user_uuid", "u.username", "e.extension_uuid", "e.extension", "e.enabled", "d.domain_name", "e.domain_uuid", "e.follow_me_uuid", "unit.unit_uuid", "unit.unit_name").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr(`
			CASE
			WHEN a.call_center_agent_uuid IS NOT NULL THEN true
			ELSE false
			END AS is_link_call_center
			`)
	if len(filter.ManageExtensionUuids) > 0 {
		query = query.Where("eu.extension_uuid IN (?)", bun.In(filter.ManageExtensionUuids))
	}
	if len(filter.Common) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("? ILIKE ?", bun.Ident("e.extension"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("u.username"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("c.contact_name_given"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("c.contact_name_middle"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("c.contact_name_family"), "%"+filter.Common+"%")
		})
	} else {
		if len(filter.Extension) > 0 {
			query = query.Where("e.extension ILIKE ?", "%"+filter.Extension+"%")
		}
		if len(filter.Username) > 0 {
			query = query.Where("u.username ILIKE ?", "%"+filter.Username+"%")
		}
		if len(filter.Fullname) > 0 {
			query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.WhereOr("c.contact_name_given ILIKE ?", "%"+filter.Fullname+"%").
					WhereOr("c.contact_name_middle ILIKE ?", "%"+filter.Fullname+"%").
					WhereOr("c.contact_name_family ILIKE ?", "%"+filter.Fullname+"%")
			})
		}
	}
	if len(filter.UnitUuid) > 0 {
		query.Where("u.unit_uuid = ?", filter.UnitUuid)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return extensions, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return extensions, total, nil
	}
}
func (repo *Extension) GetExtensionInfoByExtensionUuid(ctx context.Context, domainUuid, extensionUuid string) (*model.ExtensionInfoWithPassword, error) {
	extension := new(model.ExtensionInfoWithPassword)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(extension).
		Column("d.domain_name", "d.domain_uuid").
		ColumnExpr("e.extension, e.extension_uuid, e.password, e.enabled").
		Column("u.user_uuid", "u.username").
		ColumnExpr("e.follow_me_uuid").
		ColumnExpr(`
			CASE
				WHEN a.call_center_agent_uuid IS NOT NULL THEN true
				ELSE false
			END AS is_link_call_center
		`).
		Join("INNER JOIN v_domains d ON d.domain_uuid = e.domain_uuid").
		Join("LEFT JOIN v_extension_users eu ON e.extension_uuid = eu.extension_uuid").
		Join("LEFT JOIN v_users u ON u.user_uuid = eu.user_uuid").
		Join("LEFT JOIN v_call_center_agents a ON a.user_uuid = u.user_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("e.extension_uuid::text = ?", extensionUuid).
				WhereOr("e.extension = ?", extensionUuid)
		}).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return extension, nil
	}
}

func (repo *Extension) GetExtensionByIdOrExten(ctx context.Context, domainUuid string, extensionUuid string) (*model.Extension, error) {
	extension := new(model.Extension)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(extension).
		ColumnExpr("d.domain_name").
		ColumnExpr("e.extension, e.extension_uuid, e.password").
		ColumnExpr("e.follow_me_uuid").
		Join("INNER JOIN v_domains d ON d.domain_uuid = e.domain_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("e.extension_uuid::text = ?", extensionUuid).
				WhereOr("e.extension = ?", extensionUuid)
		}).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return extension, nil
	}
}

func (repo *Extension) InsertExtensionTransaction(ctx context.Context, extension *model.Extension, extensionUser *model.ExtensionUser) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(extension).Exec(ctx); err != nil {
			return err
		}
		if extensionUser != nil {
			if _, err := tx.NewInsert().Model(extensionUser).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Extension) UpdateExtensionTransaction(ctx context.Context, extension *model.Extension, extensionUser *model.ExtensionUser) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(extension).
			Where("extension_uuid = ?", extension.ExtensionUuid).
			Column("password", "enabled").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ExtensionUser)(nil)).Where("extension_uuid = ?", extension.ExtensionUuid).
			Exec(ctx); err != nil {
			return err
		}
		if extensionUser != nil {
			if _, err := tx.NewInsert().Model(extensionUser).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Extension) DeleteExtensionTransaction(ctx context.Context, domainUuid, extensionUuid string) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*model.Extension)(nil)).Where("extension_uuid = ?", extensionUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ExtensionUser)(nil)).Where("extension_uuid = ?", extensionUuid).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (repo *Extension) UpdateExtensionFollowMeTransaction(ctx context.Context, extension *model.Extension, followMe *model.FollowMe, followMeDestination *model.FollowMeDestination) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(extension).
			Where("extension_uuid = ?", extension.ExtensionUuid).
			Column("dial_domain", "do_not_disturb", "forward_all_enabled", "forward_busy_enabled", "forward_no_answer_enabled", "forward_user_not_registered_enabled",
				"follow_me_uuid", "dial_string").
			Exec(ctx); err != nil {
			return err
		}
		if existed, err := tx.NewSelect().Model((*model.FollowMe)(nil)).
			Where("follow_me_uuid = ?", followMe.FollowMeUuid).
			Exists(ctx); err != nil {
			return err
		} else if !existed {
			if _, err = tx.NewInsert().Model(followMe).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			if _, err := tx.NewUpdate().
				Model(followMe).
				Column("dial_string").
				Where("follow_me_uuid = ?", followMe.FollowMeUuid).
				Exec(ctx); err != nil {
				return err
			}
		}
		if existed, err := tx.NewSelect().Model((*model.FollowMeDestination)(nil)).
			Where("follow_me_uuid = ?", followMeDestination.FollowMeUuid).
			Where("follow_me_destination_uuid= ?", followMeDestination.FollowMeDestinationUuid).
			Exists(ctx); err != nil {
			return err
		} else if !existed {
			if _, err = tx.NewInsert().Model(followMeDestination).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			if _, err := tx.NewUpdate().
				Model(followMeDestination).
				Column("follow_me_destination", "follow_me_timeout", "follow_me_delay", "follow_me_prompt").
				Where("follow_me_uuid = ?", followMeDestination.FollowMeUuid).
				Where("follow_me_destination_uuid = ?", followMeDestination.FollowMeDestinationUuid).
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Extension) DeleteExtensionFollowMeTransaction(ctx context.Context, extensionUuid string, followMeUuid string) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model((*model.Extension)(nil)).
			Where("extension_uuid = ?", extensionUuid).
			Set("follow_me_uuid = NULL").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.FollowMe)(nil)).
			Where("follow_me_uuid = ?", followMeUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewSelect().Model((*model.FollowMeDestination)(nil)).
			Where("follow_me_uuid = ?", followMeUuid).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (repo *Extension) UpdateExtensionRingGroupTransaction(ctx context.Context, dialplans []model.Dialplan, dialplanDetailsMap map[string][]model.DialplanDetail, ringGroups []model.RingGroup, ringGroupDestinations []model.RingGroupDestination) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, dialplan := range dialplans {
			if isExisted, err := tx.NewSelect().Model((*model.Dialplan)(nil)).
				Where("dialplan_number = ?", dialplan.DialplanNumber).
				Where("domain_uuid = ?", dialplan.DomainUuid).
				Exists(ctx); err != nil {
				return err
			} else if !isExisted {
				if _, err = tx.NewInsert().Model(&dialplan).
					Exec(ctx); err != nil {
					return err
				}
			} else {
				if _, err := tx.NewUpdate().
					Model(&dialplan).
					Column("dialplan_name", "dialplan_xml").
					Where("dialplan_number = ?", dialplan.DialplanNumber).
					Where("domain_uuid = ?", dialplan.DomainUuid).
					Exec(ctx); err != nil {
					return err
				}
			}
		}
		for dialplanUuid, dialplanDetails := range dialplanDetailsMap {
			if _, err := repository.FusionSqlClient.GetDB().NewDelete().Model((*model.DialplanDetail)(nil)).
				Where("dialplan_uuid = ?", dialplanUuid).
				Exec(ctx); err != nil {
				return err
			}
			if _, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&dialplanDetails).
				Exec(ctx); err != nil {
				return err
			}
		}
		for _, ringGroup := range ringGroups {
			if isExisted, err := tx.NewSelect().Model((*model.RingGroup)(nil)).
				Where("ring_group_extension = ?", ringGroup.RingGroupExtension).
				Where("domain_uuid = ?", ringGroup.DomainUuid).
				Exists(ctx); err != nil {
				return err
			} else if !isExisted {
				if _, err = tx.NewInsert().Model(&ringGroup).
					Exec(ctx); err != nil {
					return err
				}
			} else {
				if _, err := tx.NewUpdate().
					Model(&ringGroup).
					Column("ring_group_name", "ring_group_enabled", "ring_group_timeout_app", "ring_group_timeout_data").
					Where("ring_group_extension = ?", ringGroup.RingGroupExtension).
					Where("domain_uuid = ?", ringGroup.DomainUuid).
					Exec(ctx); err != nil {
					return err
				}
			}
		}
		for _, ringGroupDestination := range ringGroupDestinations {
			if isExisted, err := tx.NewSelect().Model((*model.RingGroupDestination)(nil)).
				Where("ring_group_uuid = ?", ringGroupDestination.RingGroupUuid).
				Where("destination_number = ?", ringGroupDestination.DestinationNumber).
				Exists(ctx); err != nil {
				return err
			} else if !isExisted {
				if _, err = tx.NewInsert().Model(&ringGroupDestination).
					Exec(ctx); err != nil {
					return err
				}
			} else {
				if _, err := tx.NewUpdate().
					Model(&ringGroupDestination).
					Column("destination_delay", "destination_timeout", "destination_prompt").
					Where("ring_group_uuid = ?", ringGroupDestination.RingGroupUuid).
					Where("destination_number = ?", ringGroupDestination.DestinationNumber).
					Exec(ctx); err != nil {
					return err
				}
			}
		}

		return nil
	})
	return err
}

func (repo *Extension) FindExtensionDataOfExten(ctx context.Context, domainUuid string, exten string) (*model.ExtensionData, error) {
	extension := new(model.ExtensionData)
	err := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(extension).
		ColumnExpr("DISTINCT e.extension_uuid").
		Column("e.extension", "e.domain_uuid", "eu.user_uuid").
		Join("LEFT JOIN v_extension_users eu ON eu.extension_uuid = e.extension_uuid").
		Where("e.domain_uuid = ?", domainUuid).
		Where("e.extension = ?", exten).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return extension, nil
}
