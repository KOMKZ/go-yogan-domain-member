package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体（C端用户）
// 与数据库表 users 对应
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`              // 用户名
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"` // 邮箱（唯一）
	Password  string         `gorm:"column:password;size:255" json:"-"`          // 密码（bcrypt加密，不返回）
	Age       int            `gorm:"default:0" json:"age"`                       // 年龄
	Phone     *string        `gorm:"size:20" json:"phone,omitempty"`             // 手机号
	Avatar    *string        `gorm:"size:500" json:"avatar,omitempty"`           // 头像
	Status    int            `gorm:"default:1" json:"status"`                    // 状态：0=禁用，1=启用
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 兼容字段（auth-app使用，数据库表中可能不存在）
	Username      string     `gorm:"-" json:"username,omitempty"`       // 兼容旧API
	AvatarURL     string     `gorm:"-" json:"avatar_url,omitempty"`     // 兼容旧API
	EmailVerified bool       `gorm:"-" json:"email_verified,omitempty"` // 兼容旧API
	PhoneVerified bool       `gorm:"-" json:"phone_verified,omitempty"` // 兼容旧API
	LastLoginAt   *time.Time `gorm:"-" json:"last_login_at,omitempty"`  // 兼容旧API
	LoginCount    uint       `gorm:"-" json:"login_count,omitempty"`    // 兼容旧API
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 检查账户是否激活
func (u *User) IsActive() bool {
	return u.Status == 1
}

// IsDisabled 检查账户是否禁用
func (u *User) IsDisabled() bool {
	return u.Status == 0
}

// ==================== Authenticatable 接口实现 ====================

// GetID 实现 auth.Authenticatable 接口
func (u *User) GetID() uint {
	return u.ID
}

// GetEmail 实现 auth.Authenticatable 接口
func (u *User) GetEmail() string {
	return u.Email
}

// GetPasswordHash 实现 auth.Authenticatable 接口
func (u *User) GetPasswordHash() string {
	return u.Password
}

// UserLoginLog 用户登录日志
type UserLoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	Username  string    `gorm:"size:50" json:"username"`
	IP        string    `gorm:"size:50;not null" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"userAgent"`
	DeviceID  string    `gorm:"size:100" json:"deviceId"`
	City      string    `gorm:"size:100" json:"city"`
	Country   string    `gorm:"size:100" json:"country"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName 指定表名
func (UserLoginLog) TableName() string {
	return "user_login_logs"
}
