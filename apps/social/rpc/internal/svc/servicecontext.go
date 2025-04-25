package svc

import (
	"github.com/DullJZ/zeroim/apps/social/model"
	"github.com/DullJZ/zeroim/apps/social/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	model.FriendsModel
	model.FriendRequestsModel
	model.GroupsModel
	model.GroupMembersModel
	model.GroupRequestsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.Datasource)
	return &ServiceContext{
		Config: c,

		FriendsModel:        model.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: model.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         model.NewGroupsModel(sqlConn, c.Cache),
		GroupMembersModel:   model.NewGroupMembersModel(sqlConn, c.Cache),
		GroupRequestsModel:  model.NewGroupRequestsModel(sqlConn, c.Cache),
	}
}
