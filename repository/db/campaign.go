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

type Campaign struct {
}

func NewCampaign() repository.ICampaign {
	repo := &Campaign{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()
	return repo
}

func (repo *Campaign) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Campaign)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.CampaignUser)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.CampaignGroup)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.CampaignSchedule)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.CampaignCustomData)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Campaign) InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_campaign_name").Column("campaign_name").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_domain_uuid").Column("domain_uuid").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_type").Column("type").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_local_start_time").Column("local_start_time").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_local_end_time").Column("local_end_time").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_created_at").Column("created_at").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewCreateIndex().Model((*model.Campaign)(nil)).IfNotExists().Index("idx_campaign_active").Column("active").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewAddColumn().Model((*model.Campaign)(nil)).IfNotExists().ColumnExpr("enable_encrypt bool DEFAULT 'false'").Exec(ctx); err != nil {
		log.Error(err)
	}
	if _, err := repository.FusionSqlClient.GetDB().NewAddColumn().Model((*model.Campaign)(nil)).IfNotExists().ColumnExpr("orig_campaign_uuid uuid NULL DEFAULT NULL").Exec(ctx); err != nil {
		log.Error(err)
	}
}

func (repo *Campaign) InitColumn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "status text NULL DEFAULT NULL"); err != nil {
	// 	log.Error(err)
	// }
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "mode_call text NULL DEFAULT NULL"); err != nil {
	// 	log.Error(err)
	// }
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "network text NULL DEFAULT NULL"); err != nil {
	// 	log.Error(err)
	// }
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "run_id text NULL DEFAULT NULL"); err != nil {
	// 	log.Error(err)
	// }
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "type_autocall text NULL DEFAULT NULL"); err != nil {
	// 	log.Error(err)
	// }
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "priority_recall text NULL DEFAULT NULL"); err != nil {
	// 	log.Error(err)
	// }
	// if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "call_timeout integer NULL DEFAULT 60"); err != nil {
	// 	log.Error(err)
	// }
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "callback_url text NULL DEFAULT NULL"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "lead_expired_at timestamp NULL DEFAULT NULL"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "allow_select_status boolean NULL DEFAULT false"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "hot_lead_inner_campaign_uuids text[] DEFAULT NULL"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "is_contract_campaign boolean NULL DEFAULT false"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.Campaign)(nil), "created_by char(36) NULL DEFAULT NULL"); err != nil {
		log.Error(err)
	}
}

func (repo *Campaign) GetCampaignsOptionView(ctx context.Context, domainUuid string, filter model.CampaignFilter, limit, offset int) (*[]model.CampaignOptionView, int, error) {
	campaigns := new([]model.CampaignOptionView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		Where("domain_uuid = ?", domainUuid).
		Limit(limit).
		Offset(offset)
	if len(filter.ManageCampaignUuids) > 0 {
		query = query.Where("campaign_uuid IN (?)", bun.In(filter.ManageCampaignUuids))
	}
	count, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return campaigns, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return campaigns, count, nil
	}
}

func (repo *Campaign) GetCampaignsInfo(ctx context.Context, domainUuid string, filter model.CampaignFilter, limit, offset int) (*[]model.CampaignView, int, error) {
	campaigns := new([]model.CampaignView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		ColumnExpr("ca.*").
		ColumnExpr("(?) as created_by_username", repository.FusionSqlClient.GetDB().
			NewSelect().TableExpr("v_users as u").
			Column("u.username").
			Where("u.user_uuid = ca.created_by").
			Limit(1)).
		ColumnExpr("ccq.queue_strategy as call_center_queue_strategy").
		ColumnExpr("im.ivr_menu_name as template_name").
		ColumnExpr("carrier.carrier_name").
		Join("LEFT JOIN v_call_center_queues ccq ON ccq.call_center_queue_uuid = ca.call_center_queue_uuid").
		Join("LEFT JOIN v_ivr_menus im ON im.ivr_menu_uuid = ca.template_uuid").
		Join("LEFT JOIN carrier ON ca.carrier_uuid = carrier.carrier_uuid").
		Where("ca.domain_uuid = ?", domainUuid).
		Order("created_at DESC")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	if len(filter.CampaignName) > 0 {
		query.Where("ca.campaign_name LIKE ?", "%"+filter.CampaignName+"%")
	}
	if len(filter.Types) > 0 {
		query.Where("ca.type IN (?)", bun.In(filter.Types))
	}
	if len(filter.ManageCampaignUuids) > 0 {
		query.Where("ca.campaign_uuid IN (?)", bun.In(filter.ManageCampaignUuids))
	}
	count, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return campaigns, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return campaigns, count, nil
	}
}

