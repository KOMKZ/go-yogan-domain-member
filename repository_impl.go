package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/KOMKZ/go-yogan-domain-member/model"
	"gorm.io/gorm"
)

// GORMRepository 用户仓储 GORM 实现
type GORMRepository struct {
	db *gorm.DB
}

// 编译时检查
var _ Repository = (*GORMRepository)(nil)

// NewGORMRepository 创建用户仓储
func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{db: db}
}

// Create 创建用户
func (r *GORMRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据ID查询用户
func (r *GORMRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindAll 获取所有用户
func (r *GORMRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

// Update 更新用户
func (r *GORMRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户（软删除）
func (r *GORMRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// Paginate 分页查询
func (r *GORMRepository) Paginate(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByEmail 根据邮箱查询用户
func (r *GORMRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("根据邮箱查询用户失败: %w", err)
	}
	return &user, nil
}

// FindByName 根据用户名查询用户
func (r *GORMRepository) FindByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("根据用户名查询用户失败: %w", err)
	}
	return &user, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *GORMRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// LoginLogGORMRepository 登录日志仓储 GORM 实现
type LoginLogGORMRepository struct {
	db *gorm.DB
}

// 编译时检查
var _ LoginLogRepository = (*LoginLogGORMRepository)(nil)

// NewLoginLogGORMRepository 创建登录日志仓储
func NewLoginLogGORMRepository(db *gorm.DB) *LoginLogGORMRepository {
	return &LoginLogGORMRepository{db: db}
}

// Create 创建登录日志
func (r *LoginLogGORMRepository) Create(ctx context.Context, log *model.UserLoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// FindByUserID 根据用户ID查询登录日志
func (r *LoginLogGORMRepository) FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginLog, int64, error) {
	var logs []model.UserLoginLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserLoginLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// Paginate 分页查询（支持日期范围过滤）
func (r *LoginLogGORMRepository) Paginate(ctx context.Context, page, pageSize int, userID uint, startDate, endDate string) ([]model.UserLoginLog, int64, error) {
	var logs []model.UserLoginLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserLoginLog{})

	// 用户ID过滤
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// 日期范围过滤
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	// 计数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
