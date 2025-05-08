package friend

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/api/internal/types"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	_, err = l.svcCtx.SocialRpc.FriendPutIn(l.ctx, &social.FriendPutInReq{
		UserId:  req.UserId,
		ReqUid:  uid,
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})
	return &types.FriendPutInResp{}, err
}
