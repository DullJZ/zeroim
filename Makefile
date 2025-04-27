user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/make/user-api.mk release-test

social-rpc-dev:
	@make -f deploy/make/social-rpc.mk release-test

social-api-dev:
	@make -f deploy/make/social-api.mk release-test

