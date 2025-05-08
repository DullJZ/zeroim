package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DullJZ/zeroim/apps/social/model"
	"github.com/DullJZ/zeroim/apps/social/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/social/rpc/social"
	"github.com/DullJZ/zeroim/pkg/constants"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	logc.Alert(l.ctx, "入口点")
	// 检查是否已经处理
	r, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil || err == model.ErrNotFound {
		return nil, err
	}
	if r.HandleResult.Int64 != int64(constants.GroupRequestStatusWait) {
		return nil, fmt.Errorf("申请已处理")
	}
	// 修改申请结果
	logc.Alert(l.ctx, "修改申请结果")
	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新申请结果
		err := l.svcCtx.GroupRequestsModel.Update(ctx, &model.GroupRequests{
			Id:            r.Id,
			ReqId:         r.ReqId,
			GroupId:       r.GroupId,
			ReqMsg:        r.ReqMsg,
			ReqTime:       r.ReqTime,
			JoinSource:    r.JoinSource,
			InviterUserId: r.InviterUserId,
			HandleUserId: sql.NullString{
				String: in.HandleUid,
				Valid:  true,
			},
			HandleTime: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
			HandleResult: sql.NullInt64{
				Int64: int64(in.HandleResult),
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		
		// 加入群组成员
		_, err = l.svcCtx.GroupMembersModel.Insert(ctx, &model.GroupMembers{
			GroupId:   r.GroupId,
			UserId:    r.ReqId,
			RoleLevel: int64(constants.GroupRoleMember),
			JoinTime: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
			JoinSource: sql.NullInt64{
				Int64: r.JoinSource.Int64,
				Valid: true,
			},
			InviterUid: sql.NullString{
				String: r.InviterUserId.String,
				Valid:  true,
			},
			OperatorUid: sql.NullString{
				String: in.HandleUid,
				Valid:  true,
			},
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return &social.GroupPutInHandleResp{
		GroupId: in.GroupId,
	}, nil
}
