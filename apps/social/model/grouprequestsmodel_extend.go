package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (m *customGroupRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customGroupRequestsModel) FindByUserId(ctx context.Context, uid string) ([]*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", groupRequestsRows, m.table)
	var resp []*GroupRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customGroupRequestsModel) FindByGroupId(ctx context.Context, groupId string) ([]*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ?", groupRequestsRows, m.table)
	var resp []*GroupRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customGroupRequestsModel) FindByUserIdAndGroupIdAndState(ctx context.Context, groupId, userId string, state int) (*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `req_id` = ? and `handle_result` = %d ", groupRequestsRows, m.table, state)
	var groupReq GroupRequests
	err := m.QueryRowNoCacheCtx(ctx, &groupReq, query, groupId, userId)
	switch err {
	case nil:
		return &groupReq, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

