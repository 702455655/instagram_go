package routine

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxy"
	"makemoney/goinsta"
	"os"
)

type DataBaseConfig struct {
	MogoUri string `json:"mogo_uri"`
}

var dbConfig DataBaseConfig

func InitRoutine(proxyPath string) {
	err := common.LoadJsonFile("./config/dbconfig.json", &dbConfig)
	if err != nil {
		log.Error("load db config error:%v", err)
		panic(err)
	}
	common.InitMogoDB(dbConfig.MogoUri)
	err = proxy.InitProxyPool(proxyPath)
	if err != nil {
		log.Error("init ProxyPool error:%v", err)
		panic(err)
	}
	intas := goinsta.LoadAllAccount()
	if len(intas) == 0 {
		log.Error("there have no account!")
		os.Exit(0)
	}
	goinsta.InitAccountPool(intas)

	log.Info("load account count: %d", goinsta.AccountPool.Available.Len())
	//common.InitResource("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture", "C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\user_nameraw.txt")
}

func ReqAccount() (*goinsta.Instagram, error) {
	inst := goinsta.AccountPool.GetOne()
	if inst == nil {
		return nil, &common.MakeMoneyError{ErrType: common.NoMoreError, ErrStr: "no more account"}
	}

	if !SetProxy(inst) {
		return nil, &common.MakeMoneyError{ErrType: common.NoMoreError, ErrStr: "no more proxy"}
	}

	return inst, nil
}

func SetProxy(inst *goinsta.Instagram) bool {
	var _proxy *proxy.Proxy
	if inst.Proxy.ID != "" {
		_proxy = proxy.ProxyPool.Get(inst.Proxy.ID)
		if _proxy == nil {
			log.Warn("find insta proxy %s error!", inst.Proxy.ID)
		}
	}

	if _proxy == nil {
		_proxy = proxy.ProxyPool.GetNoRisk(false, false)
		if _proxy == nil {
			log.Error("get insta proxy error!")
		}
	}

	if _proxy != nil {
		inst.SetProxy(_proxy)
	} else {
		return false
	}
	return true
}