func (repo *Campaign) GetCampaigns(ctx context.Context, domainUuid string, filter model.GeneralFilter, limit, offset int) (*[]model.Campaign, int, error) {
	campaigns := new([]model.Campaign)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		Where("domain_uuid = ?", domainUuid).
		Order("created_at DESC")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	if len(filter.ManageCampaignUuids) > 0 {
		query.Where("campaign_uuid IN (?)", bun.In(filter.ManageCampaignUuids))
	}
	if len(filter.CampaignUuids) > 0 {
		query.Where("campaign_uuid IN (?)", bun.In(filter.CampaignUuids))
	}
	if len(filter.CampaignTypes) > 0 {
		query.Where("type IN (?)", bun.In(filter.CampaignTypes))
	}
	if len(filter.CarrierUuids) > 0 {
		query.Where("carrier_uuid IN (?)", bun.In(filter.CarrierUuids))
	}
	count, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return campaigns, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return campaigns, count, nil
	}
}

func (repo *Campaign) GetCampaignInfoById(ctx context.Context, domainUuid, campaignUuid string) (*model.CampaignView, error) {
	campaign := new(model.CampaignView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaign).
		Relation("Users", func(q *bun.SelectQuery) *bun.SelectQuery {
			q = q.
				Column("u.user_uuid", "e.extension", "u.domain_uuid").
				Join("LEFT JOIN v_extension_users eu ON u.user_uuid = eu.user_uuid").
				Join("LEFT JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid")
			return q
		}).
		Relation("Groups").
		Relation("Statuses", func(q *bun.SelectQuery) *bun.SelectQuery {
			q = q.
				ColumnExpr("s.status_uuid, s.status_code, s.status_name, s.campaign_uuid").
				ColumnExpr("s.selectable, s.scheduled_callback").
				ColumnExpr("s.category_status_uuid, cs.category_status_name, cs.category_status_code").
				Join("INNER JOIN category_status cs ON s.category_status_uuid = cs.category_status_uuid").
				Order("s.status_code")
			return q
		}).
		Relation("Schedules").
		ColumnExpr("ca.*").
		ColumnExpr("ccq.queue_strategy as call_center_queue_strategy").
		Join("LEFT JOIN v_call_center_queues ccq ON ccq.call_center_queue_uuid = ca.call_center_queue_uuid").
		Where("ca.domain_uuid = ?", domainUuid).
		Where("ca.campaign_uuid = ?", campaignUuid).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return campaign, nil
	}
}

func (repo *Campaign) GetCampaignsStatusInfo(ctx context.Context, domainUuid string, limit, offset int) (*[]model.CampaignStatusView, int, error) {
	campaigns := new([]model.CampaignStatusView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		Relation("Statuses", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Join("INNER JOIN category_status ON category_status.category_status_uuid = status.category_status_uuid").Order("status_code")
		}).
		ColumnExpr("ca.*").
		Where("ca.domain_uuid = ?", domainUuid)
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	count, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return campaigns, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		return campaigns, count, nil
	}
}

func (repo *Campaign) GetCampaignById(ctx context.Context, domainUuid string, campaignUuid string) (*model.Campaign, error) {
	campaign := model.Campaign{}
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&campaign).
		Where("domain_uuid = ?", domainUuid).
		Where("campaign_uuid = ?", campaignUuid)

	err := query.Scan(ctx, &campaign)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &campaign, nil
	}
}

func (repo *Campaign) GetCampaignByName(ctx context.Context, domainUuid string, campaignName string) (*model.Campaign, error) {
	campaign := model.Campaign{}
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(&campaign).
		Where("domain_uuid = ?", domainUuid).
		Where("campaign_name = ?", campaignName).
		Limit(1)

	err := query.Scan(ctx, &campaign)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &campaign, nil
	}
}

func (repo *Campaign) InsertCampaign(ctx context.Context, campaign *model.Campaign) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(campaign)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert campaign failed")
	}
	return nil
}

