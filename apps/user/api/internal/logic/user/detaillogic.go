package user

import (
	"context"

	"github.com/DullJZ/zeroim/apps/user/api/internal/svc"
	"github.com/DullJZ/zeroim/apps/user/api/internal/types"
	"github.com/DullJZ/zeroim/apps/user/rpc/user"
	"github.com/DullJZ/zeroim/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	detailResp, err := l.svcCtx.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}
	return &types.UserInfoResp{
		Info: types.User{
			Id:       detailResp.User.Id,
			Mobile:   detailResp.User.Phone,
			Nickname: detailResp.User.Nickname,
			Sex:      byte(detailResp.User.Sex),
			Avatar:   detailResp.User.Avatar,
		},
	}, nil
}
