package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

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