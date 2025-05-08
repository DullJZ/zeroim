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

type GroupPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInLogic) GroupPutIn(in *social.GroupPutInReq) (*social.GroupPutInResp, error) {
	// 检查是否已经在群内
	_, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndUserId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member err %v req %v", err, in)
	}
	// 检查是否已经申请过并且未处理
	_, err = l.svcCtx.GroupRequestsModel.FindByUserIdAndGroupIdAndState(l.ctx, in.ReqId, in.GroupId, int(constants.GroupRequestStatusWait))
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group request err %v req %v", err, in)
	}
	// 创建申请记录
	_, err = l.svcCtx.GroupRequestsModel.Insert(l.ctx, &model.GroupRequests{
		ReqId:   in.ReqId,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group request err %v req %v", err, in)
	}

	return &social.GroupPutInResp{}, nil
}
