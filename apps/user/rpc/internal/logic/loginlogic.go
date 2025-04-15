package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/DullJZ/zeroim/apps/user/model"
	"github.com/DullJZ/zeroim/apps/user/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"
	"github.com/DullJZ/zeroim/pkg/ctxdata"
	"github.com/DullJZ/zeroim/pkg/encrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 查找手机号是否已经注册
	u, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("用户手机号未注册")
		}
		return nil, err
	}
	// 校验密码
	if !encrypt.ValidatePasswordHash(in.Password, u.Password.String) {
		return nil, fmt.Errorf("密码错误")
	}
	// 生成token
	token, err := ctxdata.GetJwtToken(
		l.svcCtx.Config.Jwt.AccessSecret,
		time.Now().Unix(),
		l.svcCtx.Config.Jwt.AccessExpire,
		u.Id,
	)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		Token:  token,
		Expire: time.Now().Unix() + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
