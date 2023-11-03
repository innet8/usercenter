package service

import (
	"authentik-go/app/interfaces"
	"authentik-go/app/model"
	"authentik-go/core"
	"authentik-go/utils/common"
	"authentik-go/utils/django"

	e "authentik-go/utils/error"
)

var UserService = userService{}

type userService struct{}

// 生成 token
// func (p userService) GenerateToken(user *model.AuthentikCoreUser, refresh bool) string {
// 	var token string
// 	if refresh {
// 		days := 0
// 		if user.Bot != 0 {
// 			systemSetting := SystemProviders.SettingFind("system", "token_valid_days", 30)
// 			if systemSetting != nil {
// 				days = systemSetting.(int)
// 			}
// 		}
// 		token = p.TokenEncode(user.ID, user.Email, user.Encrypt, days)
// 	} else {
// 		token = p.UserToken(user.Userid)
// 	}
// 	// user.Encrypt = ""
// 	// user.Password = ""
// 	user.Token = token
// 	return user.Token
// }

// 生成token（编码token）
// func (p userService) TokenEncode(userid int, email, encrypt string, days int) string {
// 	if days == 0 {
// 		days = 15 // 默认有效时间（天）
// 	}
// 	// 创建JWT令牌
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
// 		Userid:  userid,
// 		Email:   email,
// 		Encrypt: encrypt,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix(),
// 		},
// 	})
// 	// 签名JWT令牌
// 	tokenString, err := token.SignedString([]byte(config.CONF.Jwt.SecretKey))
// 	if err != nil {
// 		return ""
// 	}
// 	return tokenString
// }

// 登录
// func (uSrv userService) Login(clientIp string, userModel *model.User) error {
// 	tokenString := uSrv.TokenEncode(userModel.Userid, userModel.Email, userModel.Encrypt, 0)

// 	// 更新登录信息
// 	err := core.DB.Model(&userModel).Updates(model.User{
// 		LoginNum: userModel.LoginNum + 1,
// 		LastIp:   clientIp,
// 		LastAt:   core.TsTime(time.Now().Unix()),
// 		Token:    tokenString,
// 	}).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// 用户登录
// func (uSrv userService) UserLogin(clientIp string, param interfaces.UserLoginReq) (*model.User, map[string]any, error) {
// 	var userModel model.User
// 	//
// 	codes := map[string]any{"code": "need"}
// 	if model.SettingModel.IsNeedCode(param.Email) {
// 		if len(param.Code) == 0 || len(param.CodeId) == 0 {
// 			return nil, codes, e.New(constant.ErrCaptchaCode)
// 		}
// 		err := captcha.VerifyCode(param.CodeId, param.Code)
// 		if err != nil {
// 			return nil, codes, e.New(constant.ErrCaptchaCode)
// 		}
// 	}
// 	//
// 	err := core.DB.Where("email = ?", param.Email).First(&userModel).Error
// 	if err != nil || userModel.Password != common.StringMd52(param.Password, userModel.Encrypt) {
// 		core.Cache.Set("code::"+param.Email, "need", 2*time.Hour)
// 		return nil, codes, e.New(constant.ErrPasswordFailed) //帐号或密码错误
// 	}
// 	if err = uSrv.Login(clientIp, &userModel); err != nil {
// 		return nil, map[string]any{}, err
// 	}
// 	core.Cache.Set("code::"+param.Email, "no", 2*time.Hour)
// 	// 去除重要信息返回
// 	userModel.FilterSensitiveFields()
// 	//
// 	return &userModel, map[string]any{}, nil
// }

// 用户注册
func (uSrv userService) UserReg(clientIp string, param interfaces.UserRegReq) (*model.AuthentikCoreUser, error) {
	var userModel model.AuthentikCoreUser
	//
	// 邮箱
	if !common.IsEmail(param.Email) {
		return nil, e.New("无效的邮箱地址")
	}
	// //检测密码策略是否符合
	// err := common.PasswordPolicy(param.Password, (system["password_policy"] == "complex"))
	// if err != nil {
	// 	return nil, err, false
	// }
	//
	userModel.Path = param.Source
	userModel.Source = param.Source
	userModel.Username = param.Email
	userModel.Name = param.Email
	userModel.Email = param.Email
	userModel.Password = django.GeneratePassword(param.Password)
	userModel.Email = param.Email
	err := core.DB.Create(&userModel).Error
	if err != nil {
		return nil, err
	}
	//
	return &userModel, nil
}

