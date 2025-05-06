package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	groups, err := l.svcCtx.GroupsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	var list []*social.Group
	for _, group := range groups {
		list = append(list, &social.Group{
			Id: group.Id,
			Name: group.Name,
			Icon: group.Icon,
			Status: int32(group.Status.Int64),
			CreatorUid: group.CreatorUid,
		})
	}
	return &social.GroupListResp{
		List: list,
	}, nil
}
