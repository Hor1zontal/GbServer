package service

import (
	"aliens/common/character"
	"aliens/log"
	"gangbu/constant"
	"gangbu/exception"
	"gangbu/module/game/config"
	"gangbu/module/game/db"
	"gangbu/module/game/service/gumBalls"
	"gangbu/module/game/service/myjwt"
	"gangbu/module/game/service/wx"
	"gangbu/module/game/words"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func LoginByCode(code string, channel int32) (*db.DBUser) {
	username := ""
	nickname := ""
	openid := ""
	//session_key := ""
	switch channel {
	case constant.CHANNEL_WEB:
		username = "gb_" + code
		nickname = words.RandomName()
	case constant.CHANNEL_WECHAT:
		var err error
		openid, _, err = wx.WxLogin(code)
		username = "gb_" + openid
		if err != nil {
			exception.GameExceptionCustom("LoginByCode-WxLogin", exception.WxLogin, err)
		}
	}
	user, isNew := GetUserUsername(username, openid, nickname, channel)
	user.LoginLog(isNew, channel)
	return user
}

//获取user信息若不存在，创建
func GetUserUsername(username, openid, nickname string, channel int32) (*db.DBUser, bool) {
	user := db.GetUserInfo(username)
	isNew := false
	if user == nil {
		userNew := db.NewUser(username, openid, nickname, channel)
		user = userNew
		isNew = true
	}
	return user, isNew
}

func UpdateUser(uid int32, nickname, avatar string) {
	user := db.GetUserByUid(uid)
	if user == nil {
		exception.GameException(exception.UserNotFound)
	}
	user.Nickname = nickname
	user.Avatar = avatar
	db.UpdateOne(user)
}


//func NewUserByOpenID(openid, session_key string) *db.DBUser {
//	username := "gb_" + openid
//	user := db.NewUser(username, openid, "")
//	return user
//}

// 生成令牌
func GenerateToken(uid int32) string {
	j := myjwt.NewJWT(config.Server.JWTSecret)
	claims := myjwt.CustomClaims{
		UID:uid,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),	// 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + config.Server.JWTExpired),	// 过期时间 一小时
			Issuer:    "aliens",					//签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		log.Error(err.Error())
		return ""
	}

	//log.Debug(token)
	return token
}

func GetGumBallsGift(uid int32) (int32, string, string) {
	status := ""
	message := ""
	user := db.GetUserByUid(uid)
	if user == nil {
		exception.GameException(exception.UserNotFound)
	}
	openID := user.OpenID
	if openID == "" {
		exception.GameException(exception.OpenidIsNil)
	}
	role := db.GetRoleByUid(uid)
	if role.GiftCode == "" {
		status, message = gumBalls.DrawGift(openID)
		code := character.StringToInt32(status)
		if code == 1 || code == 60009 && role.GiftCode != message {
			role.GiftCode = message
			db.UpdateOne(role)
		}
		//log.Debug("status:%v message:%v", status, message)
	}
	return character.StringToInt32(status), message, role.GiftCode
}

func DeleteUser(uid int32) {
	db.DeleteUser(uid)
	db.DeleteRole(uid)
}

func GetFlags(uid int32) []*db.DBFlag {
	flags := db.GetFlags(uid)
	for _, flag := range flags {
		_, refresh := constant.FLAG_NEED_REFRESH[flag.ID.Flag]
		if refresh {
			flag.RefreshFlag()
			db.UpdateOne(flag)
		}
	}
	return flags
}

func UseFlag(uid int32, flag_id constant.ROLE_FLAG) {
	confValue, ok := constant.FLAG_CONFIG_VALUE[flag_id]
	if !ok {
		exception.GameException(exception.FlagNotFound)
	}
	flag := db.GetFlag(uid, flag_id, confValue)
	if !flag.UseValue() {
		exception.GameException(exception.FlagValueNotEnough)
	}
	db.UpdateOne(flag)
}

func SetFlagTrue(uid int32, flag_id constant.ROLE_FLAG) {
	flag := db.GetFlag(uid, flag_id,  constant.FLAG_VALUE_FALSE)
	if flag.Value == constant.FLAG_VALUE_TRUE {
		return
	}
	flag.SetValue(constant.FLAG_VALUE_TRUE)
	db.UpdateOne(flag)
}