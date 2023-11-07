package v1

import (
	"authentik-go/app/constant"
	"authentik-go/app/helper"
	"authentik-go/app/interfaces"
	"authentik-go/app/service"
)

var ipRegistry = make(map[string]int)

// @Tags System
// @Summary 登录
// @Description 登录
// @Accept json
// @Param request body interfaces.UserLoginReq true "request"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/login [post]
func (api *BaseApi) Login() {
	var param = interfaces.UserLoginReq{}
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
	//
	result, err := service.UserService.UserLogin(param.Username, param.Password)
	if err != nil {
		helper.ApiResponse.Error(api.Context, err.Error())
		return
	}
	//
	helper.ApiResponse.Success(api.Context, result)
}

// @Tags System
// @Summary 注册
// @Description 注册
// @Accept json
// @Param request body interfaces.UserRegReq true "request"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/register [post]
func (api *NotAuthBaseApi) Register() {
	var param = interfaces.UserRegReq{}
	if err := api.BaseApi.Context.ShouldBindJSON(&param); err != nil {
		helper.ApiResponse.Error(api.BaseApi.Context, constant.ErrInvalidParameter)
		return
	}
	//
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
	//
	helper.ApiResponse.Success(api.BaseApi.Context, result)
}
