package constants

// 处理结果 1.未处理 2.通过 3.拒绝
type HandlerResult int

const (
	NoHandlerResult     HandlerResult = iota + 1 // 未处理
	PassHandlerResult                            // 通过
	RefuseHandlerResult                          // 拒绝
	CancelHandlerResult                          // 取消
)


// 群聊用户角色
type GroupRole int

const (
	GroupRoleOwner GroupRole = iota + 1 // 群主
	GroupRoleAdmin                      // 管理员
	GroupRoleMember                     // 普通成员
)

// 群聊申请状态
type GroupRequestStatus int

const (
	GroupRequestStatusWait GroupRequestStatus = iota + 1 // 等待处理
	GroupRequestStatusPass                              // 通过
	GroupRequestStatusRefuse                            // 拒绝
	GroupRequestStatusCancel                            // 取消
)
