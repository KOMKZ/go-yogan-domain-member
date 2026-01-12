package member

import (
	"context"
	"time"

	"github.com/KOMKZ/go-yogan-domain-member/model"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service 会员服务
type Service struct {
	repo   Repository
	logger *logger.CtxZapLogger
}

// NewService 创建会员服务
func NewService(repo Repository, log *logger.CtxZapLogger) *Service {
	return &Service{
		repo:   repo,
		logger: log,
	}
}

// CreateInput 创建会员输入
type CreateInput struct {
	Name     string
	Email    string
	Password string
	Age      int
	Phone    string
	Status   int
}

// Create 创建会员
func (s *Service) Create(ctx context.Context, input *CreateInput) (*model.User, error) {
	exists, err := s.repo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		s.logger.ErrorCtx(ctx, "检查邮箱失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	var passwordHash string
	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.ErrorCtx(ctx, "密码加密失败", zap.Error(err))
			return nil, ErrInternalError.Wrap(err)
		}
		passwordHash = string(hash)
	}

	var phone *string
	if input.Phone != "" {
		phone = &input.Phone
	}

	status := input.Status
	if status == 0 {
		status = 1
	}

	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  passwordHash,
		Age:       input.Age,
		Phone:     phone,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.ErrorCtx(ctx, "创建会员失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "创建会员成功", zap.Uint("user_id", user.ID), zap.String("email", user.Email))
	return user, nil
}

// GetByID 根据ID获取会员
func (s *Service) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorCtx(ctx, "查询会员失败", zap.Uint("id", id), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user, nil
}

// Paginate 分页查询会员
func (s *Service) Paginate(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.repo.Paginate(ctx, page, pageSize, keyword)
	if err != nil {
		s.logger.ErrorCtx(ctx, "分页查询会员失败", zap.Error(err))
		return nil, 0, ErrDatabaseError.Wrap(err)
	}
	return users, total, nil
}

// UpdateInput 更新会员输入
type UpdateInput struct {
	Name   *string
	Email  *string
	Age    *int
	Phone  *string
	Status *int
}

// Update 更新会员
func (s *Service) Update(ctx context.Context, id uint, input *UpdateInput) (*model.User, error) {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Email != nil && *input.Email != user.Email {
		exists, err := s.repo.ExistsByEmail(ctx, *input.Email)
		if err != nil {
			s.logger.ErrorCtx(ctx, "检查邮箱失败", zap.Error(err))
			return nil, ErrDatabaseError.Wrap(err)
		}
		if exists {
			return nil, ErrEmailExists
		}
		user.Email = *input.Email
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Age != nil {
		user.Age = *input.Age
	}
	if input.Phone != nil {
		if *input.Phone == "" {
			user.Phone = nil
		} else {
			user.Phone = input.Phone
		}
	}
	if input.Status != nil {
		user.Status = *input.Status
	}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.ErrorCtx(ctx, "更新会员失败", zap.Uint("id", id), zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "更新会员成功", zap.Uint("user_id", user.ID))
	return user, nil
}

// Delete 删除会员
func (s *Service) Delete(ctx context.Context, id uint) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.ErrorCtx(ctx, "删除会员失败", zap.Uint("id", id), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "删除会员成功", zap.Uint("user_id", id))
	return nil
}

// BatchDelete 批量删除会员
func (s *Service) BatchDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		if err := s.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ==================== 兼容旧 API ====================

// CreateUserWithPassword 创建会员（兼容旧API）
func (s *Service) CreateUserWithPassword(ctx context.Context, name, email string, age int, passwordHash string) (*model.User, error) {
	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.ErrorCtx(ctx, "检查邮箱失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	user := &model.User{
		Name:      name,
		Email:     email,
		Password:  passwordHash,
		Age:       age,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.ErrorCtx(ctx, "创建会员失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "创建会员成功", zap.Uint("user_id", user.ID), zap.String("email", user.Email))
	return user, nil
}

// GetUser 获取会员（兼容旧API）
func (s *Service) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.GetByID(ctx, id)
}
