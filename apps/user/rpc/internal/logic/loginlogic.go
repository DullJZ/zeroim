package logic

import (
	"context"
	"time"

	"github.com/DullJZ/zeroim/apps/user/model"
	"github.com/DullJZ/zeroim/apps/user/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"
	"github.com/DullJZ/zeroim/pkg/ctxdata"
	"github.com/DullJZ/zeroim/pkg/encrypt"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号未注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
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
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}
	// 校验密码
	if !encrypt.ValidatePasswordHash(in.Password, u.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}
	// 生成token
	token, err := ctxdata.GetJwtToken(
		l.svcCtx.Config.Jwt.AccessSecret,
		time.Now().Unix(),
		l.svcCtx.Config.Jwt.AccessExpire,
		u.Id,
	)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "get jwt token err %v", err)
	}

	return &user.LoginResp{
		Token:  token,
		Expire: time.Now().Unix() + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
