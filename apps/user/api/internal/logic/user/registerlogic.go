package user

import (
	"context"

	"github.com/DullJZ/zeroim/apps/user/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/api/internal/types"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResp, err := l.svcCtx.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, err
	}
	return &types.RegisterResp{
		Token:  registerResp.Token,
		Expire: registerResp.Expire,
	}, nil
}
