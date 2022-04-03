package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxys"
	"makemoney/common/verification"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	ProxyPath       string                        `json:"proxy_path"`
	ResIcoPath      string                        `json:"res_ico_path"`
	ResUsernamePath string                        `json:"res_username_path"`
	Coro            int                           `json:"coro"`
	Country         string                        `json:"country"`
	ProviderName    string                        `json:"provider_name"`
	Gmail           *verification.GmailConfig     `json:"gmail"`
	Guerrilla       *verification.GuerrillaConfig `json:"guerrilla"`
	Taxin           *verification.PhoneInfo       `json:"taxin"`
}

var ConfigPath = flag.String("config", "./config/register.json", "")
var RegisterCount = flag.Int("count", 0, "")

var config Config

var Count int32 = 0
var SuccessCount = 0
var ErrorCreateCount = 0
var ErrorSendCodeCount = 0
var ErrorRecvCodeCount = 0
var ErrorCodeCount = 0
var ErrorCheckAccountCount = 0

var ErrorChallengeRequired = 0
var ErrorFeedback = 0
var ErrorOther = 0

var WaitAll sync.WaitGroup
var logTicker *time.Ticker

var PhoneProvider verification.VerificationCodeProvider
var Guerrilla verification.VerificationCodeProvider

func LogStatus() {
	for range logTicker.C {
		log.Info("success: %d, create err: %d, send err: %d, recv err: %d, challenge: %d, feedback: %d, check err: %d",
			SuccessCount,
			ErrorCreateCount,
			ErrorSendCodeCount,
			ErrorRecvCodeCount,
			ErrorChallengeRequired,
			ErrorFeedback,
			ErrorCheckAccountCount,
		)
	}
}

func statError(err error) {
	if common.IsError(err, common.ChallengeRequiredError) {
		ErrorChallengeRequired++
	} else if common.IsError(err, common.FeedbackError) {
		ErrorFeedback++
	} else {
		ErrorOther++
	}
}

func GenAddressBook() []goinsta.AddressBook {
	addr := make([]goinsta.AddressBook, common.GenNumber(20, 30))
	for index := range addr {
		addr[index].EmailAddresses = []string{common.GenString(common.CharSet_All, common.GenNumber(0, 10)) + "@gmail.com"}
		addr[index].PhoneNumbers = []string{"+1 " + "410 " + "895 " + common.GenString(common.CharSet_123, 4)}
		addr[index].LastName = common.GenString(common.CharSet_All, common.GenNumber(0, 10))
		addr[index].FirstName = common.GenString(common.CharSet_All, common.GenNumber(0, 10))
	}
	return addr[:]
}

