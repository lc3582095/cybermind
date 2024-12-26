package model

// 错误码定义
const (
	Success       = 0    // 成功
	ParamError    = 1001 // 参数错误
	Unauthorized  = 1002 // 未授权
	Forbidden     = 1003 // 禁止访问
	NotFound      = 1004 // 资源不存在
	SystemError   = 1005 // 系统错误
	AdminNotExist = 2001 // 管理员不存在
	WrongPassword = 2002 // 密码错误
	AdminDisabled = 2003 // 账号已禁用
) 