package db

import (
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type CallCenter struct {
}

func NewCallCenter() repository.ICallCenter {
	repo := &CallCenter{}
	return repo
}

func (repo *CallCenter) GetCallCenterQueueById(ctx context.Context, callCenterQueueUuid string) (*model.CallCenterQueue, error) {
	callCenterQueue := new(model.CallCenterQueue)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(callCenterQueue).
		Where("call_center_queue_uuid = ?", callCenterQueueUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return callCenterQueue, nil
	}
}

func (repo *CallCenter) GetCallCenterQueueByExtension(ctx context.Context, domainUuid string, extension string) (*model.CallCenterQueue, error) {
	callCenter := model.CallCenterQueue{}
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Table("v_call_center_queues").
		Where("domain_uuid = ?", domainUuid).
		Where("queue_extension = ?", extension).
		Limit(1)

	err := query.Scan(ctx, &callCenter)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &callCenter, nil
	}
}

func (repo *CallCenter) GetCallCenterAgentById(ctx context.Context, callCenterAgentUuid string) (*model.CallCenterAgent, error) {
	callCenterAgent := new(model.CallCenterAgent)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(callCenterAgent).
		Where("call_center_agent_uuid = ?", callCenterAgentUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return callCenterAgent, nil
	}
}

func (repo *CallCenter) UpdateCallcenterAgentStatus(ctx context.Context, domainUuid, agentId, status string) error {
	callCenterAgent := new(model.CallCenterAgent)
	callCenterAgent.AgentStatus = status
	callCenterAgent.CallCenterAgentUuid = agentId

	query := repository.FusionSqlClient.GetDB().NewUpdate().Model(callCenterAgent).
		Where("call_center_agent_uuid = ?", agentId).
		Column("agent_status")

	res, err := query.Exec(ctx)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == -1 {
		return errors.New("update call_center_agent fail")
	}
	return nil
}

func (repo *CallCenter) GetCallCenterAgentByUserId(ctx context.Context, domainUuid, userUuid string) (*model.CallCenterAgent, error) {
	callCenterAgent := new(model.CallCenterAgent)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(callCenterAgent).
		Where("domain_uuid = ?", domainUuid).
		Where("user_uuid = ?", userUuid).
		Limit(1)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return callCenterAgent, nil
	}
}

func (repo *CallCenter) GetCallCenterTierOfAgent(ctx context.Context, domainUuid, agentUuid string) (*[]model.CallCenterTier, error) {
	callcenterTiers := new([]model.CallCenterTier)
	query := repository.FusionSqlClient.GetDB().NewSelect().Model(callcenterTiers).
		Where("domain_uuid = ?", domainUuid).
		Where("call_center_agent_uuid = ?", agentUuid)
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return callcenterTiers, nil
	}
}

func (repo *CallCenter) DeleteCallCenterTiersOfAgent(ctx context.Context, domainUuid, agentUuid string) error {
	query := repository.FusionSqlClient.GetDB().NewDelete().Model(&model.CallCenterTier{}).
		Where("domain_uuid = ?", domainUuid).
		Where("call_center_agent_uuid = ?", agentUuid)
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected == -1 {
		return errors.New("update call_center_agent fail")
	}
	return nil
}

func (repo *CallCenter) InsertCallCenterTier(ctx context.Context, callCenterTier *model.CallCenterTier) error {
	query := repository.FusionSqlClient.GetDB().NewInsert().Model(callCenterTier)
	resp, err := query.Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected != 1 {
		return errors.New("insert call_center_tier failed")
	}
	return nil
}

func (repo *CallCenter) GetCallCenterTiersReady(ctx context.Context, callCenterQueueUuid string) (*[]model.CallCenterTierWithExtension, error) {
	callcenterTiers := new([]model.CallCenterTierWithExtension)
	query := repository.FusionSqlClient.GetDB().NewSelect().
		Model(callcenterTiers).
		ColumnExpr("cct.*").
		ColumnExpr("e.extension").
		Join("INNER JOIN v_call_center_agents cca ON cct.call_center_agent_uuid = cca.call_center_agent_uuid").
		Join("INNER JOIN v_users u ON cca.user_uuid = u.user_uuid").
		Join("INNER JOIN user_live ul ON u.user_uuid = ul.user_uuid").
		Join("INNER JOIN v_extension_users eu ON u.user_uuid = eu.user_uuid").
		Join("INNER JOIN v_extensions e ON eu.extension_uuid = e.extension_uuid").
		Where("cct.call_center_queue_uuid = ?", callCenterQueueUuid).
		Where("ul.status = ?", "available")
	err := query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return callcenterTiers, nil
	}
}

func (repo *CallCenter) UpdateCallCenterQueueStrategy(ctx context.Context, callCenterQueueUuid, strategy string) error {
	res, err := repository.FusionSqlClient.GetDB().NewUpdate().Model((*model.CallCenterQueue)(nil)).
		Where("call_center_queue_uuid = ?", callCenterQueueUuid).
		Set("queue_strategy = ?", strategy).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected == -1 {
		return errors.New("update call_center_agent fail")
	} else {
		return nil
	}
}

func (repo *CallCenter) DeleteCallCenterQueue(ctx context.Context, callCenterQueueUuid string) error {
	query := repository.FusionSqlClient.GetDB().
		NewDelete().
		Model((*model.CallCenterQueue)(nil)).
		Where("call_center_queue_uuid = ?", callCenterQueueUuid)
	if _, err := query.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (repo *CallCenter) GetCallCenterQueues(ctx context.Context, domainUuid string, filter model.CallCenterQueueFilter, limit, offset int) ([]model.CallCenterQueue, int, error) {
	callCenterQueues := make([]model.CallCenterQueue, 0)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(&callCenterQueues)
	if len(filter.CallCenterQueueUuids) > 0 {
		query.Where("call_center_queue_uuid IN (?)", bun.In(filter.CallCenterQueueUuids))
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	return callCenterQueues, total, nil
}

func (repo *CallCenter) GetCallCenterAgentsOfCallCenterQueues(ctx context.Context, callCenterQueueUuid ...string) (*[]model.CallCenterAgent, error) {
	callCenterAgents := new([]model.CallCenterAgent)
	query := repository.FusionSqlClient.GetDB().
		NewSelect().
		Model(callCenterAgents).
		Join("INNER JOIN v_call_center_tiers cct ON cct.call_center_agent_uuid = cca.call_center_agent_uuid").
		Join("INNER JOIN v_call_center_queues ccq ON cct.call_center_queue_uuid = ccq.call_center_queue_uuid")
	if len(callCenterQueueUuid) > 0 {
		query.Where("ccq.call_center_queue_uuid IN (?)", bun.In(callCenterQueueUuid))
	}
	err := query.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return callCenterAgents, nil
}
