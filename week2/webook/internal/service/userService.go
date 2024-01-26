package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-work/week2/webook/internal/domain"
	"go-work/week2/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户邮箱或者密码错误")
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}
func (userService *UserService) SignUp(ctx context.Context, user domain.User) error {
	bcPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(bcPassword)
	return userService.userRepository.Create(ctx, user)
}

func (userService *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := userService.userRepository.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	//查出用户后检测用户密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

func (userService *UserService) Edit(ctx *gin.Context, user domain.User) error {
	return userService.userRepository.UpdateUserInfo(ctx, user)
}

func (userService *UserService) Profile(ctx *gin.Context, uid int64) (domain.User, error) {
	return userService.userRepository.FindById(ctx, uid)
}
