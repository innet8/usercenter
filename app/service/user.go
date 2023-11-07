package service

import (
	"authentik-go/app/interfaces"
	"authentik-go/app/model"
	"authentik-go/core"
	"authentik-go/utils/common"
	"authentik-go/utils/django"
	"time"

	e "authentik-go/utils/error"

	"github.com/google/uuid"
)

var UserService = userService{}

type userService struct{}

// 用户登录
func (uSrv userService) UserLogin(username, password string) (*model.AuthentikCoreUser, error) {
	var user model.AuthentikCoreUser
	core.DB.Where("username = ? OR email = ?", username, username).First(&user)
	if user.ID == 0 {
		return nil, e.New("用户不存在")
	}
	if !django.CheckPassword(password, user.Password) {
		return nil, e.New("密码错误")
	}
	if !user.IsActive {
		return nil, e.New("该帐户已被禁用")
	}
	// 更新登录时间
	core.DB.Model(&user).Where("id = ?", user.ID).Update("last_login", time.Now())
	//
	var result model.AuthentikCoreUser
	err := common.StructToStruct(user, &result, "password")
	if err != nil {
		return nil, err
	}
	//
	return &result, nil
}

// 用户注册
func (uSrv userService) UserReg(clientIp string, param interfaces.UserRegReq) (*model.AuthentikCoreUser, error) {
	var count int64
	core.DB.Model(&model.AuthentikCoreUser{}).Where("date_joined > ?", time.Now().Add(-1*time.Minute)).Count(&count)
	if count >= 3 {
		return nil, e.New("用户注册失败，操作频率过高")
	}
	// 邮箱
	if !common.IsEmail(param.Email) {
		return nil, e.New("无效的邮箱地址")
	}
	//
	var user model.AuthentikCoreUser
	core.DB.Where("username = ?", param.Email).First(&user)
	if user.ID > 0 {
		return nil, e.New("用户已存在")
	}
	//
	err := core.DB.Model(&model.AuthentikCoreUser{}).Create(map[string]interface{}{
		"Path":               param.Source,
		"Source":             param.Source,
		"Username":           param.Email,
		"Name":               param.Email,
		"Email":              param.Email,
		"Password":           django.GeneratePassword(param.Password),
		"UUID":               uuid.New(),
		"DateJoined":         time.Now(),
		"IsActive":           true,
		"LastName":           "",
		"FirstName":          "",
		"PasswordChangeDate": time.Now(),
		"Attributes":         "{}",
		"Type":               "internal",
		"RegIP":              clientIp,
	}).Error
	if err != nil {
		return nil, err
	}
	//
	var userModel model.AuthentikCoreUser
	core.DB.Where("username = ? OR email = ?", param.Email).First(&userModel)
	//
	return &userModel, nil
}
