package logic

import (
	"context"

	"github.com/DullJZ/zeroim/apps/user/model"
	"github.com/DullJZ/zeroim/apps/user/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrFindUserById = errors.New("未找到该用户ID")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.WithStack(ErrFindUserById)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "get user info err %v", err)
	}

	return &user.GetUserInfoResp{
		User: &user.UserEntity{
			Id:       u.Id,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
			Phone:    u.Phone,
			Status:   int32(u.Status.Int64),
			Sex:      int32(u.Sex.Int64),
		},
	}, nil
}
