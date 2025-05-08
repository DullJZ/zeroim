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
	"github.com/DullJZ/zeroim/pkg/wuid"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群业务: 创建群、修改群、群公告、申请群、用户群列表、群成员、申请群、群退出..
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	var groupId string
	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 创建群组
		groupId = wuid.GenUid(l.svcCtx.Config.Mysql.Datasource)
		_, err := l.svcCtx.GroupsModel.Insert(ctx, &model.Groups{
			Id:         groupId,
			Name:       in.Name,
			Icon:       in.Icon,
			CreatorUid: in.CreatorUid,
			Status: sql.NullInt64{
				Int64: int64(in.Status),
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		// 2. 添加创建者为群主
		_, err = l.svcCtx.GroupMembersModel.Insert(ctx, &model.GroupMembers{
			GroupId:   groupId,
			UserId:    in.CreatorUid,
			RoleLevel: int64(constants.GroupRoleOwner),
			JoinTime: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
		})
		return err
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "create group err %v req %v", err, in)
	}

	return &social.GroupCreateResp{
		GroupId: groupId,
	}, nil
}
