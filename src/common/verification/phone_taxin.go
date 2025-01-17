package verification

import (
	"makemoney/common"
	"makemoney/common/log"
	"strings"
	"time"
)

//http://h5.do889.com:81/info
//741852
type PhoneTaxin struct {
	*PhoneInfo
}

type BaseRespPhoneTaxin struct {
	Stat    bool   `json:"stat"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type PhoneTaxinLogin struct {
	BaseRespPhoneTaxin
	Data struct {
		Money string `json:"money"`
		Cash  string `json:"cash"`
		Token string `json:"token"`
	} `json:"data"`
}

func (this *PhoneTaxin) Login() error {
	var respJson PhoneTaxinLogin
	err := common.HttpDoJson(this.Client, &common.RequestOpt{
		ReqUrl: this.UrlLogin,
		Params: map[string]string{
			"username": this.Username,
			"password": this.Password,
			"type":     "json",
		}}, &respJson)
	if err != nil {
		return err
	}
	this.Token = respJson.Data.Token

	log.Info("phone token %s", respJson.Data.Token)
	return nil
}

type PhoneTaxinRequirephone struct {
	BaseRespPhoneTaxin
	Data string `json:"data"`
}

func (this *PhoneTaxin) RequireAccount() (string, error) {
	var respJson PhoneTaxinRequirephone
	err := common.HttpDoJson(this.Client, &common.RequestOpt{
		ReqUrl: this.UrlReqPhoneNumber,
		Params: map[string]string{
			"token": this.Token,
			"id":    this.ProjectID,
			//"area":  this.City,
			"loop": "1",
			"card": "1",
			"type": "json",
		}}, &respJson)

	if err != nil {
		return "", err
	}
	if respJson.Message != "ok" {
		return "", &common.MakeMoneyError{ErrStr: respJson.Message, ErrType: common.RequirePhoneError}
	}

	return respJson.Data, err
}

func (this *PhoneTaxin) RequireCode(number string) (string, error) {
	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		resp, err := common.HttpDo(this.Client, &common.RequestOpt{ReqUrl: this.UrlReqPhoneCode,
			Params: map[string]string{
				"token": this.Token,
				"id":    this.ProjectID,
				"phone": number,
			}})
		sp := strings.Split(resp, "|")
		if len(sp) != 2 || err != nil {
			log.Warn("to getting phone %s code request error: %v", number, err)
		} else if sp[0] == "0" {
			log.Warn("to getting phone %s code error: %v", number, resp)
		} else if sp[0] == "1" {
			code := common.GetCode(sp[1])
			if code != "" {
				return code, nil
			} else {
				log.Warn("to getting phone %s code parse error", number)
			}
		} else {
			log.Warn("to getting phone %s code error: %v", number, resp)
		}
		time.Sleep(this.RetryDelayDuration)
	}

	return "", &common.MakeMoneyError{ErrStr: "require code timeout", ErrType: common.RecvPhoneCodeError}
}

type PhoneTaxinReleasePhone struct {
	BaseRespPhoneTaxin
}

func (this *PhoneTaxin) ReleaseAccount(number string) error {
	var respJson PhoneTaxinReleasePhone
	err := common.HttpDoJson(this.Client, &common.RequestOpt{
		ReqUrl: this.UrlReqReleasePhone,
		Params: map[string]string{
			"token": this.Token,
			"phone": number,
			"id":    this.ProjectID,
			"type":  "json",
		}}, &respJson)

	if err != nil {
		return err
	}
	if respJson.Message != "ok" {
		return &common.MakeMoneyError{ErrStr: respJson.Message}
	}
	return nil
}

type PhoneTaxinBlackPhone struct {
	BaseRespPhoneTaxin
}

func (this *PhoneTaxin) BlackAccount(number string) error {
	var respJson PhoneTaxinReleasePhone
	err := common.HttpDoJson(this.Client, &common.RequestOpt{ReqUrl: this.UrlReqBlackPhone,
		Params: map[string]string{
			"token": this.Token,
			"id":    this.ProjectID,
			"phone": number,
			"type":  "json",
		}}, &respJson)

	if err != nil {
		return err
	}
	if respJson.Message != "ok" {
		return &common.MakeMoneyError{ErrStr: respJson.Message}
	}
	return nil
}

type PhoneTaxinBalance struct {
	BaseRespPhoneTaxin
	Data struct {
		Money string `json:"money"`
		Cash  string `json:"cash"`
	} `json:"data"`
}

func (this *PhoneTaxin) GetBalance() (string, error) {
	var respJson PhoneTaxinBalance
	err := common.HttpDoJson(this.Client, &common.RequestOpt{ReqUrl: this.UrlReqBalance,
		Params: map[string]string{
			"token": this.Token,
			"type":  "json",
		}}, &respJson)

	if err != nil {
		return "", err
	}
	if respJson.Message != "ok" {
		return "", &common.MakeMoneyError{ErrStr: respJson.Message}
	}
	return respJson.Data.Money, nil
}
