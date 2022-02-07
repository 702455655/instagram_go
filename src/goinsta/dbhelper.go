package goinsta

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxy"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strings"
	"time"
)

type MogoDBHelper struct {
	Client         *mongo.Client
	Phone          *mongo.Collection
	Account        *mongo.Collection
	UploadIDRecord *mongo.Collection
}

var MogoHelper *MogoDBHelper = nil

func InitMogoDB(mogoUri string) {
	//"mongodb://xbyl:XBYLxbyl1234@62.216.92.183:27017"
	clientOptions := options.Client().ApplyURI(mogoUri)
	//clientOptions := options.Client().ApplyURI("mongodb://xbyl:xbyl741852JHK@192.168.187.1:27017")

	var err error
	MogoHelper = &MogoDBHelper{}

	MogoHelper.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("mongo %v", err)
	}

	err = MogoHelper.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error("mongo %v", err)
	}

	MogoHelper.Phone = MogoHelper.Client.Database("inst").Collection("phone")
	MogoHelper.Account = MogoHelper.Client.Database("inst").Collection("account")
	MogoHelper.UploadIDRecord = MogoHelper.Client.Database("inst").Collection("upload_id")
}

func GetDB(name string) *mongo.Database {
	return MogoHelper.Client.Database(name)
}

type PhoneStorage struct {
	Area          string        `bson:"area"`
	Phone         string        `bson:"phone"`
	SendCount     string        `bson:"send_count"`
	RegisterCount int           `bson:"register_count"`
	Provider      string        `bson:"provider"`
	LastUseTime   time.Duration `bson:"last_use_time"`
}

