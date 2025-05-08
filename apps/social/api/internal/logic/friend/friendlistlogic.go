package friend

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/api/internal/types"
	"github.com/DullJZ/zeroim/apps/social/rpc/socialclient"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"
	"github.com/DullJZ/zeroim/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	friends, err := l.svcCtx.SocialRpc.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}
	friendUids := make([]string, 0, len(friends.List))
	for _, friend := range friends.List {
		friendUids = append(friendUids, friend.FriendUid)
	}

	// 获取好友信息
	users, err := l.svcCtx.UserRpc.FindUser(l.ctx, &user.FindUserReq{
		Ids: friendUids,
	})
	if err != nil {
		return nil, err
	}

	// 组装好友信息
	var friendList []*types.Friends
	for _, friend := range friends.List {
		for _, user := range users.User {
			if friend.FriendUid == user.Id {
				friendList = append(friendList, &types.Friends{
					Id:        int32(friend.Id),
					FriendUid: friend.FriendUid,
					Nickname:  user.Nickname,
					Avatar:    user.Avatar,
					Remark:    friend.Remark,
				})
			}
		}
	}

	return &types.FriendListResp{
		List: friendList,
	}, nil
}
