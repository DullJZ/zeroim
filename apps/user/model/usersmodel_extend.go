package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	cacheUsersPhonePrefix = "cache:users:phone:"
	cacheUsersNamePrefix  = "cache:users:name:"
	cacheUsersIdsPrefix   = "cache:users:ids:"
)

func (m *defaultUsersModel) FindOneByPhone(ctx context.Context, phone string) (*Users, error) {
	usersPhoneKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, phone)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, usersPhoneKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, phone)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
	usersPhoneKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, data.Phone)

	result, err := m.defaultUsersModel.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	// 手动删除Phone相关缓存，防止重复注册
	err = m.defaultUsersModel.CachedConn.DelCache(usersPhoneKey)
	if err != nil {
		return result, nil
	}

	return result, nil
}

func (m *defaultUsersModel) ListByName(ctx context.Context, name string) ([]*Users, error) {
	usersNameKey := fmt.Sprintf("%s%v", cacheUsersNamePrefix, name)
	var resp []*Users

	err := m.QueryRowCtx(ctx, &resp, usersNameKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `nickname` like ? ", usersRows, m.table)
		return conn.QueryRowsCtx(ctx, v, query, fmt.Sprint("%", name, "%"))
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultUsersModel) ListByIds(ctx context.Context, ids []string) ([]*Users, error) {
	// 对于多个ID，我们使用一个组合的key
	usersIdsKey := fmt.Sprintf("%s%v", cacheUsersIdsPrefix, strings.Join(ids, "_"))
	var resp []*Users

	err := m.QueryRowCtx(ctx, &resp, usersIdsKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` in ('%s') ", usersRows, m.table, strings.Join(ids, "','"))
		return conn.QueryRowsCtx(ctx, v, query)
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}
