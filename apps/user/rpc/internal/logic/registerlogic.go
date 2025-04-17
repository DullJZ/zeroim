package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/DullJZ/zeroim/apps/user/model"
	"github.com/DullJZ/zeroim/apps/user/rpc/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"
	"github.com/DullJZ/zeroim/pkg/ctxdata"
	"github.com/DullJZ/zeroim/pkg/encrypt"
	"github.com/DullJZ/zeroim/pkg/wuid"
	"github.com/DullJZ/zeroim/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	ErrPhoneAlreadyRegistered = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号已注册")
	ErrPasswordTooShort       = xerr.New(xerr.SERVER_COMMON_ERROR, "密码小于六位")
)

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 1. 检查手机号是否已注册
	userEntity, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if userEntity != nil {
		return nil, ErrPhoneAlreadyRegistered
	}

	userEntity = &model.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex:      sql.NullInt64{Int64: int64(in.Sex), Valid: true},
	}

	// 2. 检查密码位数
	if len(in.Password) < 6 {
		return nil, ErrPasswordTooShort
	}

	genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
	if err != nil {
		return nil, err
	}
	userEntity.Password = sql.NullString{
		String: string(genPassword),
		Valid:  true,
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, userEntity)

	if err != nil {
		return nil, err
	}

	// 3. 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil

}