func (repo *Campaign) InsertCampaignAutodialerTransaction(ctx context.Context, campaign *model.Campaign, callCenterQueue *model.CallCenterQueue, dialplan *model.Dialplan) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(campaign).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(callCenterQueue).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(dialplan).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (repo *Campaign) InsertCampaignUsers(ctx context.Context, campaignUuid string, campaignUsers *[]model.CampaignUser) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(&model.CampaignUser{}).Where("campaign_uuid = ?", campaignUuid).Exec(ctx); err != nil {
			return err
		}
		if len(*campaignUsers) > 0 {
			if _, err := tx.NewInsert().Model(campaignUsers).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Campaign) InsertCampaignGroups(ctx context.Context, campaignUuid string, campaignGroups *[]model.CampaignGroup) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(&model.CampaignGroup{}).Where("campaign_uuid = ?", campaignUuid).Exec(ctx); err != nil {
			return err
		}
		if len(*campaignGroups) > 0 {
			if _, err := tx.NewInsert().Model(campaignGroups).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Campaign) AppendCampaignGroups(ctx context.Context, campaignUuid string, campaignGroups ...model.CampaignGroup) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, cg := range campaignGroups {
			tmp := cg
			existed, err := repository.FusionSqlClient.GetDB().NewSelect().Model(&cg).
				Where("campaign_uuid = ?", cg.CampaignUuid).
				Where("group_uuid = ?", cg.GroupUuid).
				Exists(ctx)
			if err != nil {
				return err
			} else if !existed {
				cg = tmp
				resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&cg).
					Exec(ctx)
				if err != nil {
					return err
				} else if affected, _ := resp.RowsAffected(); affected < 0 {
					return errors.New("insert fail")
				}
			}
		}
		return nil
	})
	return err
}

func (repo *Campaign) AppendCampaignUsers(ctx context.Context, campaignUuid string, campaignUsers ...model.CampaignUser) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, cg := range campaignUsers {
			tmp := cg
			existed, err := repository.FusionSqlClient.GetDB().NewSelect().Model(&cg).
				Where("campaign_uuid = ?", cg.CampaignUuid).
				Where("user_uuid = ?", cg.UserUuid).
				Exists(ctx)
			if err != nil {
				return err
			} else if !existed {
				cg = tmp
				resp, err := repository.FusionSqlClient.GetDB().NewInsert().Model(&cg).
					Exec(ctx)
				if err != nil {
					return err
				} else if affected, _ := resp.RowsAffected(); affected < 0 {
					return errors.New("insert fail")
				}
			}
		}
		return nil
	})
	return err
}

func (repo *Campaign) GetCampaignUserByUserUuid(ctx context.Context, campaignUuid string, userUuid string) (*model.CampaignUser, error) {
	campaignUser := new(model.CampaignUser)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaignUser).
		Where("campaign_uuid = ?", campaignUuid).
		Where("user_uuid = ?", userUuid).
		Limit(1)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return campaignUser, nil
	}
}

func (repo *Campaign) GetCampaignUsersByCampaignUuid(ctx context.Context, campaignUuid string) (*[]model.CampaignUser, error) {
	campaignUsers := new([]model.CampaignUser)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaignUsers).
		Where("campaign_uuid = ?", campaignUuid)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return campaignUsers, nil
}

func (repo *Campaign) GetCampaignsActive(ctx context.Context) (*[]model.CampaignView, error) {
	campaigns := new([]model.CampaignView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		Column("campaign_uuid", "type").
		Where("active = ?", true)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return campaigns, nil
	} else if err != nil {
		return nil, err
	} else {
		return campaigns, nil
	}
}

func (repo *Campaign) GetCampaignsActiveContainHopper(ctx context.Context) (*[]model.CampaignView, error) {
	campaigns := new([]model.CampaignView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		ColumnExpr("ca.*").
		ColumnExpr("domain.domain_name").
		Join("INNER JOIN v_domains as domain ON domain.domain_uuid = ca.domain_uuid").
		Join("LEFT JOIN hopper ON ca.campaign_uuid = hopper.campaign_uuid").
		Where("ca.active = ?", true).
		Group("domain.domain_uuid", "ca.campaign_uuid").
		Order("ca.created_at DESC").
		Having("count(hopper.id) > 0")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return campaigns, nil
	} else if err != nil {
		return nil, err
	} else {
		return campaigns, nil
	}
}

