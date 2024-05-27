package db

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	sqlclient "contactcenter-api/internal/sqlclient"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
}

func NewUser() repository.IUser {
	repo := &User{}
	repo.InitTable()
	repo.InitColumn()
	return repo
}

func (repo *User) InitTable() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := repository.CreateTable(repository.FusionSqlClient, ctx, (*model.UserCustomData)(nil)); err != nil {
		panic(err)
	}
}

func (repo *User) InitColumn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.User)(nil), "role_uuid char(36) NULL DEFAULT NULL"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.User)(nil), "enable_webrtc bool DEFAULT false"); err != nil {
		log.Error(err)
	}
	if err := repository.AddColumn(repository.FusionSqlClient, ctx, (*model.User)(nil), "last_login_date timestamp NULL DEFAULT NULL"); err != nil {
		log.Error(err)
	}
}

func (repo *User) GetUsersInfo(ctx context.Context, domainUuid string, limit, offset int, filter model.UserFilter) (int, *[]model.UserView, error) {
	users := new([]model.UserView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(users).
		ColumnExpr("u.user_uuid, e.extension_uuid, u.domain_uuid, u.user_enabled, u.user_status").
		ColumnExpr("u.username, e.extension, u.level").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		ColumnExpr("role_group.role_group_uuid, role_group.role_group_name").
		ColumnExpr("unit.unit_uuid, unit.unit_name").
		Relation("Groups").
		// Column("u.role_uuid", "role.role_name").
		// Join("LEFT JOIN role ON role.role_uuid = u.role_uuid").
		Join("LEFT JOIN v_extension_users as eu ON eu.user_uuid = u.user_uuid").
		Join("LEFT JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid").
		Join("LEFT JOIN role_group ON role_group.role_group_uuid = u.role_uuid").
		Join("LEFT JOIN unit ON unit.unit_uuid = u.unit_uuid")
	if len(domainUuid) > 0 && domainUuid != "all" {
		query = query.Where("u.domain_uuid = ?", domainUuid)
	}

	if len(filter.UserUuid) > 0 {
		query = query.Where("u.user_uuid IN (?)", bun.In(filter.UserUuid))
	}

	if len(filter.StartTime) > 0 {
		query = query.Where("u.add_date >= ? ", filter.StartTime)
	}

	if len(filter.EndTime) > 0 {
		query = query.Where("u.add_date <= ? ", filter.EndTime)
	}

	if len(filter.Name) > 0 {
		query = query.Where("u.username LIKE ?", "%"+filter.Name+"%")
	}

	if len(filter.Level) > 0 {
		query = query.
			Where("u.level = ?", filter.Level)
	}
	if len(filter.Levels) > 0 {
		query = query.
			Where("u.level IN (?)", bun.In(filter.Levels))
	}
	if len(filter.Enabled) > 0 {
		query = query.Where("u.user_enabled = ?", filter.Enabled)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if len(filter.Email) > 0 {
		query = query.Where("ce.email_address like ?", "%"+filter.Email+"%")
	}
	if !filter.IsAll {
		if len(filter.ManageUserUuids) > 0 {
			query = query.
				Where("u.user_uuid IN (?)", bun.In(filter.ManageUserUuids))
		}
	}
	if filter.IsMapExtension == "true" {
		query = query.Where("e.extension_uuid IS NOT NULL")
	} else if filter.IsMapExtension == "false" {
		query = query.Where("e.extension_uuid IS NULL")
	}
	if len(filter.RoleUuid) > 0 {
		query = query.Where("u.role_uuid = ?", filter.RoleUuid)
	}
	// if len(filter.GroupUuid) > 0 {
	// 	subQuery := repository.FusionSqlClient.GetDB().NewSelect().
	// 		TableExpr("v_group_users as gu").
	// 		Where("gu.user_uuid = u.user_uuid").
	// 		Where("gu.group_uuid = ?", filter.GroupUuid)
	// 	query = query.Where("EXISTS (?)", subQuery)
	// }
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
		query = query.Order(order)
	}
	total, err := query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, nil, errors.New("no data")
	} else if err != nil {
		return 0, nil, err
	}
	return total, users, nil
}

