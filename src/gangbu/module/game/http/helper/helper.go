package helper

import (
	"gangbu/exception"
	"gangbu/module/game/service/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

//type ErrorCode int32

//const (
//	ErrorCodeNone ErrorCode = iota
//	WxLogin     	= 10001	//微信登录出错
//	EroorCodeWxData         = 10002
//)

func ResponseWithCode(c *gin.Context, code exception.ErrorCode) {
	c.JSON(http.StatusOK, gin.H{"errcode":code,"errmsg":exception.ErrorMapping[code]})

}

func ResponseWithData(c *gin.Context, obj interface{}) {
	res := structToMap(obj)
	res["errcode"] = exception.CodeSuccess
	res["errmsg"] = exception.ErrorMapping[exception.CodeSuccess]
	c.JSON(http.StatusOK, res)
}

func structToMap(obj interface{}) gin.H {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Tag.Get("json")] = field.Interface()
	}
	return data
}

// 获取通过token解析出的用户信息字段
func GetClaimUser(c *gin.Context) *myjwt.CustomClaims {
	user := c.MustGet("claims")
	if user == nil {
		exception.GameException(exception.TokenInvalid)
	}
	return user.(*myjwt.CustomClaims)
}

func CheckReq (c *gin.Context, req interface{}) {
	if err := c.ShouldBind(req); err != nil {
		exception.GameExceptionCustom("CheckReq", exception.InvalidParam, err)
	}
}