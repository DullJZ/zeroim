package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (m *customFriendRequestsModel) FindByUserId(ctx context.Context, uid string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", friendRequestsRows, m.table)
	var resp []*FriendRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customFriendRequestsModel) FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ?", friendRequestsRows, m.table)
	var resp FriendRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, rid, uid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customFriendRequestsModel) Update(ctx context.Context, data *FriendRequests) error {
	friendRequestIdKey := fmt.Sprintf("%s%v", cacheFriendRequestsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRequestsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt, data.Id)
	}, friendRequestIdKey)
	return err
}
