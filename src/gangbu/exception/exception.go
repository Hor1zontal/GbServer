package exception

import "aliens/log"

type ErrorCode int32

/**
* @api {ERROR_CODE} 错误码
* @apiGroup ERROR
* @apiSuccess {Number} errcode 错误码
* @apiSuccess {String} errmsg 错误信息
* @apiSuccess (errcode) {Number} 10001 参数错误
* @apiSuccess (errcode) {Number} 10002 数据库操作异常
* @apiSuccess (errcode) {Number} 10003 内部异常
* @apiSuccess (errcode) {Number} 20001 微信登录异常
* @apiSuccess (errcode) {Number} 20002 无效的token
* @apiSuccess (errcode) {Number} 20003 未携带token
*/
const (
	CodeSuccess   ErrorCode = iota
	InvalidParam            = 10001 //参数错误
	DatabaseError           = 10002 //数据库操作异常
	InternalError           = 10003 //内部异常
	TimeParseError			= 10004 //时间格式错误
	SignError				= 10005 //签名错误

	WxLogin      			= 20001 //微信登录异常
	TokenInvalid 			= 20002 //无效的token
	TokenIsNil				= 20003 //未携带token
	TokenExpired			= 20004 //token过期
	TokenCheckError			= 20005 //token验证错误
	EnergyNotEnough			= 20006 //体力不足
	UpdateScoreError		= 20007 //上传分数错误
	GameDataNotFound		= 20008 //游戏数据未找到
	PropNotFound			= 20009 //游戏道具未找到
	PropNotEnough			= 20010 //道具不足
	RoleNotFound			= 20011 //用户角色未找到
	ComposeEquipError		= 20012 //装备合成失败
	OpenidIsNil				= 20013 //openId为空
	GumBallsGiftError		= 20014 //冈布奥礼包领取错误
	UserNotFound			= 20015 //用户未找到
	FlagNotFound			= 20016 //flag未找到
	FlagValueNotEnough		= 20017 //flag值不够
	//ErrorCodeUserNotFound           = 20002 //用户未找到
	AlreadyDrawLetter		= 20018 //已经领取过公开信奖励
)

func GameException(this ErrorCode) {
	panic(this)
}

var ErrorMapping = map[ErrorCode]string{
	CodeSuccess:	"success",
	InvalidParam:	"invalid param",
	DatabaseError:	"database error",
	InternalError:	"internal error",
	TimeParseError: "time parse error",
	SignError: "sign error",

	WxLogin:		"weChat login error",
	TokenInvalid:	"token invalid",
	TokenIsNil:		"token is nil",
	TokenExpired:	"token expired",
	TokenCheckError:"token check error",
	EnergyNotEnough:"energy not enough",
	UpdateScoreError:"you cheated",
	GameDataNotFound:"game data not found",
	PropNotFound:"prop not found",
	PropNotEnough: "prop not enough",
	RoleNotFound: "role not found",
	ComposeEquipError: "compose equip error",
	OpenidIsNil: "openid is nil",
	GumBallsGiftError: "gum balls gift draw error",
	UserNotFound: "user not found",
	FlagNotFound: "flag not found",
	FlagValueNotEnough: "flag value not enough",
	AlreadyDrawLetter: "already draw letter",
}


func GameExceptionCustom(api string, code ErrorCode, err error) {
	log.Error(api + " -- code:%v --- error:%v", code, err)
	GameException(code)
}
