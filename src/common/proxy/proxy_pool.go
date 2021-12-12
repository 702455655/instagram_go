package proxy

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"net/url"
)

type ProxyType int
type BlackType int

var (
	ProxyHttp   ProxyType = 0
	ProxySocket ProxyType = 1

	BlackType_NoBlack      BlackType = 0
	BlackType_Risk         BlackType = 1
	BlackType_Conn         BlackType = 2
	BlackType_RegisterRisk BlackType = 3
)

type Proxy struct {
	ID              string    `json:"id"`
	Ip              string    `json:"ip"`
	Port            string    `json:"port"`
	Username        string    `json:"username"`
	Passwd          string    `json:"passwd"`
	Rip             string    `json:"rip"`
	ProxyType       ProxyType `json:"proxy_type"`
	NeedAuth        bool      `json:"need_auth"`
	Country         string    `json:"country"`
	IsUsed          bool      `json:"is_used"`
	IsBusy          bool      `json:"is_busy"`
	RegisterSuccess int       `json:"register_success"`
	RegisterError   int       `json:"register_error"`
	BlackType       BlackType `json:"black_type"`
}

func (this *Proxy) GetProxy() *http.Transport {
	if this.ProxyType == 0 {
		var proxyUrl string
		if this.NeedAuth {
			proxyUrl = "http://" + this.Username + ":" + this.Passwd + "@" + this.Ip + ":" + this.Port
		} else {
			proxyUrl = "http://" + this.Ip + ":" + this.Port
		}
		_url, _ := url.Parse(proxyUrl)
		return &http.Transport{
			Proxy:           http.ProxyURL(_url),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		var auth *proxy.Auth = &proxy.Auth{}
		if this.NeedAuth {
			auth.User = this.Username
			auth.Password = this.Passwd
		} else {
			auth = nil
		}
		dialer, _ := proxy.SOCKS5("tcp", this.Ip+":"+this.Port, auth, proxy.Direct)
		var httpTran = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		httpTran.Dial = dialer.Dial
		return httpTran
	}
}

type ProxyPoolt interface {
	GetNoRisk(busy bool, used bool) *Proxy
	Get(id string) *Proxy
	Black(proxy *Proxy, _type BlackType)
	Remove(proxy *Proxy)
	Dumps()
}

var ProxyPool ProxyPoolt

type ProxyConfigt struct {
	Provider string `json:"provider"`
	Url      string `json:"url"`
	//ProxyType ProxyType `json:"proxy_type"`
}

var ProxyConfig ProxyConfigt

func InitProxyPool(configPath string) error {
	err := common.LoadJsonFile(configPath, &ProxyConfig)
	if err != nil {
		log.Error("load proxy config error: %v", err)
		return err
	}

	switch ProxyConfig.Provider {
	case "dove":
		ProxyPool, err = InitDovePool(ProxyConfig.Url)
		break
	case "luminati":
		ProxyPool, err = InitLuminatiPool(ProxyConfig.Url)
		break
	default:
		return &common.MakeMoneyError{ErrStr: fmt.Sprintf("proxy config provider error: %s",
			ProxyConfig.Provider), ErrType: common.OtherError}
	}
	return err
}
