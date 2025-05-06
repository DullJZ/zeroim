package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (m *customGroupsModel) FindByUserId(ctx context.Context, uid string) ([]*Groups, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", groupMembersRows, `group_members`)
	var groupMembers []*GroupMembers
	err := m.QueryRowsNoCacheCtx(ctx, &groupMembers, query, uid)
	if err != nil {
		return nil, err
	}
	var resp []*Groups
	for _, groupMember := range groupMembers {
		group, err := m.FindOne(ctx, groupMember.GroupId)
		if err != nil {
			return nil, err
		}
		resp = append(resp, group)
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customGroupsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}
