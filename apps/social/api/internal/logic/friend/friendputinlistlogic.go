package friend

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/api/internal/types"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请列表
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	list, err := l.svcCtx.SocialRpc.FriendPutInList(l.ctx, &social.FriendPutInListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}
	respList := make([]*types.FriendRequests, 0, len(list.List))
	for _, item := range list.List {
		respList = append(respList, &types.FriendRequests{
			Id:           int64(item.Id),
			UserId:       item.UserId,
			ReqUid:       item.ReqUid,
			ReqMsg:       item.ReqMsg,
			ReqTime:      item.ReqTime,
			HandleResult: int(item.HandleResult),
		})
	}

	return &types.FriendPutInListResp{
		List: respList,
	}, nil
}
