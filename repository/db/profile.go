package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"context"
)

type (
	Profile struct {
	}
)

func NewProfile() repository.IProfile {
	repo := &Profile{}
	repo.InitTable()
	repo.InitColumn()
	repo.InitIndex()

	return repo
}

func (repo *Profile) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repository.FusionSqlClient.GetDB().RegisterModel((*model.GroupUser)(nil))

	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.Profile)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ProfileChannel)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ProfilePhone)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ProfileEmail)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ProfileOwner)(nil)); err != nil {
		panic(err)
	}
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.ProfileNote)(nil)); err != nil {
		panic(err)
	}
}

func (repo *Profile) InitColumn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := repository.FusionSqlClient.GetDB().NewAddColumn().Model((*model.Profile)(nil)).IfNotExists().ColumnExpr("social_mapping_contact text").Exec(ctx); err != nil {
		panic(err)
	}
}

func (repo *Profile) InitIndex() {

}

func (repo *Profile) GetProfileInfoById(ctx context.Context, domainUuid string, profileId string) (*model.ProfileView, error) {
	profile := new(model.ProfileView)
	err := repository.FusionSqlClient.GetDB().NewSelect().Model(profile).
		Relation("Emails").
		Relation("Phones").
		Relation("UserOwners").
		Relation("ListRelatedProfile").
		Where("p.domain_uuid = ?", domainUuid).
		Where("p.profile_uuid = ?", profileId).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return profile, err
}

func (repo *Profile) GetProfilesInfo(ctx context.Context, domainUuid string, filter model.ProfileFilter, limit, offset int) (int, *[]model.ProfileView, error) {
	profiles := new([]model.ProfileView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(profiles).
		Relation("Emails").
		Relation("Phones").
		Relation("UserOwners").
		Relation("ListRelatedProfile").
		Where("p.domain_uuid = ?", domainUuid).
		Order("p.created_at DESC")
	if len(filter.Common) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("? ILIKE ?", bun.Ident("p.fullname"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("p.email"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("p.phone_number"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("p.profile_code"), "%"+filter.Common+"%")
		})
	} else {
		if len(filter.PhoneNumber) > 0 {
			query.Where("p.phone_number = ?", filter.PhoneNumber)
		}
		if len(filter.Email) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.email"), "%"+filter.Email+"%")
		}
		if len(filter.Fullname) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.fullname"), "%"+filter.Fullname+"%")
		}
		if len(filter.ProfileCode) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.profile_code"), "%"+filter.ProfileCode+"%")
		}
		if len(filter.JobTitle) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.job_title"), "%"+filter.JobTitle+"%")
		}
		if len(filter.Birthday) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.birthday"), "%"+filter.Birthday+"%")
		}
		if len(filter.Gender) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.gender"), filter.Gender)
		}
		if len(filter.Address) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.address"), "%"+filter.Address+"%")
		}
		if len(filter.Country) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.country"), "%"+filter.Country+"%")
		}
		if len(filter.Province) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.province"), "%"+filter.Province+"%")
		}
		if len(filter.District) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.district"), "%"+filter.District+"%")
		}
		if len(filter.Ward) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.ward"), "%"+filter.Ward+"%")
		}
		if filter.Status.Valid {
			query.Where("p.status = ?", filter.Status.Bool)
		}
		if len(filter.RefCode) > 0 {
			query = query.Where("p.ref_code LIKE ?", "%"+filter.RefCode+"%")
		}
		if len(filter.RefId) > 0 {
			query = query.Where("p.ref_id LIKE ?", "%"+filter.RefId+"%")
		}
		if len(filter.UserOwerUuid) > 0 {
			subQuery := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("profile_owner as po").
				Where("p.profile_uuid = po.profile_uuid").
				Where("po.profile_uuid= ?", filter.UserOwerUuid)
			query = query.Where("EXISTS (?)", subQuery)
		}
		if len(filter.IdentityNumber) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.identity_number"), "%"+filter.IdentityNumber+"%")
		}
	}

	if len(filter.Passport) > 0 {
		query = query.Where("? ILIKE ?", bun.Ident("p.passport"), "%"+filter.Passport+"%")
	}
	if len(filter.ProfileType) > 0 {
		query = query.Where("profile_type = ?", filter.ProfileType)
	}

	if !filter.StartTime.IsZero() {
		query = query.Where("p.created_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("p.created_at <= ?", filter.EndTime)
	}

	if len(filter.FacebookUserId) > 0 {
		query.Where("social_mapping_contact::jsonb ->>'facebook' = ?", filter.FacebookUserId)
	}

	if len(filter.ZaloUserId) > 0 {
		query.Where("social_mapping_contact::jsonb ->>'zalo' = ?", filter.ZaloUserId)
	}

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, profiles, nil
	}
	return total, profiles, err
}

