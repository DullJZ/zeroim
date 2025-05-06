package model

import (
	"context"
	"fmt"
)

func (m *customGroupMembersModel) FindByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ?", groupMembersRows, m.table)
	var resp []*GroupMembers
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customGroupMembersModel) FindByGroupIdAndUserId(ctx context.Context, groupId, userId string) (*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `user_id` = ?", groupMembersRows, m.table)
	var resp GroupMembers
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId, userId)
	switch err {
	case nil:
		return &resp, nil
	case ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