// // UserEditData 用户编辑/个人设置
// func (uSrv userService) UserEditData(userID int, param interfaces.UserEditDataReq) (*model.User, error) {
// 	paramMap, _ := common.StructToMap(param)
// 	update, err := model.UserModel.UpdateUserByID(userID, paramMap)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return update, nil
// }

// // UserEditPassword 密码设置
// func (uSrv userService) UserEditPassword(userID int, param interfaces.UserEditPasswordReq) (*model.User, error) {
// 	user, err := model.UserModel.GetUserByID(userID, true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//
// 	if param.OldPassword == param.NewPassword {
// 		return nil, e.New(constant.ErrUserPasswordSame)
// 	}
// 	if user.Password != common.StringMd52(param.OldPassword, user.Encrypt) {
// 		return nil, e.New(constant.ErrUserOldPassword)
// 	}
// 	//检测密码策略是否符合
// 	system, _ := model.SettingModel.GetSetting(model.SettingSystemKey)
// 	err = common.PasswordPolicy(param.NewPassword, (system["password_policy"] == "complex"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	//
// 	user.Encrypt = common.RandString(8)
// 	user.Password = common.StringMd52(param.NewPassword, user.Encrypt)
// 	if err = core.DB.Save(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	user.FilterSensitiveFields()
// 	return user, nil
// }

// // UserEditEmail 修改邮箱
// func (uSrv userService) UserEditEmail(userID int, param interfaces.UserEditEmailReq) (*model.User, error) {
// 	paramMap, _ := common.StructToMap(param)
// 	update, err := model.UserModel.UpdateUserByID(userID, paramMap)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return update, nil
// }

// // DeleteAccount 删除帐号
// func (uSrv userService) DeleteAccount(user *model.User, reason string) error {
// 	result := core.DB.Create(&model.UserDelete{
// 		Userid:   user.Userid,
// 		Operator: user.Userid,
// 		Email:    user.Email,
// 		Reason:   reason,
// 		Cache:    common.StructToJson(user),
// 	})
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	// 删除未读
// 	core.DB.Where("userid = ?", user.Userid).Delete(&model.WebSocketDialogMsgRead{})
// 	// 删除待办
// 	core.DB.Where("userid = ?", user.Userid).Delete(&model.WebSocketDialogMsgTodo{})
// 	// 删除邮箱验证记录
// 	core.DB.Where("email = ?", user.Email).Delete(&model.UserEmailVerification{})
// 	//
// 	core.DB.Where("email = ?", user.Email).Delete(&user)
// 	//
// 	return nil
// }

// // UserDeleteAccount 用户删除帐号
// func (uSrv userService) UserDeleteAccount(user *model.User, param interfaces.UserDeleteAccountReq) error {
// 	emailSetting, _ := model.SettingModel.GetSetting(model.SettingEmailKey)
// 	if emailSetting["reg_verify"] == "open" {
// 		_, err := UserEmailService.UserEmailVerification(param.Code)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		if len(param.Password) == 0 {
// 			return e.New(constant.ErrPasswordEnterFailed)
// 		}
// 		if user.Password != common.StringMd52(param.Password, user.Encrypt) {
// 			return e.New(constant.ErrPasswordFailed)
// 		}
// 	}
// 	if param.Type == "confirm" {
// 		uSrv.DeleteAccount(user, param.Reason)
// 	}
// 	return nil
// }

// // isExistEmail 邮箱是否存在
// func (uSrv userService) IsExistEmail(email string) bool {
// 	var user *model.User
// 	err := core.DB.Where("email = ?", email).Find(&user).Error
// 	if err != nil {
// 		return false
// 	}
// 	if user.Userid > 0 {
// 		return true
// 	}
// 	return false
// }

// // GetUserBasic 获取指定会员基础信息
// func (uSrv userService) GetUserBasic(user *model.User, ids []string) ([]map[string]any, error) {
// 	var userList []*model.User
// 	err := core.DB.Model(&model.User{}).Where("userid in ?", ids).Find(&userList).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	userInfoList := []map[string]any{}
// 	for _, v := range userList {
// 		v.Token = ""
// 		v.FilterSensitiveFields()
// 		userMap, _ := common.StructToMap(v)
// 		userMap["online"] = v.GetOnlineStatus()
// 		userMap["department_name"] = v.GetDepartmentName()
// 		userInfoList = append(userInfoList, userMap)
// 	}
// 	return userInfoList, nil
// }

