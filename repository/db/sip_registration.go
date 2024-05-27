package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
)

type SipRegistration struct {
}

func NewSipRegistration() repository.ISipRegistration {
	repo := &SipRegistration{}
	return repo
}

func (repo *SipRegistration) GetSipRegistrationsNotCallOfDomain(ctx context.Context, domainName string) (*[]model.SipRegistration, error) {
	sipRegistrations := new([]model.SipRegistration)
	query := repository.FreeswitchSqlClient.GetDB().NewSelect().
		Model(sipRegistrations).
		Join("LEFT JOIN sip_dialogs sd ON sd.contact_user = sr.sip_user AND (sr.sip_host = sd.sip_to_host OR sr.sip_host = sd.sip_from_host)").
		Where("sd.call_id IS NULL").
		Where("sr.sip_host = ?", domainName)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return sipRegistrations, nil
	} else if err != nil {
		return nil, err
	} else {
		return sipRegistrations, nil
	}
}

func (repo *SipRegistration) GetSipRegistrationOfExtension(ctx context.Context, domainName, extension string) (any, error) {
	var sipRegistration model.SipRegistration
	query := repository.FreeswitchSqlClient.GetDB().NewSelect().
		Table("sip_registrations").
		Column("sip_user", "sip_host", "status", "sip_realm", "mwi_user", "mwi_host").
		Where("sip_user = ?", extension).
		Where("sip_host = ?", domainName).
		Limit(1)
	err := query.
		Scan(ctx, &sipRegistration)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return sipRegistration, nil
}

func (repo *SipRegistration) GetSipRegistrations(ctx context.Context, domainName string) (*[]model.SipRegistrationMonitorView, error) {
	sipRegistrations := new([]model.SipRegistrationMonitorView)
	query := repository.FreeswitchSqlClient.GetDB().NewSelect().
		TableExpr("sip_registrations sr").
		ColumnExpr("sr.sip_user, sr.server_host, sr.network_ip, sr.user_agent, sd.uuid as call_id").
		ColumnExpr("b_ch.initial_cid_num as caller, a_ch.dest as callee, c.call_created as call_time").
		ColumnExpr(`
			CASE
				WHEN b_ch.callstate = 'EARLY' THEN 'RINGING'
				WHEN b_ch.callstate = 'ACTIVE' THEN 'ONCALL'
				ELSE b_ch.callstate
			END AS call_state
		`).
		ColumnExpr(`
			CASE
				WHEN a_ch.direction = 'local' THEN 'local'
				WHEN sd.uuid = c.caller_uuid THEN 'outbound'
				WHEN sd.uuid = c.callee_uuid THEN 'inbound'
				ELSE ''
			END AS direction
		`).
		ColumnExpr(`
			CASE
				WHEN sd.uuid = c.caller_uuid THEN b_ch.initial_cid_num
				WHEN sd.uuid = c.callee_uuid THEN a_ch.initial_dest
				ELSE ''
			END AS hotline
		`).
		Join("LEFT JOIN sip_dialogs sd ON sd.contact_user = sr.sip_user AND (sd.sip_to_host = sr.sip_host OR sd.sip_from_host = sr.sip_host)").
		Join("LEFT JOIN calls c ON sd.uuid = c.caller_uuid OR sd.uuid = c.callee_uuid").
		Join("LEFT JOIN channels a_ch ON a_ch.uuid = c.caller_uuid").
		Join("LEFT JOIN channels b_ch ON b_ch.uuid = c.callee_uuid").
		Where("sr.sip_host = ?", domainName).
		Order("call_state")
	err := query.
		Scan(ctx, sipRegistrations)
	if err == sql.ErrNoRows {
		return sipRegistrations, nil
	} else if err != nil {
		return nil, err
	}
	return sipRegistrations, nil
}
