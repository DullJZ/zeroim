package friend

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/api/internal/types"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请处理
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	_, err = l.svcCtx.SocialRpc.FriendPutInHandle(l.ctx, &social.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		UserId:       ctxdata.GetUid(l.ctx),
		HandleResult: req.HandleResult,
	})
	return &types.FriendPutInHandleResp{}, err
}
