package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go-work/week2/webook/internal/repository"
	"go-work/week2/webook/internal/repository/dao"
	"go-work/week2/webook/internal/service"
	"go-work/week2/webook/internal/web/middleware"
	"go-work/week2/webook/internal/web/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {

	db := initDB()

	server := initWebServer()

	initUser(db, server)

	err := server.Run(":8081")
	if err != nil {
		panic(err)
	}

}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "root:123456@tcp(47.97.188.171:13316)/webook?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:        258,  // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision: true, // disable datetime precision support, which not supported before MySQL 5.6
		//DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,  // drop 订单 11 create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  //
		SkipInitializeWithVersion: false, // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)

	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "xn_lyh.com")
		},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}))

	login := &middleware.LoginMiddlewareBuilder{}
	//存储userID的地方
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store))
	server.Use(login.CheckLogin())
	return server
}

func initUser(db *gorm.DB, server *gin.Engine) {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	hdl := user.NewUserHandler(userService)
	hdl.RegisterRoutes(server)
}
