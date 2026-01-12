// Package user 用户领域模块（C端用户）
package member

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-member/model"
)

// Repository 用户仓储接口
type Repository interface {
	// 基础 CRUD
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error

	// 分页查询
	Paginate(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error)

	// 业务特定查询
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByName(ctx context.Context, name string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// LoginLogRepository 登录日志仓储接口
type LoginLogRepository interface {
	Create(ctx context.Context, log *model.UserLoginLog) error
	FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginLog, int64, error)
	// 分页查询（支持日期范围过滤）
	Paginate(ctx context.Context, page, pageSize int, userID uint, startDate, endDate string) ([]model.UserLoginLog, int64, error)
}

