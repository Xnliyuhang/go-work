package user

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-work/week2/webook/internal/domain"
	"go-work/week2/webook/internal/service"
	"net/http"
	"time"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	userService    *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		userService:    userService,
	}
}

func (userHandler *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", userHandler.Signup)
	ug.POST("/login", userHandler.Login)
	ug.GET("/profile", userHandler.Profile)
	ug.POST("/edit", userHandler.Edit)
}

const (
	emailRegexPattern    = "^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$"
	passwordRegexPattern = "^(?![a-zA-Z]+$)(?!\\d+$)(?![^\\da-zA-Z\\s]+$).{1,64}$" //由字母、数字、特殊字符，任意2种组成，1-9位
)

func (userHandler *UserHandler) Signup(context *gin.Context) {
	type SignupRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var signupReq SignupRequest
	if err := context.Bind(&signupReq); err != nil {
		return
	}

	isEmail, err := userHandler.emailRexExp.MatchString(signupReq.Email)
	if err != nil {
		context.String(http.StatusOK, "系统错误(执行超时)")
		return
	}
	if !isEmail {
		context.String(http.StatusOK, "邮箱格式错误")
		return
	}

	isPassword, err := userHandler.passwordRexExp.MatchString(signupReq.Password)
	if err != nil {
		context.String(http.StatusOK, "系统错误(执行超时)")
		return
	}
	if !isPassword {
		context.String(http.StatusOK, "密码格式错误")
		return
	}

	if signupReq.Password != signupReq.ConfirmPassword {
		context.String(http.StatusOK, "两次密码输入不一致")
		return
	}

	err = userHandler.userService.SignUp(context, domain.User{
		Email:    signupReq.Email,
		Password: signupReq.Password,
	})
	switch err {
	case nil:
		context.String(http.StatusOK, "完成注册")
	case service.ErrUserDuplicateEmail:
		context.String(http.StatusOK, "邮箱冲突")
	default:
		context.String(http.StatusOK, "系统错误")

	}

}

func (userHandler *UserHandler) Login(context *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest
	if err := context.Bind(&loginReq); err != nil {
		return
	}

	user, err := userHandler.userService.Login(context, loginReq.Email, loginReq.Password)
	switch err {
	case nil:

		session := sessions.Default(context)
		session.Set("userID", user.Id)
		session.Options(sessions.Options{
			//保存时间15分钟
			MaxAge:   900,
			HttpOnly: true,
		})
		err = session.Save()
		if err != nil {
			context.String(http.StatusOK, "系统错误")
			return
		}
		context.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		context.String(http.StatusOK, "用户邮箱或者密码错误")
	default:
		context.String(http.StatusOK, "系统错误")
	}

}

func (userHandler *UserHandler) Profile(context *gin.Context) {
	session := sessions.Default(context)
	uid := session.Get("userID").(int64)
	user, err := userHandler.userService.Profile(context, uid)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}

	context.JSON(http.StatusOK, User{
		Nickname: user.Nickname,
		Email:    user.Email,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday.Format(time.DateOnly),
	})
}

func (userHandler *UserHandler) Edit(context *gin.Context) {
	type editRequest struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var editReq editRequest
	if err := context.Bind(&editReq); err != nil {
		return
	}

	session := sessions.Default(context)
	uid := session.Get("userID").(int64)

	birthday, err := time.Parse(time.DateOnly, editReq.Birthday)
	if err != nil {
		context.String(http.StatusOK, "生日格式错误")
		return
	}

	err = userHandler.userService.Edit(context, domain.User{
		Id:       uid,
		Nickname: editReq.Nickname,
		Birthday: birthday,
		AboutMe:  editReq.AboutMe,
	})
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}

	context.String(http.StatusOK, "修改用户基本信息成功")
}
