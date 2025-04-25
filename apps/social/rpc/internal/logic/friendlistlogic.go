package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// 获取好友列表
	friends, err := l.svcCtx.FriendsModel.ListByUserid(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend list err %v req %v", err, in)
	}

	var resp []*social.Friend
	for _, friend := range friends {
		resp = append(resp, &social.Friend{
			Id:        int32(friend.Id),
			UserId:    friend.UserId,
			Remark:    friend.Remark.String,
			AddSource: int32(friend.AddSource.Int64),
			FriendUid: friend.FriendUid,
		})
	}

	return &social.FriendListResp{
		List: resp,
	}, nil
}
