package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (m *customFriendsModel) FindByUidAndFid(ctx context.Context, uid, fid string) (*Friends, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `friend_uid` = ?", friendsRows, m.table)
	var resp Friends
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid, fid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFriendsModel) Inserts(ctx context.Context, session sqlx.Session, data ...*Friends) (sql.Result, error) {
	var sql strings.Builder
	var args []interface{}
	if len(data) == 0 {
		return nil, nil
	}
	sql.WriteString(fmt.Sprintf("insert into %s (%s) values ", m.table, friendsRowsExpectAutoSet))

	for i, v := range data {
		sql.WriteString("(?, ?, ?, ?)")
		args = append(args, v.UserId, v.FriendUid, v.Remark, v.AddSource)
		if i == len(data)-1 {
			break
		}

		sql.WriteString(",")
	}
	return session.ExecCtx(ctx, sql.String(), args...)
}

func (m *customFriendsModel) ListByUserid(ctx context.Context, userId string) ([]*Friends, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", friendsRows, m.table)
	var resp []*Friends
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
