package member

import (
	"context"
	"time"

	"github.com/KOMKZ/go-yogan-domain-member/model"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
)

// LoginLogService 登录日志服务
type LoginLogService struct {
	repo   LoginLogRepository
	logger *logger.CtxZapLogger
}

// NewLoginLogService 创建登录日志服务
func NewLoginLogService(repo LoginLogRepository, log *logger.CtxZapLogger) *LoginLogService {
	return &LoginLogService{
		repo:   repo,
		logger: log,
	}
}

// CreateLog 创建登录日志
func (s *LoginLogService) CreateLog(ctx context.Context, userID uint, username, ip, userAgent, deviceID, city, country string) error {
	log := &model.UserLoginLog{
		UserID:    userID,
		Username:  username,
		IP:        ip,
		UserAgent: userAgent,
		DeviceID:  deviceID,
		City:      city,
		Country:   country,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, log); err != nil {
		s.logger.ErrorCtx(ctx, "创建会员登录日志失败", zap.Uint("user_id", userID), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "会员登录日志创建成功", zap.Uint("user_id", userID), zap.String("ip", ip))
	return nil
}

// ListPageInput 分页查询输入
type ListPageInput struct {
	Page      int
	Size      int
	UserID    uint
	StartDate string
	EndDate   string
}

// ListPage 分页查询登录日志
func (s *LoginLogService) ListPage(ctx context.Context, input *ListPageInput) ([]model.UserLoginLog, int64, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Size <= 0 || input.Size > 100 {
		input.Size = 10
	}

	logs, total, err := s.repo.Paginate(ctx, input.Page, input.Size, input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		s.logger.ErrorCtx(ctx, "查询会员登录日志失败", zap.Error(err))
		return nil, 0, ErrDatabaseError.Wrap(err)
	}

	return logs, total, nil
}

// GetByUserID 根据会员ID获取登录日志
func (s *LoginLogService) GetByUserID(ctx context.Context, userID uint, page, pageSize int) ([]model.UserLoginLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	logs, total, err := s.repo.FindByUserID(ctx, userID, page, pageSize)
	if err != nil {
		s.logger.ErrorCtx(ctx, "查询会员登录日志失败", zap.Uint("user_id", userID), zap.Error(err))
		return nil, 0, ErrDatabaseError.Wrap(err)
	}

	return logs, total, nil
}
