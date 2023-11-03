package v1

import (
	"authentik-go/app/constant"
	"authentik-go/app/helper"
	"authentik-go/app/interfaces"
	"authentik-go/app/model"
	"authentik-go/app/service"
	"authentik-go/core"
	"authentik-go/utils/common"
	"authentik-go/utils/django"
	"time"
)

// @Tags System
// @Summary 登录
// @Description 登录
// @Accept json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/login [post]
func (api *BaseApi) Login() {
	var param = model.AuthentikCoreUser{}
	if err := api.Context.ShouldBind(&param); err != nil {
		helper.ApiResponse.Error(api.Context, constant.ErrInvalidParameter)
		return
	}
	if param.Username == "" {
		helper.ApiResponse.Error(api.Context, "用户名不能为空")
		return
	}
	if param.Password == "" {
		helper.ApiResponse.Error(api.Context, "密码不能为空")
		return
	}
	var user model.AuthentikCoreUser
	core.DB.Where("username = ? OR email = ?", param.Username, param.Username).First(&user)
	if user.ID == 0 {
		helper.ApiResponse.Error(api.Context, "用户不存在")
		return
	}
	if !django.CheckPassword(param.Password, user.Password) {
		helper.ApiResponse.Error(api.Context, "密码错误")
		return
	}
	if !user.IsActive {
		helper.ApiResponse.Error(api.Context, "该帐户已被禁用")
		return
	}
	// 更新登录时间
	core.DB.Model(&user).Where("id = ?", user.ID).Update("last_login", time.Now())
	//
	var result model.AuthentikCoreUser
	err := common.StructToStruct(user, &result, "password")
	if err != nil {
		helper.ApiResponse.Error(api.Context, err.Error())
		return
	}
	//
	helper.ApiResponse.Success(api.Context, result)
}

// @Tags System
// @Summary 获取客户端列表
// @Description 获取客户端列表
// @Accept json
// @Param source query string false "来源 all"
// @Param online query string false "是否在线 all,1,0"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/register [post]
func (api *NotAuthBaseApi) Register() {
	var param = interfaces.UserRegReq{}
	if err := api.BaseApi.Context.ShouldBindJSON(&param); err != nil {
		helper.ApiResponse.Error(api.BaseApi.Context, constant.ErrInvalidParameter)
		return
	}
	if len(param.Password) == 0 {
		helper.ApiResponse.Error(api.BaseApi.Context, "密码不能为空")
		return
	}
	if len(param.Email) == 0 {
		helper.ApiResponse.Error(api.BaseApi.Context, "邮箱不能为空")
		return
	}
	if len(param.Source) == 0 {
		helper.ApiResponse.Error(api.BaseApi.Context, "来源不能为空")
		return
	}
	result, err := service.UserService.UserReg(api.BaseApi.Context.ClientIP(), param)
	if err != nil {
		helper.ApiResponse.Error(api.BaseApi.Context, err.Error())
		return
	}
	helper.ApiResponse.Success(api.BaseApi.Context, result)
}