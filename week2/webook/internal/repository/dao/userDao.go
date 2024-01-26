package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound     = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (userDao *UserDao) Insert(ctx context.Context, user UserEntity) error {
	now := time.Now().UnixMilli()
	user.Utime = now
	user.Ctime = now
	err := userDao.db.WithContext(ctx).Create(&user).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number == 1062 {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (userDao *UserDao) FindByEmail(ctx context.Context, email string) (UserEntity, error) {
	var user UserEntity
	err := userDao.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	return user, err
}

func (userDao *UserDao) Update(ctx *gin.Context, user UserEntity) error {
	return userDao.db.WithContext(ctx).Model(&user).Update("nickname", user.Nickname).Update("birthday", user.Birthday).Update("about_me", user.AboutMe).Error
}

func (userDao *UserDao) FindById(ctx *gin.Context, uid int64) (UserEntity, error) {
	var user UserEntity
	err := userDao.db.WithContext(ctx).First(&user, uid).Error
	return user, err
}

type UserEntity struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Nickname string `gorm:"type=varchar(128)"`
	Birthday time.Time
	AboutMe  string `gorm:"type=varchar(4096)"`
	Ctime    int64
	Utime    int64
}
