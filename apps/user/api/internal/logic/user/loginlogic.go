package user

import (
	"context"

	"github.com/DullJZ/zeroim/apps/user/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/api/internal/types"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResq, err := l.svcCtx.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token:  loginResq.Token,
		Expire: loginResq.Expire,
	}, nil
}
