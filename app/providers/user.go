package providers

import (
	"dootask-go/app/model"
	"dootask-go/config"
	"dootask-go/core"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	e "dootask-go/utils/error"
)

type userProviders struct{}

var UserProviders = userProviders{}

type Claims struct {
	Userid  int    `json:"userid"`
	Email   string `json:"email"`
	Encrypt string `json:"encrypt"`
	jwt.StandardClaims
}

// IsExistUser 判断用户是否存在
func (p userProviders) IsExistUser(userid int) bool {
	user := model.User{}
	core.DB.Where("id = ?", userid).First(&user)
	return user.Userid > 0
}

// 生成 token
func (p userProviders) GenerateToken(user *model.User, refresh bool) string {
	var token string
	if refresh {
		days := 30 //（天）
		token = p.TokenEncode(user.Userid, user.Email, user.Encrypt, days)
	} else {
		token = p.UserToken(user.Userid)
	}
	user.Token = token
	return user.Token
}

// 生成token（编码token）
func (p userProviders) TokenEncode(userid int, email, encrypt string, days int) string {
	if days == 0 {
		days = 15 // 默认有效时间（天）
	}
	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Userid:  userid,
		Email:   email,
		Encrypt: encrypt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix(),
		},
	})
	// 签名JWT令牌
	tokenString, err := token.SignedString([]byte(config.CONF.Jwt.SecretKey))
	if err != nil {
		return ""
	}
	return tokenString
}

// 当前会员token（来自请求的token）
func (p userProviders) UserToken(userid int) string {
	user := model.User{}
	core.DB.Where("userid = ?", userid).First(&user)
	if user.Token == "" {
		user.Token = p.TokenEncode(user.Userid, user.Email, user.Encrypt, 0)
		core.DB.Model(&user).Update("token", user.Token)
	}
	return user.Token
}

// 验证登录
func (p userProviders) VerifyLogin(tokenString string) (*model.User, error) {
	if len(tokenString) == 0 {
		return nil, e.New("")
	}

	// 缓存解析后的 token 和 claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.CONF.Jwt.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if claims.Userid == 0 {
		return nil, e.New("")
	}

	// 缓存 user info
	info, err := model.UserModel.GetUserByID(int(claims.Userid))
	if err != nil {
		return nil, e.New("")
	}

	if info.Token != tokenString {
		return nil, e.New("")
	}

	// // 验证 claims 中所有值是否相等
	// userClaims := Claims{
	// 	Userid:  info.Userid,
	// 	Email:   info.Email,
	// 	Encrypt: info.Encrypt,
	// }
	// if !reflect.DeepEqual(userClaims, claims) {
	// 	return nil, e.New("")
	// }

	return info, nil
}