func (repo *Profile) GetProfileByPhoneNumber(ctx context.Context, domainUuid string, phoneNumber ...string) (*model.ProfileView, error) {
	profile := new(model.ProfileView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(profile).
		Where("p.domain_uuid = ?", domainUuid).
		Where("p.phone_number IN (?)", bun.In(phoneNumber)).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return profile, err
}

func (repo *Profile) InsertProfileTransaction(ctx context.Context, profile *model.Profile, profilePhones []model.ProfilePhone, profileEmails []model.ProfileEmail, profileOwners []model.ProfileOwner, profileNotes []model.ProfileNote) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(profile).Exec(ctx); err != nil {
			return err
		}
		if len(profilePhones) > 0 {
			if _, err := tx.NewInsert().Model(&profilePhones).Exec(ctx); err != nil {
				return err
			}
		}
		if len(profileEmails) > 0 {
			if _, err := tx.NewInsert().Model(&profileEmails).Exec(ctx); err != nil {
				return err
			}
		}
		if len(profileOwners) > 0 {
			if _, err := tx.NewInsert().Model(&profileOwners).Exec(ctx); err != nil {
				return err
			}
		}
		if len(profileNotes) > 0 {
			if _, err := tx.NewInsert().Model(&profileNotes).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Profile) UpdateProfile(ctx context.Context, profile *model.Profile) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(profile).
		Where("profile_uuid = ?", profile.ProfileUuid)
	_, err := query.Exec(ctx)
	return err
}

func (repo *Profile) UpdateProfileTransaction(ctx context.Context, profileUuid string, profile *model.Profile, profilePhones []model.ProfilePhone, profileEmails []model.ProfileEmail, profileOwners []model.ProfileOwner) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(profile).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ProfilePhone)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if len(profilePhones) > 0 {
			if _, err := tx.NewInsert().Model(&profilePhones).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ProfileEmail)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if len(profileEmails) > 0 {
			if _, err := tx.NewInsert().Model(&profileEmails).Exec(ctx); err != nil {
				return err
			}
		}
		if _, err := tx.NewDelete().Model((*model.ProfileOwner)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if len(profileOwners) > 0 {
			if _, err := tx.NewInsert().Model(&profileOwners).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *Profile) GetProfileById(ctx context.Context, domainUuid, profileId string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := repository.FusionSqlClient.GetDB().NewSelect().Model(profile).
		Where("domain_uuid = ?", domainUuid).
		Where("profile_uuid = ?", profileId).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return profile, err
}

func (repo *Profile) DeleteProfileTransaction(ctx context.Context, profileUuid string) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*model.Profile)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ProfilePhone)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ProfileEmail)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ProfileOwner)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ProfileNote)(nil)).Where("profile_uuid = ?", profileUuid).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}

func (repo *Profile) UpdateProfileByField(ctx context.Context, domainUuid string, profile model.Profile) error {
	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(&profile).
		Set("social_mapping_contact = ?", profile.SocialMappingContact).
		WherePK("profile_uuid")
	_, err := query.Exec(ctx)
	return err
}

