package svc

import (
	"github.com/DullJZ/zeroim/apps/social/api/internal/config"
	"github.com/DullJZ/zeroim/apps/social/rpc/socialclient"
	"github.com/DullJZ/zeroim/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc userclient.User
	SocialRpc socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		SocialRpc: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
