package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/social/model"
	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/constants"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrFriendReqBeforePass = errors.New("该好友申请已通过")
var ErrFriendReqBeforeRefuse = errors.New("该好友申请已拒绝")

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 获取好友申请记录
	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend request err %v req %v", err, in)
	}

	// 验证是否有处理
	switch constants.HandlerResult(friendReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}
	friendReq.HandleResult.Int64 = int64(in.HandleResult)

	// 修改申请结果
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, friendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v req %v", err, friendReq)
		}
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}
		friends := []*model.Friends{
			{
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			},
			{
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}
		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friends err %v req %v", err, friends)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update friend request err %v req %v", err, in)
	}

	return &social.FriendPutInHandleResp{}, nil
}