func RegisterByPhone() {
	provider := PhoneProvider
	for true {
		var err error
		curCount := atomic.AddInt32(&Count, 1)
		if curCount > int32(*RegisterCount) {
			break
		}
		_proxy := proxys.ProxyPool.GetNoRisk(config.Country, true, true)
		if _proxy == nil {
			log.Error("get proxy error: %v", _proxy)
			break
		}

		inst := goinsta.New("", "", _proxy)
		inst.AccountInfo.Register.RegisterIpCountry = _proxy.Country
		prepare := inst.PrepareNewClient()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))

		username := common.Resource.ChoiceUsername()
		password := "XBYLxbyl1234"
		regisert := goinsta.Register{
			Inst:         inst,
			RegisterType: "phone",
			//Account:      account,
			Username: username,
			Password: password,
			AreaCode: provider.GetArea(),
			Year:     fmt.Sprintf("%d", common.GenNumber(1995, 2000)),
			Month:    fmt.Sprintf("%02d", common.GenNumber(1, 11)),
			Day:      fmt.Sprintf("%02d", common.GenNumber(1, 27)),
		}

		err = regisert.GetSignupConfig()
		err = regisert.GetCommonEmailDomains()
		err = regisert.PrecheckCloudId()
		err = regisert.IgUser()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))

		var account string
		account, err = provider.RequireAccount()
		if err != nil {
			log.Error("require account error: %v", err)
			break
		}
		regisert.Account = account
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(1000, 2000)))

		_, err = regisert.SendSignupSmsCode()
		if err != nil {
			ErrorSendCodeCount++
			statError(err)
			provider.ReleaseAccount(regisert.Account)
			log.Error("phone %s send error: %v", account, err)
			continue
		}
		code, err := provider.RequireCode(account)
		if err != nil {
			ErrorRecvCodeCount++
			statError(err)
			log.Error("phone %s require code error: %v", account, err)
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.ValidateSignupSmsCode(code)
		if err != nil {
			ErrorCodeCount++
			statError(err)
			log.Error("phone %s check code error: %v", account, err)
			continue
		}

		regisert.GenUsername()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CheckAgeEligibility()
		_, err = regisert.NewUserFlowBegins()

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CreatePhone()
		if err != nil {
			ErrorCreateCount++
			statError(err)
			log.Error("phone %s create error: %v", account, err)
			continue
		}

		prepareStr, _ := json.Marshal(prepare)
		_, err = regisert.GetSteps()

		if err == nil {
			log.Info("phone: %s register success! prepare: %s, account: %s password: %s", account, prepareStr, inst.User, inst.Pass)
			_ = goinsta.SaveInstToDB(inst)
			_, err = regisert.NewAccountNuxSeen()
			_, err = inst.AddressBookLink(GenAddressBook())
			var uploadID string
			uploadID, _, err = inst.GetUpload().UploadPhotoFromPath(common.Resource.ChoiceIco(), nil)
			err = inst.GetAccount().ChangeProfilePicture(uploadID)
			SuccessCount++
		} else {
			_ = goinsta.SaveInstToDB(inst)
			accInfo, _ := json.Marshal(inst.AccountInfo)
			log.Error("phone %s create error: %v, prepare: %s, account info: %s", account, err, prepareStr, accInfo)
			statError(err)
			ErrorCreateCount++

			//if common.IsError(err, common.ChallengeRequiredError) {
			//	log.Error("phone: %s had been challenge_required", account)
			//	continue
			//} else if common.IsError(err, common.FeedbackError) {
			//	ErrorCreateCount++
			//	log.Error("phone: %s had been feedback_required", account)
			//	continue
			//}
		}

	}
	WaitAll.Done()
}

func RegisterByEmail() {
	var mail verification.VerificationCodeProvider
	if config.ProviderName == "gmail" {
		mail = verification.GetGMails()
	} else {
		mail = Guerrilla
	}

	if mail == nil {
		log.Error("get mail error,so return!")
		return
	}

	for true {
		curCount := atomic.AddInt32(&Count, 1)
		if curCount > int32(*RegisterCount) {
			break
		}
		_proxy := proxys.ProxyPool.GetNoRisk(config.Country, true, true)
		if _proxy == nil {
			log.Error("get proxy error: %v", _proxy)
			break
		}
		var err error

		account, err := mail.RequireAccount()
		//account := "admin1@followmebsix.com"
		if err != nil {
			log.Error("require account error: %v", err)
			break
		}

		username := common.Resource.ChoiceUsername()
		password := common.GenString(common.CharSet_abc, 4) +
			common.GenString(common.CharSet_123, 4)
		//password := "xbyl1234"

		inst := goinsta.New("", "", _proxy)
		regisert := goinsta.Register{
			Inst:         inst,
			RegisterType: "email",
			Account:      account,
			Username:     username,
			Password:     password,
			Year:         fmt.Sprintf("%02d", common.GenNumber(1995, 2000)),
			Month:        fmt.Sprintf("%02d", common.GenNumber(1, 11)),
			Day:          fmt.Sprintf("%02d", common.GenNumber(1, 27)),
		}
		inst.AccountInfo.Register.RegisterIpCountry = _proxy.Country
		inst.PrepareNewClient()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))
		err = regisert.GetSignupConfig()

		err = regisert.GetCommonEmailDomains()
		err = regisert.PrecheckCloudId()
		err = regisert.IgUser()

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))
		_, err = regisert.CheckEmail()
		if err != nil {
			ErrorCheckAccountCount++
			statError(err)
			log.Error("email %s check error: %v", account, err)
			continue
		}

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(1000, 2000)))
		_, err = regisert.SendVerifyEmail()
		if err != nil {
			ErrorSendCodeCount++
			statError(err)
			log.Error("email %s send error: %v", account, err)
			continue
		}
		code, err := mail.RequireCode(account)
		if err != nil {
			ErrorRecvCodeCount++
			statError(err)
			log.Error("email %s require code error: %v", account, err)
			continue
		}

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CheckConfirmationCode(code)
		if err != nil {
			ErrorCodeCount++
			statError(err)
			log.Error("email %s check code error: %v", account, err)
			continue
		}

		regisert.GenUsername()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CheckAgeEligibility()
		_, err = regisert.NewUserFlowBegins()

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CreateEmail()
		if err != nil {
			ErrorCreateCount++
			statError(err)
			log.Error("email %s create error: %v", account, err)
			continue
		}
		_, err = regisert.GetSteps()

		if err == nil {
			log.Info("email: %s register success!   account: %s password: %s", account, inst.User, inst.Pass)
			_ = goinsta.SaveInstToDB(inst)
			_, err = regisert.NewAccountNuxSeen()
			_, err = inst.AddressBookLink(GenAddressBook())
			var uploadID string
			uploadID, _, err = inst.GetUpload().UploadPhotoFromPath(common.Resource.ChoiceIco(), nil)
			err = inst.GetAccount().ChangeProfilePicture(uploadID)
			SuccessCount++
		} else {
			_ = goinsta.SaveInstToDB(inst)
			accInfo, _ := json.Marshal(inst.AccountInfo)
			log.Error("email %s create error: %v,   account info: %s", account, err, accInfo)
			statError(err)
			ErrorCreateCount++
		}
	}
	WaitAll.Done()
}

