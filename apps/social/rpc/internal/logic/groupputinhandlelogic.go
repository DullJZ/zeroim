package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// 检查是否已经处理
	r, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, err
	}
	if r.HandleResult.Int64 != int64(constants.GroupRequestStatusWait) {
		return nil, err
	}
	// 修改申请结果
	r.HandleResult.Int64 = int64(in.HandleResult)
	err = l.svcCtx.GroupRequestsModel.Update(l.ctx, r)
	if err != nil {
		return nil, err
	}
	return &social.GroupPutInHandleResp{}, nil
}