func UpdatePhoneSendOnce(provider string, area string, number string) error {
	_, err := MogoHelper.Phone.UpdateOne(context.TODO(),
		bson.D{
			{"area", area},
			{"phone", number},
		}, bson.D{{"$set", bson.D{{"area", area},
			{"phone", number},
			{"provider", provider},
			{"last_use_time", time.Now()},
		}}, {"$inc", bson.M{"send_count": 1}},
		}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func UpdatePhoneRegisterOnce(area string, number string) error {
	_, err := MogoHelper.Phone.UpdateOne(context.TODO(),
		bson.D{
			{"area", area},
			{"phone", number},
		}, bson.D{{"$inc", bson.M{"register_count": 1}}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

type AccountCookies struct {
	ID                  int64                    `json:"id" bson:"id"`
	Username            string                   `json:"username" bson:"username"`
	Passwd              string                   `json:"passwd" bson:"passwd"`
	HttpHeader          map[string]string        `json:"http_header" bson:"http_header"`
	ProxyID             string                   `json:"proxy_id" bson:"proxy_id"`
	IsLogin             bool                     `json:"is_login" bson:"is_login"`
	Token               string                   `json:"token" bson:"token"`
	Cookies             []*http.Cookie           `json:"cookies" bson:"cookies"`
	CookiesB            []*http.Cookie           `json:"cookies_b" bson:"cookies_b"`
	Device              *InstDeviceInfo          `json:"device" bson:"device"`
	RegisterEmail       string                   `json:"register_email" bson:"register_email"`
	RegisterPhoneNumber string                   `json:"register_phone_number" bson:"register_phone_number"`
	RegisterPhoneArea   string                   `json:"register_phone_area" bson:"register_phone_area"`
	RegisterIpCountry   string                   `json:"register_ip_country" bson:"register_ip_country"`
	RegisterTime        int64                    `json:"register_time" bson:"register_time"`
	Status              string                   `json:"status" bson:"status"`
	LastSendMsgTime     int                      `json:"last_send_msg_time" bson:"last_send_msg_time"`
	Tag                 string                   `json:"tag" bson:"tag"`
	SpeedControl        map[string]*SpeedControl `json:"speed_control" bson:"speed_control"`
}

func SaveNewAccount(account AccountCookies) error {
	_, err := MogoHelper.Account.UpdateOne(
		context.TODO(),
		bson.M{"username": account.Username},
		bson.M{"$set": account},
		options.Update().SetUpsert(true))
	return err
}

func LoadDBAccountByTags(tag string) ([]AccountCookies, error) {
	cursor, err := MogoHelper.Account.Find(context.TODO(), bson.M{"tags": tag}, nil)
	if err != nil {
		return nil, err
	}
	var ret []AccountCookies
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func LoadDBAllAccount() ([]AccountCookies, error) {
	cursor, err := MogoHelper.Account.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var ret []AccountCookies
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CleanStatus() error {
	_, err := MogoHelper.Account.UpdateMany(
		context.TODO(),
		bson.M{},
		bson.M{"$set": bson.M{"status": ""}},
		options.Update().SetUpsert(true))
	return err
}

func SaveInstToDB(inst *Instagram) error {
	url, _ := neturl.Parse(InstagramHost)
	urlb, _ := neturl.Parse(InstagramHost_B)

	Cookies := AccountCookies{
		ID:                  inst.ID,
		Username:            inst.User,
		Passwd:              inst.Pass,
		Token:               inst.token,
		Device:              inst.Device,
		Cookies:             inst.c.Jar.Cookies(url),
		CookiesB:            inst.c.Jar.Cookies(urlb),
		HttpHeader:          inst.httpHeader,
		ProxyID:             inst.Proxy.ID,
		IsLogin:             inst.IsLogin,
		RegisterEmail:       inst.RegisterEmail,
		RegisterPhoneNumber: inst.RegisterPhoneNumber,
		RegisterPhoneArea:   inst.RegisterPhoneArea,
		RegisterIpCountry:   inst.RegisterIpCountry,
		RegisterTime:        inst.RegisterTime,
		Status:              inst.Status,
		LastSendMsgTime:     inst.LastSendMsgTime,
		SpeedControl:        inst.SpeedControl,
	}
	return SaveNewAccount(Cookies)
}

func ConvConfig(config *AccountCookies) (*Instagram, error) {
	url, err := neturl.Parse(InstagramHost)
	if err != nil {
		return nil, err
	}
	urlb, err := neturl.Parse(InstagramHost_B)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(url, config.Cookies)
	jar.SetCookies(urlb, config.CookiesB)

	inst := &Instagram{
		ID:                  config.ID,
		User:                config.Username,
		Pass:                config.Passwd,
		token:               config.Token,
		Device:              config.Device,
		httpHeader:          config.HttpHeader,
		IsLogin:             config.IsLogin,
		RegisterEmail:       config.RegisterEmail,
		RegisterPhoneNumber: config.RegisterPhoneNumber,
		RegisterPhoneArea:   config.RegisterPhoneArea,
		RegisterIpCountry:   config.RegisterIpCountry,
		Status:              config.Status,
		SpeedControl:        config.SpeedControl,
		sessionID:           strings.ToUpper(common.GenUUID()),
		LastSendMsgTime:     config.LastSendMsgTime,
		RegisterTime:        config.RegisterTime,
		c: &http.Client{
			Jar: jar,
		},
	}

	if inst.Device == nil {
		inst.Device = GenInstDeviceInfo()
	}
	if inst.SpeedControl == nil {
		inst.SpeedControl = make(map[string]*SpeedControl)
	}
	for key, value := range inst.SpeedControl {
		ReSetRate(value, key)
	}

	inst.graph = &Graph{inst: inst}
	inst.Proxy = &proxy.Proxy{ID: config.ProxyID}
	common.DebugHttpClient(inst.c)

	return inst, nil
}

func LoadAccountByTags(tag string) []*Instagram {
	config, err := LoadDBAccountByTags(tag)
	if err != nil {
		return nil
	}
	var ret []*Instagram
	for item := range config {
		inst, err := ConvConfig(&config[item])
		if err != nil {
			log.Warn("conv config to inst error:%v", err)
			continue
		}
		ret = append(ret, inst)
	}
	return ret
}

func LoadAllAccount() []*Instagram {
	config, err := LoadDBAllAccount()
	if err != nil {
		return nil
	}
	var ret []*Instagram
	for item := range config {
		inst, err := ConvConfig(&config[item])
		if err != nil {
			log.Warn("conv config to inst error:%v", err)
			continue
		}
		ret = append(ret, inst)
	}
	return ret
}

type UploadIDRecord struct {
	FileMd5  string `bson:"file_md5"`
	Username string `bson:"username"`
	FileType string `bson:"file_type"`
	FileName string `bson:"file_name"`
	UploadID string `bson:"upload_id"`
}

func SaveUploadID(record *UploadIDRecord) error {
	_, err := MogoHelper.UploadIDRecord.UpdateOne(
		context.TODO(),
		bson.M{"file_md5": record.FileMd5},
		bson.M{"$set": record},
		options.Update().SetUpsert(true))
	return err
}

func FindUploadID(username string, fileMd5 string) (*UploadIDRecord, error) {
	cursor, err := MogoHelper.UploadIDRecord.Find(context.TODO(),
		bson.D{{"$and",
			bson.D{
				{"username", username},
				{"file_md5", fileMd5},
			},
		}}, nil)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var result = &UploadIDRecord{}
	if cursor.Next(context.TODO()) {
		err = cursor.Decode(result)
		return result, err
	}

	return nil, &common.MakeMoneyError{ErrType: common.NoMoreError}
}
