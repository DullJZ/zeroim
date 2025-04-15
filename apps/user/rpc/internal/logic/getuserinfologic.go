package logic

import (
	"context"
	"fmt"

	"github.com/DullJZ/zeroim/apps/user/model"
	"github.com/DullJZ/zeroim/apps/user/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

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
			return nil, fmt.Errorf("未找到该用户ID")
		}
		return nil, err
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
