package words

import (
	"encoding/csv"
	"gangbu/module/game/config"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

var firstName []string = []string{}
var lastName []string = []string{}

//敏感词汇
var filters []string = []string{}

var regExp *regexp.Regexp

//func ValidateNickname(nickname string) {
//	if !regExp.MatchString(nickname) {
//		exception.GameException(exception.NicknameInvalidWord)
//	}
//
//	for _, filter := range filters {
//		if strings.Contains(nickname, filter) {
//			//log.Debug("%v", filter)
//			exception.GameException(exception.NicknameSensitiveWord)
//		}
//	}
//}

//随机姓名
func RandomName() string {
	//rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().UnixNano())
	return firstName[rand.Intn(len(firstName))] + lastName[rand.Intn(len(lastName))]
}

func Init() {
	regexp_, _ := regexp.Compile("^[\u4e00-\u9fa5a-zA-Z0-9_]{2,8}$")
	regExp = regexp_

	rfile, _ := os.Open(config.GetConfigPath("/data/name.csv"))
	r := csv.NewReader(rfile)
	for {
		strs, err := r.Read()
		if err != nil {
			break
		}
		firstName = append(firstName, strs[0])
		lastName = append(lastName, strs[1])
	}
	//LoadSensitiveConf("/sensitiveEx.csv")
	//LoadSensitiveConf("/sensitive.csv")
}


func LoadSensitiveConf(path string){
	rSensitive, _ := os.Open(config.GetConfigPath(path))
	r := csv.NewReader(rSensitive)
	for {
		strs, err := r.Read()
		if err != nil {
			break
		}
		for _, filter := range strs {
			filter = strings.TrimSpace(filter)
			if filter != "" {
				filters = append(filters, filter)
			}
		}
	}
}