func initParams() {
	flag.Parse()
	log.InitDefaultLog("register", true, true)
	err := common.LoadJsonFile(*ConfigPath, &config)
	if err != nil {
		log.Error("load config error: %v", err)
		os.Exit(0)
	}
	if config.ProxyPath == "" {
		log.Error("proxy path is null")
		os.Exit(0)
	}
	if config.ResIcoPath == "" {
		log.Error("ResourcePath is null")
		os.Exit(0)
	}
	if config.ResUsernamePath == "" {
		log.Error("ResUsernamePath is null")
		os.Exit(0)
	}
	if *RegisterCount == 0 {
		log.Error("RegisterCount is 0")
		os.Exit(0)
	}
}

//girlchina001
//a123456789
func main() {
	common.UseCharles = false
	common.UseTruncation = true
	initParams()
	routine.InitRoutine(config.ProxyPath)
	var err error
	switch config.ProviderName {
	case "taxin":
		PhoneProvider, err = verification.InitTaxin(config.Taxin)
		break
	case "gmail":
		err = verification.InitDefaultGMail(config.Gmail)
		break
	case "guerrilla":
		Guerrilla, err = verification.InitGuerrilla(config.Guerrilla)
		break
	}

	if err != nil {
		log.Error("create provider error! %v", err)
		os.Exit(0)
	}

	err = common.InitResource(config.ResIcoPath, config.ResUsernamePath)
	if err != nil {
		log.Error("InitResource error!%v", err)
		os.Exit(0)
	}

	WaitAll.Add(config.Coro)

	if config.ProviderName == "gmail" || config.ProviderName == "guerrilla" {
		for i := 0; i < config.Coro; i++ {
			go RegisterByEmail()
		}
	} else if config.ProviderName == "taxin" {
		for i := 0; i < config.Coro; i++ {
			go RegisterByPhone()
		}
	}

	logTicker = time.NewTicker(time.Second * 10)
	go LogStatus()
	WaitAll.Wait()
	logTicker.Stop()
	log.Info("success: %d, create err: %d, send err: %d, recv err: %d, challenge: %d, feedback: %d, check err: %d",
		SuccessCount,
		ErrorCreateCount,
		ErrorSendCodeCount,
		ErrorRecvCodeCount,
		ErrorChallengeRequired,
		ErrorFeedback,
		ErrorCheckAccountCount)
	log.Info("task finish!")
}
