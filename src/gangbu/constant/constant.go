package constant

const (
	ENERGY_INIT			= 6
	ENERGY_LIMIT		= 6
	ENERGY_COST			= 1
	ENERGY_RESTORE_TIME	= 60*60 //s
	ENERGY_RESTORE_NUM	= 1
	ENERGY_AD_RESTORE	= 1

	WATCH_AD_INTERVAL  = 60*60 //s
	MAX_DAY_AD_RESTORE = 6     //每天看广告恢复体力的次数

	FLOOR_FIRST			= 1
	MAX_FLOOR_SCORE		= 100000

	PROP_PER_COST = 1 //每次道具消耗的数量

	RANK_MAX_LIMIT = 5

	MAX_DAY_SHARE_HELP_NUM = 3 //每天最多的分享获得帮助的次数
	RELIVE_ITEM_ID = 3 //复活药剂的id
	HELP_ADD_NUM   = 1 //帮助加的复活药剂的数量

	DAY_GIFT_ITEM_ID = 3 //每日礼包复活药水
	DAY_GIFT_NUM = 3	 //每日礼包药水的数量
)

const (
	PROP_TYPE_HP		= 	1 //血
	PROP_TYPE_SHIELD	=	2 //盾
	PROP_TYPE_RELIVE	=	3 //复活
	PROP_TYPE_CLEAN		=	4 //清除debuff
)

const (
	CHANNEL_WECHAT	=	1
	CHANNEL_WEB		=	2
)

const (
	KEYWORD_EVENT = "event"
	EVENT_SUBSCRIBE = "subscribe"
	EVENT_UNSUBSCRIBE = "unsubscribe"
	EVENT_CLICK = "CLICK"

	EVENT_DAY_GIFT = "DAY_GIFT"
)

const (
	PLATFORM_WECHAT = "wx"
	PLATFORM_TOUTIAO = "tt"
)

const (
	ITEM_PROP 		= 1	//道具(药水) (both)
	ITEM_RELIC 		= 2 //碎片 (冈布奥)
	ITEM_EQUIPMENT 	= 3 //装备 (冈布奥)
	ITEM_ARCHIVE 	= 4 //伴手礼 (赛几)
)

const (
	DB_COMMAND_INSERT = "I"
	DB_COMMAND_UPDATE = "U"
	DB_COMMAND_DELETE = "D"
	DB_COMMAND_FUPDATE = "FU"
	DB_COMMAND_CONDITION_UPDATE = "CU"
	DB_COMMAND_CONDITION_DELETE = "CD"

	DB_DEBUG                 = true	 //是否开启数据库操作
	DB_SIGAL_TIMEOUT float64 = 1 //数据库操作超时告警阈值 1秒
)

type ROLE_FLAG int32

const (
	FLAG_NONE         ROLE_FLAG = iota
	FLAG_RELIVE_SHARE   = 1 //分享立即复活
	FLAG_DAY_GIFT		= 2
	FLAG_RELIVE_PROP_SHARE = 3 //分享获得复活道具
	FLAG_GUIDE			= 4

	FLAG_DRAW_LETTER = 21
	FLAG_SUBSCRIBE   = 22
	FLAG_BEAT_PETER  = 23 //打败peter怪

	FLAG_VALUE_FALSE = 0
	FLAG_VALUE_TRUE  = 1
)

var FLAG_CONFIG_VALUE = map[ROLE_FLAG]int32 {   //flag - count
	FLAG_RELIVE_SHARE: 1,
	FLAG_DAY_GIFT: 1,
	FLAG_RELIVE_PROP_SHARE: 1,

}

var FLAG_NEED_REFRESH = map[ROLE_FLAG]struct{} {
	FLAG_RELIVE_SHARE: {},
	FLAG_DAY_GIFT: {},
	FLAG_RELIVE_PROP_SHARE: {},
}

var CONFIG_LETTER_GIFT = map[int32]int32{ //relicId - num
	141:1, 161:1, 221:1, 241:1,
}