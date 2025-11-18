package dal

import (
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/es"
	"judgeMore/biz/dal/mysql"
)

func Init() {
	mysql.Init()
	cache.Init()
	es.Init()
}
