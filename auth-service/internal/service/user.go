package service

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "cybermind/auth-service/internal/model"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Phone    string `json:"phone" binding:"required,len=11"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
    // 检查用户名是否已存在
    var existingUser model.User
    if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
        return nil, errors.New("username or email already exists")
    } else if err != gorm.ErrRecordNotFound {
        return nil, err
    }

    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // 判断是否为管理员账号
    role := 0
    if req.Email == "2986309418@qq.com" {
        role = 1 // 设置为管理员角色
    }

    user := &model.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        Phone:        req.Phone,
        Status:       1,
        Role:         role,
    }

    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }

    return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*model.User, error) {
    var user model.User
    if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    if user.Status != 1 {
        return nil, errors.New("user is disabled")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid password")
    }

    return &user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
    var user model.User
    if err := s.db.First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
} 