func (repo *Profile) GetManageProfiles(ctx context.Context, domainUuid string, filter model.ProfileFilter, limit, offset int) (int, *[]model.ProfileManageView, error) {
	profiles := new([]model.ProfileManageView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(profiles).
		ColumnExpr("p.*").
		ColumnExpr("tk.*").
		ColumnExpr("primary_profile.*").
		ColumnExpr("CONCAT(c.contact_name_given, c.contact_name_middle, c.contact_name_family) as user_created_by").
		Join("LEFT JOIN v_users u ON p.created_by::uuid = u.user_uuid").
		Join("LEFT JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Where("p.domain_uuid = ?", domainUuid).
		Order("p.created_at DESC")
	if len(filter.Common) > 0 {
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("? ILIKE ?", bun.Ident("p.fullname"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("p.email"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("p.phone_number"), "%"+filter.Common+"%").
				WhereOr("? ILIKE ?", bun.Ident("p.profile_code"), "%"+filter.Common+"%")
		})
	} else {
		if len(filter.PhoneNumber) > 0 {
			query.Where("p.phone_number = ?", filter.PhoneNumber)
		}
		if len(filter.Email) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.email"), "%"+filter.Email+"%")
		}
		if len(filter.Fullname) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.fullname"), "%"+filter.Fullname+"%")
		}
		if len(filter.ProfileCode) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.profile_code"), "%"+filter.ProfileCode+"%")
		}
		if len(filter.JobTitle) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.job_title"), "%"+filter.JobTitle+"%")
		}
		if len(filter.Birthday) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.birthday"), "%"+filter.Birthday+"%")
		}
		if len(filter.Gender) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.gender"), filter.Gender)
		}
		if len(filter.Address) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.address"), "%"+filter.Address+"%")
		}
		if len(filter.Country) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.country"), "%"+filter.Country+"%")
		}
		if len(filter.Province) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.province"), "%"+filter.Province+"%")
		}
		if len(filter.District) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.district"), "%"+filter.District+"%")
		}
		if len(filter.Ward) > 0 {
			query.Where("? ILIKE ?", bun.Ident("p.ward"), "%"+filter.Ward+"%")
		}
		if filter.Status.Valid {
			query.Where("p.status = ?", filter.Status.Bool)
		}
		if len(filter.RefCode) > 0 {
			query = query.Where("p.ref_code LIKE ?", "%"+filter.RefCode+"%")
		}
		if len(filter.RefId) > 0 {
			query = query.Where("p.ref_id LIKE ?", "%"+filter.RefId+"%")
		}
		if len(filter.UserOwerUuid) > 0 {
			subQuery := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("profile_owner as po").
				Where("p.profile_uuid = po.profile_uuid").
				Where("po.profile_uuid = ?", filter.UserOwerUuid)
			query = query.Where("EXISTS (?)", subQuery)
		}
		if len(filter.IdentityNumber) > 0 {
			query = query.Where("? ILIKE ?", bun.Ident("p.identity_number"), "%"+filter.IdentityNumber+"%")
		}
	}

	if len(filter.Passport) > 0 {
		query = query.Where("? ILIKE ?", bun.Ident("p.passport"), "%"+filter.Passport+"%")
	}
	if len(filter.ProfileType) > 0 {
		query = query.Where("profile_type = ?", filter.ProfileType)
	}

	if !filter.StartTime.IsZero() {
		query = query.Where("p.created_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("p.created_at <= ?", filter.EndTime)
	}

	if len(filter.FacebookUserId) > 0 {
		query.Where("social_mapping_contact::jsonb ->>'facebook' = ?", filter.FacebookUserId)
	}

	if len(filter.ZaloUserId) > 0 {
		query.Where("social_mapping_contact::jsonb ->>'zalo' = ?", filter.ZaloUserId)
	}

	if len(filter.ProfileUuids) > 0 {
		query.Where("p.profile_uuid IN (?)", bun.In(filter.ProfileUuids))
	}

	if len(filter.CreatedBy) > 0 {
		query.Where("p.created_by = ?", filter.CreatedBy)
	}

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	subQuery := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("ticket as tk").
		ColumnExpr(`COUNT(CASE WHEN tk.status='OPEN' THEN tk.ticket_uuid END) AS total_open`).
		ColumnExpr(`COUNT(CASE WHEN tk.status='PROCESSING' THEN tk.ticket_uuid END) AS total_processing`).
		ColumnExpr(`COUNT(CASE WHEN tk.status='WAITING' THEN tk.ticket_uuid END) AS total_waiting`).
		ColumnExpr(`COUNT(CASE WHEN tk.status='SOLVED' THEN tk.ticket_uuid END) AS total_solved`).
		ColumnExpr(`COUNT(CASE WHEN tk.status='PENDING' THEN tk.ticket_uuid END) AS total_pending`).
		ColumnExpr(`COUNT(CASE WHEN tk.status='REOPEN' THEN tk.ticket_uuid END) AS total_reopen`).
		Where("tk.profile_uuid=p.profile_uuid")
	if !filter.StartTime.IsZero() {
		subQuery.Where("tk.created_at >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		subQuery.Where("tk.created_at <= ?", filter.EndTime)
	}

	subQuery2 := repository.FusionSqlClient.GetDB().NewSelect().TableExpr("profile as pp").
		ColumnExpr("pp.profile_uuid as primary_profile_uuid, pp.fullname as primary_profile_fullname, pp.profile_type as primary_profile_type, pp.profile_name as primary_profile_name").
		Where("pp.profile_uuid=p.related_profile_uuid")

	query.Join("LEFT JOIN LATERAL (?) as tk ON true", subQuery).
		Join("LEFT JOIN LATERAL (?) as primary_profile ON true", subQuery2)

	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, profiles, nil
	}
	return total, profiles, err
}

func (repo *Profile) DeleteProfile(ctx context.Context, profile []model.Profile) error {
	_, err := repository.FusionSqlClient.GetDB().NewDelete().
		Model(&profile).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Profile) DeleteProfileWithTicketTransaction(ctx context.Context, domainUuid string, profiles []model.Profile, tickets []model.Ticket) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if len(profiles) > 0 {
			if _, err := repository.FusionSqlClient.GetDB().NewDelete().Model(&profiles).WherePK().Exec(ctx); err != nil {
				return err
			}
		}

		if len(tickets) > 0 {
			if _, err := repository.FusionSqlClient.GetDB().NewDelete().Model(&tickets).WherePK().Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}
