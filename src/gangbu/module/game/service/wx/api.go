package wx

import (
	"aliens/log"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gangbu/constant"
	"gangbu/module/game/config"
	"gangbu/module/game/db"
	"gangbu/module/game/service/wx/model"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	PARAM_APPID = "appid"
	PARAM_SECRET = "secret"
	PARAM_CODE = "js_code"
	PARAM_GRANT_TYPE = "grant_type"
)



func WxLogin(code string) (openID string, sessionKey string,error error) {
	var url string
	//if config.Server.Platform == constant.PLATFORM_TOUTIAO {
	//	url = "https://developer.toutiao.com/api/apps/jscode2session?appid=" + config.Server.WxConfig.AppID +
	//		"&secret=" + config.Server.WxConfig.AppSecret +
	//		"&code=" + code
	//} else if config.Server.Platform == constant.PLATFORM_WECHAT {
		url = "https://api.weixin.qq.com/sns/jscode2session?appid=" + config.Server.WxConfig.AppID +
			"&secret=" + config.Server.WxConfig.AppSecret +
			"&grant_type=authorization_code" +
			"&js_code=" + code
	//} else {
	//	log.Fatal("platform config is nil")
	//}

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var openid string
	var session_key string
	if bytes.Contains(body, []byte("openid")) {
		atr := &model.LoginResponse{}
		err = json.Unmarshal(body, &atr)
		openid = atr.OpenID
		session_key = atr.SessionKey
		if err != nil {
			return "", "", err
		}
	} else {
		ater := &model.ErrorResponse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return "", "", err
		}
		return "", "", fmt.Errorf("%s", ater.ErrMsg)
	}
	return openid, session_key, nil
}

func RefreshToken() {
	//TODO 只需要一个服务器刷新即可
	//只有主节点才需要刷新token
	result, err := FetchAccessToken()
	if err != nil {
		log.Error("refresh accessToken error : %v", err)
		return
	}
	log.Info("refresh accessToken success : %v", result)
	//clustercache.Cluster.SetAccessToken(result.AccessToken, int(result.ExpiresIn))
}

//获取wx_AccessToken 拼接get请求 解析返回json结果 返回 AccessToken和err
func FetchAccessToken() (*model.AccessTokenResponse, error) {
	if wxConfig.AppID == "" {
		return nil, errors.New("accessToken param is not configuration")
	}
	//appID string, appSecret string, accessTokenFetchUrl string
	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=",
		wxConfig.AppID,
		"&secret=",
		wxConfig.AppSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(body, []byte("access_token")) {
		atr := &model.AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return nil, err
		}
		return atr, nil
	} else {
		ater := &model.ErrorResponse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", ater.ErrMsg)
	}
}

func BuildTextResponse(fromUserName, toUserName, content, nonce, timestamp string, encrypt bool) []byte {
	if encrypt {
		respBody, e := MakeEncryptResponseBody(fromUserName, toUserName, content, nonce, timestamp)
		if e != nil {
			log.Error("make encrypt response err : %v", e)
		}
		return respBody
	} else {
		respBody, e := MakeTextResponseBody(fromUserName, toUserName, content)
		if e != nil {
			log.Error("make text response err : %v", e)
		}
		return respBody
	}
}

//func HandleClickEvent(eventID string, openID string) string {
//	//激活游戏特权
//	if eventID == constant.EVENT_ACTIVE_PRIVILEGE {
//		_, err := UpdateUnionID(openID)
//		if err != nil {
//			log.Error("%v", err)
//		}
//	} else if eventID == constant.EVENT_DAY_GIFT {
//		uid := GetUIDByWechatServiceOpenID(openID)
//		if uid > 0 {
//			SetDayGiftCache(uid)
//		}
//	}
//	return eventID
//}

////通过openid获取unionid
//func UpdateUnionID(openID string) (string, error) {
//	accessToken := clustercache.Cluster.GetAccessToken()
//	if accessToken == ""  {
//		return "", errors.New("accessToken can not be found")
//	}
//	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/user/info?lang=zh_CN&access_token=", accessToken, "&openid=", openID}, "")
//	resp, err := http.Get(requestLine)
//	if err != nil || resp.StatusCode != http.StatusOK {
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	log.Debug("get %v unionID : %v",openID, string(body))
//	decryptedMapping := make(map[string]interface{})
//	json.Unmarshal(body, &decryptedMapping)
//	unionID, _ := decryptedMapping["unionid"].(string)
//	if unionID != "" {
//		//log.Debug("update wx service oumapping %v - %v", openID, unionID)
//		//cache.UserCache.SetWxServiceMapping(openID, unionID)
//		//uid := cache.UserCache.GetUserIDByUnionID(unionID)
//		//if uid > 0 {
//		//	rpc.UserServiceProxy.PersistCall(uid, &protocol.C2GS{
//		//		Sequence:[]int32{160},
//		//		ActivePrivilege:true,
//		//	})
//		//}
//	}
//	return unionID, nil
//}
func HandleClickEvent(eventID string, openID string) string {
	if eventID == constant.EVENT_DAY_GIFT {
		user := db.GetUserInfo("gb_" + openID)
		if user == nil {
			return eventID
		}
		value := constant.FLAG_CONFIG_VALUE[constant.FLAG_DAY_GIFT]
		flag := db.GetFlag(user.ID, constant.FLAG_SUBSCRIBE, value)
		if flag.UseValue() {
			db.AddItem(user.ID,constant.ITEM_PROP,constant.DAY_GIFT_ITEM_ID,constant.DAY_GIFT_NUM)
			db.UpdateOne(flag)
		}
	}
	return eventID
}

func HandleOtherEvent(event string, openID string) {
	user := db.GetUserInfo("gb_"+openID)
	if user == nil {
		return
	}
	subscribeFlag := db.GetFlag(user.ID, constant.FLAG_SUBSCRIBE, constant.FLAG_VALUE_TRUE)
	if event == constant.EVENT_SUBSCRIBE {
		subscribeFlag.SetValue(constant.FLAG_VALUE_TRUE)
		db.UpdateOne(subscribeFlag)
	} else if event == constant.EVENT_UNSUBSCRIBE {
		subscribeFlag.SetValue(constant.FLAG_VALUE_FALSE)
		db.UpdateOne(subscribeFlag)
	}
}