func (repo *Campaign) UpdateCampaignTransaction(ctx context.Context, campaign *model.Campaign, campaignUsers *[]model.CampaignUser, statuses []string) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(campaign).Where("campaign_uuid = ?", campaign.Id).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.CampaignUser)(nil)).Where("campaign_uuid = ?", campaign.Id).Exec(ctx); err != nil {
			return err
		}
		if len(*campaignUsers) > 0 {
			if _, err := tx.NewInsert().Model(campaignUsers).Exec(ctx); err != nil {
				return err
			}
		}
		if len(statuses) > 0 {
			//if _, err := tx.NewUpdate().Model((*model.Status)(nil)).Where("status_uuid NOT IN (?)", bun.In(statuses)).Set("selectable = false").Exec(ctx); err != nil {
			//	return err
			//}
			if _, err := tx.NewUpdate().Model((*model.Status)(nil)).Where("status_uuid IN (?)", bun.In(statuses)).Set("selectable = true").Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Campaign) UpdateCampaign(ctx context.Context, domainUuid string, campaign *model.Campaign) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().
		Model(campaign).
		Where("campaign_uuid = ?", campaign.Id).
		Where("domain_uuid = ?", domainUuid)
	if _, err := query.
		Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *Campaign) DeleteCampaign(ctx context.Context, domainUuid string, campaignUuid string) error {
	if _, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model((*model.Campaign)(nil)).
		Where("campaign_uuid = ?", campaignUuid).
		Where("domain_uuid = ?", domainUuid).
		Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *Campaign) GetInboundCampaignByQueueExtension(ctx context.Context, domainUuid, queueExtension string) (*model.Campaign, error) {
	campaign := new(model.Campaign)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(campaign).
		Join("INNER JOIN v_call_center_queues ccq ON ccq.call_center_queue_uuid = ca.call_center_queue_uuid").
		Where("ccq.queue_extension = ?", queueExtension).
		Where("domain_uuid = ?", domainUuid).
		Where("type = ?", "inbound").
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return campaign, nil
	}
}

func (repo *Campaign) DeleteCampaignGroup(ctx context.Context, domainUuid string, campaignUuid string) error {
	if _, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model((*model.CampaignGroup)(nil)).
		Where("campaign_uuid = ?", campaignUuid).
		Where("domain_uuid = ?", domainUuid).
		Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *Campaign) DeleteCampaignUser(ctx context.Context, domainUuid string, campaignUuid string) error {
	if _, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model((*model.CampaignUser)(nil)).
		Where("campaign_uuid = ?", campaignUuid).
		Where("domain_uuid = ?", domainUuid).
		Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *Campaign) InsertCampaignSchedules(ctx context.Context, campaignUuid string, schedules ...model.CampaignSchedule) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(&model.CampaignSchedule{}).Where("campaign_uuid = ?", campaignUuid).Exec(ctx); err != nil {
			return err
		}
		if len(schedules) > 0 {
			if _, err := tx.NewInsert().Model(&schedules).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (repo *Campaign) GetCampaignUuidsOfUsers(ctx context.Context, domainUuid string, userUuids []string) ([]string, error) {
	campaignUuids := make([]string, 0)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		TableExpr("campaign as ca").
		ColumnExpr("DISTINCT ca.campaign_uuid").
		Join("LEFT JOIN campaign_users as cu ON cu.campaign_uuid = ca.campaign_uuid").
		Where("ca.active = ?", true).
		Where("ca.domain_uuid = ?", domainUuid)
	if len(userUuids) > 0 {
		query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			if len(userUuids) > 0 {
				query = query.WhereOr("cu.user_uuid IN (?)", bun.In(userUuids))
			}
			return q
		})
	}
	err := query.Scan(ctx, &campaignUuids)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return campaignUuids, nil
	}
}

func (repo *Campaign) GetCampaignsByCarrierId(ctx context.Context, domainUuid, carrierUuid string) (*[]model.Campaign, error) {
	campaigns := new([]model.Campaign)

	query := repository.FusionSqlClient.GetDB().NewSelect().Model(campaigns).
		ColumnExpr("ca.*").
		Where("ca.active = ?", true).
		Where("ca.domain_uuid = ?", domainUuid).
		Where("ca.carrier_uuid = ?", carrierUuid)

	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (repo *Campaign) GetCampaignCustomDatasByCampaignUuid(ctx context.Context, campaignUuid string) (*[]model.CampaignCustomData, error) {
	campaignCustomDatas := new([]model.CampaignCustomData)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(campaignCustomDatas).
		Where("campaign_uuid = ?", campaignUuid)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return campaignCustomDatas, nil
}

func (repo *Campaign) InsertCampaignCustomData(ctx context.Context, campaignCustomDatas ...model.CampaignCustomData) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, campaignCustomData := range campaignCustomDatas {
			if existed, err := tx.NewSelect().Model((*model.CampaignCustomData)(nil)).
				Where("campaign_uuid = ?", campaignCustomData.CampaignUuid).
				Where("key = ?", campaignCustomData.Key).
				Exists(ctx); err != nil {
				return err
			} else if !existed {
				if _, err = tx.NewInsert().Model(&campaignCustomData).
					Exec(ctx); err != nil {
					return err
				}
			} else {
				if _, err := tx.NewUpdate().
					Model(&campaignCustomData).
					Column("value").
					Where("campaign_uuid = ?", campaignCustomData.CampaignUuid).
					Where("key = ?", campaignCustomData.Key).
					Exec(ctx); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
