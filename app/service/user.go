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

// 用户注册
func (uSrv userService) UserReg(clientIp string, param interfaces.UserRegReq) (*model.AuthentikCoreUser, error) {
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