// // UserSearch 搜索会员列表
// func (uSrv userService) UserSearch(user *model.User, param interfaces.UserSearchReq, page, pageSize int) (*interfaces.Pagination, error) {
// 	projectUserTableName := core.DBTableName(&model.ProjectUser{})
// 	webSocketDialogUserTableName := core.DBTableName(&model.WebSocketDialogUser{})
// 	db := core.DB.Table(core.DBTableName(&model.User{})).Select(model.UserBasicField)
// 	if param.Keys["key"] != nil {
// 		key := "%" + param.Keys["key"].(string) + "%"
// 		if v, ok := param.Keys["key"].(string); ok && strings.Contains(v, "@") {
// 			db = db.Where("email like ?", key)
// 		} else {
// 			db = db.Where("(name like '%?%' or pinyin like '%?%')", key, key)
// 		}
// 	}
// 	//
// 	if param.Keys["disable"] == nil || param.Keys["disable"] == 0 {
// 		db = db.Where("(disable_at = 0 OR disable_at IS NULL)")
// 	} else {
// 		db = db.Where("disable_at > 0")
// 	}
// 	//
// 	if param.Keys["bot"] == nil || param.Keys["bot"] == 0 {
// 		db = db.Where("bot = 0")
// 	} else if param.Keys["bot"] == 1 {
// 		db = db.Where("bot = 1")
// 	}
// 	//
// 	if param.UpdatedTime > 0 {
// 		db = db.Where("updated_at >= ?", time.Unix(param.UpdatedTime, 0))
// 	}
// 	//
// 	if param.Keys["project_id"] != nil {
// 		db = db.Where("id IN (?)", core.DB.Table(projectUserTableName).Select("userid").Where("project_id = ?", param.Keys["project_id"]))
// 	}
// 	//
// 	if param.Keys["no_project_id"] != nil {
// 		db = db.Where("id NOT IN (?)", core.DB.Table(projectUserTableName).Select("userid").Where("project_id = ?", param.Keys["no_project_id"]))
// 	}
// 	//
// 	if param.Keys["dialog_id"] != nil {
// 		db = db.Where("id IN (?)", core.DB.Table(webSocketDialogUserTableName).Select("dialog_id").Where("dialog_id = ?", param.Keys["dialog_id"]))
// 	}
// 	//
// 	if param.Sorts["az"] != nil && strings.Contains("asc desc", param.Sorts["az"].(string)) {
// 		db = db.Order("az " + param.Sorts["az"].(string))
// 	}
// 	//
// 	var count int64
// 	var userList []model.User
// 	if err := db.Order("id").Offset((page - 1) * pageSize).Limit(pageSize).Find(&userList).Count(&count).Error; err != nil {
// 		return nil, err
// 	}
// 	//
// 	results := []map[string]any{}
// 	for _, v := range userList {
// 		tags := []string{}
// 		re := regexp.MustCompile(`\Q(M)\E$`)
// 		parts := strings.Split(v.GetDepartmentName(), ",")
// 		for _, part := range parts {
// 			if re.MatchString(part) {
// 				tag := re.ReplaceAllString(strings.Trim(part, " "), "") + "负责人"
// 				tags = append(tags, tag)
// 			}
// 		}
// 		if user.IsAdmin() {
// 			if v.IsAdmin() {
// 				tags = append(tags, "系统管理员")
// 			}
// 			if v.IsTemp() {
// 				tags = append(tags, "临时帐号")
// 			}
// 			if v.Userid > 3 && time.Unix(int64(v.CreatedAt), 0).After(time.Now().UTC().AddDate(0, 0, -30)) {
// 				tags = append(tags, "新帐号")
// 			}
// 		}
// 		//
// 		userMap, _ := common.StructToMap(v)
// 		userMap["tags"] = tags
// 		if param.State == 1 {
// 			userMap["online"] = v.GetOnlineStatus()
// 		}
// 		results = append(results, userMap)
// 	}
// 	//
// 	return interfaces.PaginationRsp(page, pageSize, count, results), nil
// }

// // UserList 会员列表
// func (uSrv userService) UserList(user *model.User, param interfaces.UserListReq, page, pageSize int) (*interfaces.Pagination, error) {
// 	if !user.IsAdmin() {
// 		return nil, e.New(constant.ErrNoPermission)
// 	}
// 	getCheckinMac := param.GetCheckinMac == 1
// 	db := core.DB.Table(core.DBTableName(&model.User{})).Select("*, nickname as name_original")

