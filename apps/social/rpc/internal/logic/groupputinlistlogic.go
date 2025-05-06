package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/model"
	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInListLogic) GroupPutInList(in *social.GroupPutInListReq) (*social.GroupPutInListResp, error) {
	groupReqs, err := l.svcCtx.GroupRequestsModel.FindByGroupId(l.ctx, in.GroupId)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	var list []*social.GroupRequest
	for _, groupReq := range groupReqs {
		list = append(list, &social.GroupRequest{
			Id:         int32(groupReq.Id),
			GroupId:    groupReq.GroupId,
			ReqId:      groupReq.ReqId,
			ReqMsg:     groupReq.ReqMsg.String,
			ReqTime:    groupReq.ReqTime.Time.Unix(),
			JoinSource: int32(groupReq.JoinSource.Int64),
			InviterUid: groupReq.InviterUserId.String,
			HandleUid:  groupReq.HandleUserId.String,
			HandleResult: int32(groupReq.HandleResult.Int64),
		})
	}

	return &social.GroupPutInListResp{
		List: list,
	}, nil
}
