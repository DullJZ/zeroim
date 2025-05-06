package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/model"
	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUsersLogic {
	return &GroupUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	groupMem, err := l.svcCtx.GroupMembersModel.FindByGroupId(l.ctx, in.GroupId)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	var list []*social.GroupMember
	for _, mem := range groupMem {
		list = append(list, &social.GroupMember{
			Id:          int32(mem.Id),
			GroupId:     mem.GroupId,
			UserId:      mem.UserId,
			RoleLevel:   int32(mem.RoleLevel),
			JoinTime:    mem.JoinTime.Time.Unix(),
			JoinSource:  int32(mem.JoinSource.Int64),
			InviterUid:  mem.InviterUid.String,
			OperatorUid: mem.OperatorUid.String,
		})
	}
	return &social.GroupUsersResp{
		List: list,
	}, nil
}
