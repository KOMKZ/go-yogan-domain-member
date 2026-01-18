package member

import (
	"github.com/KOMKZ/go-yogan-framework/event"
)

// 事件名称常量
const (
	EventUserRegistered = "user.registered"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	event.BaseEvent
	UserID   uint   // 用户ID
	UserName string // 用户名
	Email    string // 邮箱
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID uint, userName, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: event.NewEvent(EventUserRegistered),
		UserID:    userID,
		UserName:  userName,
		Email:     email,
	}
}
