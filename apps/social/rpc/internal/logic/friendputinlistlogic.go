package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend request err %v req %v", err, in)
	}
	var list []*social.FriendRequest
	for _, friendReq := range friendReqs {
		list = append(list, &social.FriendRequest{
			Id:           int32(friendReq.Id),
			UserId:       friendReq.UserId,
			ReqUid:       friendReq.ReqUid,
			ReqMsg:       friendReq.ReqMsg.String,
			ReqTime:      friendReq.ReqTime.Unix(),
			HandleResult: int32(friendReq.HandleResult.Int64),
		})
	}

	return &social.FriendPutInListResp{
		List: list,
	}, nil
}
