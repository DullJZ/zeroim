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
	// 检查是否已经处理
	r, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group request err %v req %v", err, in)
	}
	if r.HandleResult.Int64 != int64(constants.GroupRequestStatusWait) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "apply already handled, err %v req %v", err, in)
	}

	// 验证处理人是否有权限处理申请（群主或管理员）
	handler, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndUserId(l.ctx, r.GroupId, in.HandleUid)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member err %v req %v", err, in)
	}
	if handler == nil || (handler.RoleLevel != int64(constants.GroupRoleOwner) && 
		handler.RoleLevel != int64(constants.GroupRoleAdmin)) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "no permission to handle apply, err %v req %v", err, in)
	}

	// 如果是通过申请，检查用户是否已经是群成员
	if in.HandleResult == int32(constants.GroupRequestStatusPass) {
		member, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndUserId(l.ctx, r.GroupId, r.ReqId)
		if err != nil && err != model.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewDBErr(), "check user in group err %v req %v", err, in)
		}
		if member != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "user already in group, err %v req %v", err, in)
		}
	}

	// 修改申请结果
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
			return errors.Wrapf(xerr.NewDBErr(), "update group request err %v req %v", err, in)
		}
		
		if in.HandleResult == int32(constants.GroupRequestStatusPass) {
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
			if err != nil {
				return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
			}
		}
		return nil
	})
	
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "group put in handle err %v req %v", err, in)
	}
	
	return &social.GroupPutInHandleResp{
		GroupId: in.GroupId,
	}, nil
}
