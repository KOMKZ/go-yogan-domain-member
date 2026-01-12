# go-yogan-domain-member

会员领域包，提供 C 端用户（会员）的 CRUD 功能。

## 特性

- 会员 CRUD 操作
- 密码加密与验证
- 登录日志记录
- 批量操作支持
- 实现 `auth.Authenticatable` 接口

## 模型

- `User` - 会员实体（表名：users）
- `UserLoginLog` - 登录日志

## 使用方式

```go
repo := member.NewGORMRepository(db)
service := member.NewService(repo, logger)

// 创建会员
user, err := service.Create(ctx, &member.CreateInput{
    Name:     "张三",
    Email:    "zhangsan@example.com",
    Password: "password123",
})
```

## 依赖

- `github.com/KOMKZ/go-yogan-framework`
