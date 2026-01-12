package member

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

// 错误码模块
const ModuleMember = 25

// 领域错误定义
var (
	// ErrDatabaseError 数据库错误
	ErrDatabaseError = errcode.Register(errcode.New(
		ModuleMember, 1001,
		"member",
		"error.member.database_error",
		"数据库操作失败",
		http.StatusInternalServerError,
	))

	// ErrNotFound 会员不存在
	ErrNotFound = errcode.Register(errcode.New(
		ModuleMember, 1002,
		"member",
		"error.member.not_found",
		"会员不存在",
		http.StatusNotFound,
	))

	// ErrEmailExists 邮箱已存在
	ErrEmailExists = errcode.Register(errcode.New(
		ModuleMember, 1003,
		"member",
		"error.member.email_exists",
		"邮箱已存在",
		http.StatusConflict,
	))

	// ErrInvalidCredentials 无效凭证
	ErrInvalidCredentials = errcode.Register(errcode.New(
		ModuleMember, 1004,
		"member",
		"error.member.invalid_credentials",
		"邮箱或密码错误",
		http.StatusUnauthorized,
	))

	// ErrAccountDisabled 账户已禁用
	ErrAccountDisabled = errcode.Register(errcode.New(
		ModuleMember, 1005,
		"member",
		"error.member.account_disabled",
		"账户已禁用",
		http.StatusForbidden,
	))

	// ErrInternalError 内部错误
	ErrInternalError = errcode.Register(errcode.New(
		ModuleMember, 1006,
		"member",
		"error.member.internal_error",
		"内部错误",
		http.StatusInternalServerError,
	))
)
