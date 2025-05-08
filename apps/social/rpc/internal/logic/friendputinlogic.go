package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/DullJZ/zeroim/apps/social/model"
	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/constants"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 好友业务: 请求好友、通过或拒绝申请、好友列表
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 申请人和目标是否已经是好友
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend by id and fid err %v req %v", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, nil
	}
	// 是否有过申请
	friendReq, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend request by id and fid err %v req %v", err, in)
	}
	if friendReq != nil {
		return &social.FriendPutInResp{}, nil
	}

	// 创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &model.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friend request err %v req %v", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