func (repo *User) GetUserViewById(ctx context.Context, domainUuid, userUuid string) (*model.UserView, error) {
	user := new(model.UserView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.user_uuid, e.extension_uuid, u.domain_uuid, u.user_enabled, u.user_status").
		ColumnExpr("u.username, e.extension, u.level").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		Column("u.role_uuid", "role.role_name").
		Column("u.enable_webrtc").
		Relation("RoleGroups").
		Join("LEFT JOIN role ON role.role_uuid = u.role_uuid").
		Join("LEFT JOIN v_extension_users as eu ON eu.user_uuid = u.user_uuid").
		Join("LEFT JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid ").
		Where("u.user_uuid = ?", userUuid)
	if len(domainUuid) > 0 {
		query = query.Where("u.domain_uuid = ?", domainUuid)
	}
	if err := query.Scan(ctx); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) GetUserViewByIdOrUsername(ctx context.Context, domainUuid, id string) (*model.UserView, error) {
	user := new(model.UserView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.user_uuid, e.extension_uuid, u.domain_uuid, u.user_enabled, u.user_status").
		ColumnExpr("u.username, e.extension, u.level").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		Relation("Groups").
		Column("u.role_uuid", "role.role_name").
		Join("LEFT JOIN role ON role.role_uuid = u.role_uuid").
		Join("LEFT JOIN v_extension_users as eu ON eu.user_uuid = u.user_uuid").
		Join("LEFT JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Join("INNER JOIN v_contacts c ON u.contact_uuid = c.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid ").
		Where("( u.user_uuid::text = ? OR u.username = ? )", id, id)
	if len(domainUuid) > 0 {
		query = query.Where("u.domain_uuid = ?", domainUuid)
	}
	if err := query.Scan(ctx); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) GetUserInfoById(ctx context.Context, userUuid string) (*model.UserView, error) {
	user := new(model.UserView)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(user).
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = u.contact_uuid").
		ColumnExpr("u.user_uuid, u.domain_uuid").
		ColumnExpr("u.username, u.level, u.enable_webrtc").
		ColumnExpr("ce.email_address as email").
		Relation("Groups").
		Where("u.user_uuid = ?", userUuid).
		Column("u.role_uuid")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) GetUserById(ctx context.Context, domainUuid string, userUuid string) (*model.User, error) {
	user := new(model.User)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.*").
		Where("u.domain_uuid = ?", domainUuid).
		Where("u.user_uuid = ?", userUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) GetUserOfDomainByName(ctx context.Context, domainUuid string, userName string) (*model.UserView, error) {
	user := model.UserView{}
	query := repository.FusionSqlClient.GetDB().NewSelect().
		TableExpr("v_users as u").
		ColumnExpr("u.*").
		Where("u.domain_uuid = ?", domainUuid).
		Where("u.username = ?", userName).
		Limit(1)
	err := query.Scan(ctx, &user)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (repo *User) GetUserLoginInfoById(ctx context.Context, domainUuid, userUuid string) (*model.UserLoginView, error) {
	user := new(model.UserLoginView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.user_uuid, u.username, u.domain_uuid").
		ColumnExpr("e.extension, e.extension_uuid").
		ColumnExpr("d.domain_name").
		Join("INNER JOIN v_domains d ON u.domain_uuid = d.domain_uuid").
		Join("INNER JOIN v_extension_users eu ON u.user_uuid = eu.user_uuid").
		Join("INNER JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Where("u.domain_uuid = ?", domainUuid).
		Where("u.user_uuid = ?", userUuid).Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) InsertUser(ctx context.Context, user *model.User) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(user)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert user failed")
	}
	return nil
}

func (repo *User) InsertUserTransaction(ctx context.Context, user *model.User, contact *model.VContact, groupUser *[]model.GroupUser, contactEmail *model.VContactEmail, roleGroup *model.RoleGroup) error {
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

func (repo *User) GetAllUserUuidOfGroup(ctx context.Context, domainUuid, groupUuid string) ([]string, error) {
	listUserUuid := make([]string, 0)

	query := repository.FusionSqlClient.GetDB().NewSelect().
		Table("v_group_users").
		Column("user_uuid").
		Where("domain_uuid = ?", domainUuid).
		Where("group_uuid = ?", groupUuid)

	err := query.Scan(ctx, &listUserUuid)
	if err == sql.ErrNoRows {
		return listUserUuid, nil
	} else if err != nil {
		return nil, err
	} else {
		return listUserUuid, nil
	}
}

func (repo *User) GetAgentInfoById(ctx context.Context, domainUuid, userUuid string) (*model.AgentInfo, error) {
	user := new(model.AgentInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		Relation("Campaigns").
		Relation("UserLive").
		Join("INNER JOIN v_contacts c ON c.contact_uuid = u.contact_uuid").
		ColumnExpr("u.*").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		Where("u.domain_uuid = ?", domainUuid).
		Where("u.user_uuid = ?", userUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) GetLiveUsersBackground(ctx context.Context) (*[]model.UserLive, error) {
	users := new([]model.UserLive)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(users).
		WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("status = ?", "available").Where("last_update_time <= ?", time.Now().Add(-10*time.Minute).Format("2006-01-02 15:04:05"))
		}).
		WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("status != ?", "available").Where("last_update_time <= ?", time.Now().Add(-1*time.Hour).Format("2006-01-02 15:04:05"))
		})
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

func (repo *User) GetUserByIdOrUsername(ctx context.Context, domainUuid string, value string) (*model.User, error) {
	user := new(model.User)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.*").
		Where("u.domain_uuid = ?", domainUuid).
		Where("(u.user_uuid::text = ? OR u.username = ?)", value, value).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) GetUserByUsername(ctx context.Context, domainUuid string, userName string) (*model.User, error) {
	user := new(model.User)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(user).
		ColumnExpr("u.*").
		Where("u.domain_uuid = ?", domainUuid).
		Where("u.username = ?", userName).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) GetContactById(ctx context.Context, domainUuid, contactUuid string) (*model.VContact, error) {
	contact := new(model.VContact)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(contact).
		Where("contact_uuid = ?", contactUuid).
		Where("domain_uuid = ?", domainUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return contact, nil
	}
}

func (repo *User) UpdateUserTransaction(ctx context.Context, user *model.User, contact *model.VContact, contactEmail *model.VContactEmail) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(user).WherePK().
			Column("username", "password", "salt", "level", "role_uuid", "unit_uuid", "enable_webrtc").
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

func (repo *User) GetUserPasswordInfo(ctx context.Context, domainUuid, userUuid string) (*model.User, error) {
	user := new(model.User)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(user).
		ColumnExpr("u.username, u.user_uuid, u.domain_uuid, u.api_key, u.user_enabled, u.password, u.salt, u.level as level").
		Where("u.domain_uuid = ?", domainUuid).
		Where("u.user_uuid = ?", userUuid).Limit(1)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) UpdateUserPassword(ctx context.Context, userUuid string, salt, passwordEncrypted string) error {
	query := repository.FusionSqlClient.GetDB().
		NewUpdate().
		Model((*model.User)(nil)).
		Set("password = ?", passwordEncrypted).
		Set("salt = ?", salt).
		Where("user_uuid = ?", userUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected != 1 {
		return errors.New("update user failed")
	}
	return nil
}

func (repo *User) GetUserByExtension(ctx context.Context, domainUuid, extension string) (*model.UserExtensionView, error) {
	user := new(model.UserExtensionView)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.user_uuid, u.username, u.domain_uuid").
		ColumnExpr("e.extension").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		Join("INNER JOIN v_extension_users eu ON u.user_uuid = eu.user_uuid").
		Join("INNER JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Join("INNER JOIN v_contacts c ON c.contact_uuid = u.contact_uuid").
		Where("u.domain_uuid = ?", domainUuid).
		Where("e.extension = ?", extension).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) GetLiveUsers(ctx context.Context, domainUuid string, filter model.MonitorFilter) (*[]model.UserLiveView, error) {
	users := new([]model.UserLiveView)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(users).
		Join("INNER JOIN v_users u ON u.user_uuid = ul.user_uuid").
		Join("INNER JOIN v_contacts c ON c.contact_uuid = u.contact_uuid").
		ColumnExpr("ul.*").
		ColumnExpr("u.username").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		Where("ul.domain_uuid = ?", domainUuid)
	if len(filter.CampaignUuids) > 0 {
		query = query.Where("ul.campaign_uuid IN (?)", bun.In(filter.CampaignUuids))
	}
	if len(filter.GroupUuids) > 0 {
		query = query.Where("ul.campaign_uuid IN (?)", bun.In(filter.GroupUuids))
	}
	if len(filter.UserUuids) > 0 {
		query = query.Where("ul.campaign_uuid IN (?)", bun.In(filter.UserUuids))
	}
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

func (repo *User) GetUserByUsernameAndDomainName(ctx context.Context, domainName string, username string) (*model.UserAuth, error) {
	user := new(model.UserAuth)
	err := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		ColumnExpr("u.username, u.user_uuid, u.domain_uuid, u.api_key, u.user_enabled, u.password, u.salt, u.level").
		ColumnExpr("d.domain_name").
		Join("inner join v_domains d on u.domain_uuid = d.domain_uuid").
		Where("u.username = ?", username).
		Where("d.domain_name = ?", domainName).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) UpdateLoginTimeUser(ctx context.Context, domainUuid string, userUuid string) error {
	resp, err := repository.FusionSqlClient.GetDB().NewUpdate().Model((*model.User)(nil)).
		Set("last_login_date = ?", time.Now()).
		Where("user_uuid = ?", userUuid).
		Exec(ctx)
	if err != nil {
		return err
	} else if afftected, _ := resp.RowsAffected(); afftected < 0 {
		return errors.New("update last_login_date fail")
	}
	return nil
}

func (repo *User) GetUserByAPIKey(ctx context.Context, apiKey string) (*model.UserAuth, error) {
	user := new(model.UserAuth)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(user).
		ColumnExpr("u.username, u.user_uuid, u.domain_uuid, d.domain_name, u.api_key, u.level").
		Join("JOIN v_domains d on d.domain_uuid = u.domain_uuid").
		Where("u.api_key = ?", apiKey).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) GetUserDataById(ctx context.Context, userUuid string) (*model.UserData, error) {
	user := new(model.UserData)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(user).
		ColumnExpr("u.user_uuid, u.domain_uuid").
		ColumnExpr("u.username, u.level").
		Column("u.role_uuid").
		Relation("Units").
		// Relation("Extensions").
		Where("u.user_uuid = ?", userUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) GetUsersInfoOfGroupUsers(ctx context.Context, domainUuid string, userUuid, userlevel string, groupUuids []string, level []string, isExtension bool) (*[]model.UserInfoData, error) {
	users := new([]model.UserInfoData)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(users).
		ColumnExpr("DISTINCT u.user_uuid, u.username").
		ColumnExpr("c.contact_name_given as first_name, c.contact_name_middle as middle_name, c.contact_name_family as last_name").
		ColumnExpr("ce.email_address as email").
		Join("LEFT JOIN v_group_users as gu ON gu.user_uuid = u.user_uuid").
		Join("LEFT JOIN v_contacts c ON c.contact_uuid = u.contact_uuid").
		Join("LEFT JOIN v_contact_emails ce ON ce.contact_uuid = c.contact_uuid ").
		Where("u.domain_uuid = ?", domainUuid)
	if len(groupUuids) > 0 {
		query = query.Where("gu.group_uuid IN (?)", bun.In(groupUuids))
	}
	if len(level) > 0 {
		query = query.Where("u.level IN (?)", bun.In(level))
	}
	if len(userlevel) > 0 {
		switch userlevel {
		case "leader":
			query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("u.level IN (?)", bun.In([]string{"agent", "user"})).
					WhereOr("u.level = ? AND u.user_uuid = ?", "leader", userUuid)
			})
		case "manager":
			query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
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
		query = query.ColumnExpr("tmpe.*").Join("LEFT JOIN lateral (?) tmpe ON true", subQuery)
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

func (repo *User) UpdateUserAndGroupTransaction(ctx context.Context, user *model.User, contact *model.VContact, groupUsers *[]model.GroupUser, contactEmail *model.VContactEmail) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(user).WherePK().
			Column("username", "password", "salt", "level", "role_uuid", "enable_webrtc", "user_enabled").Exec(ctx); err != nil {
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
		if groupUsers != nil && len(*groupUsers) >= 0 {
			if _, err := tx.NewDelete().Model((*model.GroupUser)(nil)).
				Where("user_uuid = ?", user.UserUuid).
				Where("domain_uuid = ?", user.DomainUuid).
				Exec(ctx); err != nil {
				return err
			}
			if len(*groupUsers) > 0 {
				if _, err := tx.NewInsert().Model(groupUsers).Exec(ctx); err != nil {
					return err
				}
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

func (repo *User) UpdateUserActive(ctx context.Context, domainUuid string, userUuid, extensionUuid string, enable bool) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		enableStr := "true"
		if !enable {
			enableStr = "false"
		}
		query := tx.NewUpdate().Model((*model.User)(nil)).Where("user_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Set("user_enabled = ?", enableStr)
		if _, err := query.Exec(ctx); err != nil {
			return err
		}
		if len(extensionUuid) > 0 {
			query = tx.NewUpdate().Model((*model.Extension)(nil)).Where("extension_uuid = ?", extensionUuid).
				Where("domain_uuid = ?", domainUuid).
				Set("enabled = ?", enableStr)
			if _, err := query.Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *User) DeleteUserTransaction(ctx context.Context, domainUuid, userUuid, contactUuid string) error {
	err := repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*model.GroupUser)(nil)).
			Where("user_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.CampaignUser)(nil)).
			Where("user_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.ExtensionUser)(nil)).
			Where("user_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.Contact)(nil)).
			Where("contact_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.VContactEmail)(nil)).
			Where("contact_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*model.User)(nil)).
			Where("user_uuid = ?", userUuid).
			Where("domain_uuid = ?", domainUuid).
			Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	return err
}
func (repo *User) GetUserLiveByExtension(ctx context.Context, domainUuid string, extension string) (*model.UserLive, error) {
	userLive := new(model.UserLive)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(userLive).
		Where("extension = ?", extension).
		Where("domain_uuid = ?", domainUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return userLive, nil
	}
}

func (repo *User) InsertUserCustomData(ctx context.Context, userCustomDatas ...model.UserCustomData) error {
	return repository.FusionSqlClient.GetDB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for _, userCustomData := range userCustomDatas {
			if existed, err := tx.NewSelect().Model((*model.UserCustomData)(nil)).
				Where("user_uuid = ?", userCustomData.UserUuid).
				Where("key = ?", userCustomData.Key).
				Exists(ctx); err != nil {
				return err
			} else if !existed {
				if _, err = tx.NewInsert().Model(&userCustomData).
					Exec(ctx); err != nil {
					return err
				}
			} else {
				if _, err := tx.NewUpdate().
					Model(&userCustomData).
					Column("value").
					Where("user_uuid = ?", userCustomData.UserUuid).
					Where("key = ?", userCustomData.Key).
					Exec(ctx); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (repo *User) GetUserCustomDatasByUserUuid(ctx context.Context, userUuid string) (*[]model.UserCustomData, error) {
	userCustomDatas := new([]model.UserCustomData)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(userCustomDatas).
		Where("user_uuid = ?", userUuid)
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return userCustomDatas, nil
}

func (repo *User) GetUserBasicInfoById(ctx context.Context, domainUuid, userId string) (*model.UserBasicInfo, error) {
	user := new(model.UserBasicInfo)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(user).
		Where("domain_uuid = ?", domainUuid).
		Where("user_uuid = ?", userId)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (repo *User) SelectUserLive(ctx context.Context) (*[]model.UserLive, error) {
	userLive := new([]model.UserLive)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(userLive).
		Order("login_time DESC")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return userLive, nil
	}
}
