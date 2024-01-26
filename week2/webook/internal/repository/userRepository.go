package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-work/week2/webook/internal/domain"
	"go-work/week2/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrRecordNotFound
)

type UserRepository struct {
	userDao *dao.UserDao
}

func NewUserRepository(userDao *dao.UserDao) *UserRepository {
	return &UserRepository{
		userDao: userDao,
	}
}
func (userRepository *UserRepository) Create(ctx context.Context, user domain.User) error {
	return userRepository.userDao.Insert(ctx, dao.UserEntity{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (userRepository *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := userRepository.userDao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return userRepository.toDomain(user), err

}

func (userRepository *UserRepository) toDomain(entity dao.UserEntity) domain.User {
	return domain.User{
		Id:       entity.Id,
		Email:    entity.Email,
		Password: entity.Password,
		Nickname: entity.Nickname,
		AboutMe:  entity.AboutMe,
		Birthday: entity.Birthday,
	}
}

func (userRepository *UserRepository) UpdateUserInfo(ctx *gin.Context, user domain.User) error {
	return userRepository.userDao.Update(ctx, dao.UserEntity{
		Id:       user.Id,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	})
}

func (userRepository *UserRepository) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	user, err := userRepository.userDao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return userRepository.toDomain(user), err
}
