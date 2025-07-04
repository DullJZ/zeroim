package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupRequestsModel = (*customGroupRequestsModel)(nil)

type (
	// GroupRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupRequestsModel.
	GroupRequestsModel interface {
		groupRequestsModel

		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		FindByUserId(ctx context.Context, uid string) ([]*GroupRequests, error)
		FindByGroupId(ctx context.Context, groupId string) ([]*GroupRequests, error)
		FindByUserIdAndGroupIdAndState(ctx context.Context, uid string, groupId string, state int) (*GroupRequests, error)
	}

	customGroupRequestsModel struct {
		*defaultGroupRequestsModel
	}
)

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}