// 	// 所有查询条件处理
// 	if param.Keys != nil {
// 		//
// 		if param.Keys["key"] != nil {
// 			key := "%" + param.Keys["key"].(string) + "%"
// 			if v, ok := param.Keys["key"].(string); ok && strings.Contains(v, "@") {
// 				db = db.Where("email like ?", key)
// 			} else {
// 				db = db.Where("(email like ? or tel like ? or nickname like ? or profession like ?)", key, key, key, key)
// 			}
// 		} else {
// 			if param.Keys["email"] != nil && param.Keys["email"] != "" {
// 				db = db.Where("email like ?", "%"+param.Keys["email"].(string)+"%")
// 			}
// 			if param.Keys["tel"] != nil && param.Keys["tel"] != "" {
// 				db = db.Where("tel like ?", "%"+param.Keys["tel"].(string)+"%")
// 			}
// 			if param.Keys["name"] != nil && param.Keys["name"] != "" {
// 				db = db.Where("nickname like ?", "%"+param.Keys["name"].(string)+"%")
// 			}
// 			if param.Keys["profession"] != nil && param.Keys["profession"] != "" {
// 				db = db.Where("profession like ?", "%"+param.Keys["profession"].(string)+"%")
// 			}
// 		}
// 		//
// 		if param.Keys["identity"] != nil {
// 			if common.LeftExists(param.Keys["identity"].(string), "no") {
// 				db = db.Where("identity like ?", "%"+common.LeftDelete(param.Keys["identity"].(string), "no")+"%")
// 			} else {
// 				db = db.Where("identity like ?", "%"+param.Keys["identity"].(string)+"%")
// 			}
// 		}
// 		//
// 		if param.Keys["disable"] == "yes" {
// 			db = db.Where("disable_at > 0")
// 		} else if param.Keys["disable"] != "all" {
// 			db = db.Where("(disable_at = 0 OR disable_at IS NULL)")
// 		}
// 		//
// 		if param.Keys["email_verity"] == "yes" {
// 			db = db.Where("email_verity = 1")
// 		} else if param.Keys["email_verity"] == "no" {
// 			db = db.Where("email_verity = 0")
// 		}
// 		//
// 		if param.Keys["bot"] == "yes" {
// 			db = db.Where("bot = 1")
// 		} else if param.Keys["bot"] != "all" {
// 			db = db.Where("bot = 0")
// 		}
// 		//
// 		if param.Keys["department"] != nil {
// 			if param.Keys["department"] == "0" || param.Keys["department"] == 0 {
// 				db = db.Where("(department = '' OR department = ',,')")
// 			} else {
// 				mUserDepartment := core.DBTableName(&model.UserDepartment{})
// 				subDb := core.DB.Table(mUserDepartment).Select("owner_userid").Where("id = ?", strings.Trim(param.Keys["department"].(string), ","))
// 				db = db.Where("(department like ? OR id In (?))", "%"+param.Keys["department"].(string)+"%", subDb)
// 				db = db.Select("if(EXISTS(select id from ? where owner_userid = id and id=?),1,0) as is_principal", mUserDepartment, param.Keys["department"])
// 				db = db.Order("is_principal desc")
// 			}
// 		}
// 		//
// 		if getCheckinMac && param.Keys["checkin_mac"] != nil {
// 			db = db.Where("id IN ?", db.Table(core.DBTableName(&model.UserCheckinMac{})).Select("userid").Where("mac like ?", "%"+param.Keys["checkin_mac"].(string)+"%"))
// 		}
// 	} else {
// 		db = db.Where("(disable_at = 0 OR disable_at IS NULL)")
// 		db = db.Where("bot = 0")
// 	}
// 	// 当没有传部门id时，默认排序，不然以排泄调整为负责人放在第一位
// 	if param.Keys["department"] == nil || param.Keys["department"] == "0" || param.Keys["department"] == 0 {
// 		db = db.Order("userid desc")
// 	}
// 	// 获取总数
// 	var count int64
// 	var userList []model.User
// 	if err := db.Order("userid").Offset((page - 1) * pageSize).Limit(pageSize).Find(&userList).Count(&count).Error; err != nil {
// 		return nil, err
// 	}
// 	//
// 	results := []map[string]any{}
// 	for _, v := range userList {
// 		userMap, _ := common.StructToMap(v)
// 		if getCheckinMac {
// 			var userCheckinMacs []model.UserCheckinMac
// 			core.DB.Select("id,mac,remark").Where("userid = ?", v.Userid).Order("id").Find(&userCheckinMacs)
// 			userMap["checkin_macs"] = userCheckinMacs
// 		}
// 		results = append(results, userMap)
// 	}
// 	//
// 	return interfaces.PaginationRsp(page, pageSize, count, results), nil
// }
