package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	cacheUsersPhonePrefix = "cache:users:phone:"
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

func (m *defaultUsersModel) ListByName(ctx context.Context, name string) ([]*Users, error) {
	var resp []*Users

	query := fmt.Sprintf("select %s from %s where `nickname` like ? ", usersRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, fmt.Sprint("%", name, "%"))

	if err != nil {
		return nil, err
	}
	return resp, nil

}

func (m *defaultUsersModel) ListByIds(ctx context.Context, ids []string) ([]*Users, error) {
	var resp []*Users

	query := fmt.Sprintf("select %s from %s where `id` in ('%s') ", usersRows, m.table, strings.Join(ids, "','"))
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
