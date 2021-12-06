package goinsta

import (
	"encoding/json"
	"makemoney/common"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

var (
	InsAccountError_ChallengeRequired = "challenge_required"
)

type Instagram struct {
	User                string
	Pass                string
	androidID           string
	uuid                string
	token               string
	familyID            string
	adid                string
	wid                 string
	challengeURL        string
	ID                  int64
	httpHeader          map[string]string
	RegisterPhoneNumber string
	RegisterPhoneArea   string
	RegisterIpCountry   string
	IsLogin             bool
	UserAgent           string
	Status              string
	sessionID           string

	ReqSuccessCount  int
	ReqErrorCount    int
	ReqApiErrorCount int
	ReqContError     int

	Proxy *common.Proxy
	c     *http.Client
}

func (this *Instagram) SetCookieJar(jar http.CookieJar) error {
	url, err := neturl.Parse(goInstaHost)
	if err != nil {
		return err
	}
	// First grab the cookies from the existing jar and we'll put it in the new jar.
	cookies := this.c.Jar.Cookies(url)
	this.c.Jar = jar
	this.c.Jar.SetCookies(url, cookies)
	return nil
}

func New(username, password string, _proxy *common.Proxy) *Instagram {
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		User:      username,
		Pass:      password,
		androidID: generateDeviceID(),
		uuid:      common.GenUUID(), // both uuid must be differents
		familyID:  common.GenUUID(),
		wid:       common.GenUUID(),
		adid:      common.GenUUID(),
		UserAgent: GenUserAgent(),
		sessionID: common.GenUUID(),
		c: &http.Client{
			Jar:       jar,
			Transport: _proxy.GetProxy(),
		},
	}
	inst.Proxy = _proxy
	inst.httpHeader = make(map[string]string)
	common.DebugHttpClient(inst.c)

	return inst
}

func (this *Instagram) GetSearch(q string) *Search {
	return newSearch(this, q)
}

func (this *Instagram) GetUpload() *Upload {
	return newUpload(this)
}

func (this *Instagram) GetAccount() *Account {
	return &Account{ID: this.ID, inst: this}
}

func (this *Instagram) GetUser(id string) *User {
	pk, _ := strconv.ParseInt(id, 10, 64)
	return &User{ID: pk, inst: this}
}

func (this *Instagram) GetMessage(msgType MessageType) *Message {
	return &Message{inst: this, msgType: msgType}
}

// SetProxy sets proxy for connection.
func (this *Instagram) SetProxy(_proxy *common.Proxy) {
	this.Proxy = _proxy
	this.c.Transport = _proxy.GetProxy()
	common.DebugHttpClient(this.c)
}

func (this *Instagram) NeedReplace() bool {
	if this.ReqContError >= 3 {
		return true
	}
	return false
}

func (this *Instagram) CleanCookiesAndHeader() {
	this.httpHeader = make(map[string]string)
	this.c.Jar, _ = cookiejar.New(nil)
}

func (this *Instagram) ReadHeader(key string) string {
	return this.httpHeader[key]
}

func (this *Instagram) PrepareNewClient() {
	//facebook_dod
	_ = this.ZrToken()
	_ = this.ZrToken()
	_ = this.LogAttribution()
	_ = this.QeSync()
	_ = this.launcherSync()
	_ = this.PrefillCandidates()
	_ = this.LoggingClientEvents()
}

//func (this *Instagram) readMsisdnHeader() error {
//	_, err := this.HttpRequest(
//		&reqOptions{
//			ApiPath: urlMsisdnHeader,
//			IsPost:  true,
//			Query: map[string]interface{}{
//				"device_id": this.uuid,
//			},
//		},
//	)
//	return err
//}

func (this *Instagram) ZrToken() error {
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlZrToken,
			IsPost:  false,
			IsApiB:  true,
			Query: map[string]interface{}{
				"device_id":        this.androidID,
				"token_hash":       "",
				"custom_device_id": this.uuid,
				"fetch_reason":     "token_expired",
			},
		},
	)
	return err
}

//早于注册登录?
func (this *Instagram) LogAttribution() error {
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlLogAttribution,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query: map[string]interface{}{
				"adid": this.adid,
			},
		},
	)
	return err
}

func (this *Instagram) QeSync() error {
	var params map[string]interface{}
	if this.IsLogin {
		params = map[string]interface{}{
			"id":          this.ID,
			"_uid":        this.ID,
			"_uuid":       this.uuid,
			"experiments": goInstaExperiments,
		}
	} else {
		params = map[string]interface{}{
			"id":                      this.uuid,
			"experiments":             goInstaExperiments,
			"server_config_retrieval": "1",
		}
	}
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlQeSync,
			Query:   params,
			IsPost:  true,
			Signed:  true,
			IsApiB:  true,
		},
	)
	return err
}

func (this *Instagram) launcherSync() error {
	var query map[string]interface{}
	var isApiB bool
	if this.IsLogin {
		query = map[string]interface{}{
			"id":                      this.ID,
			"_uid":                    this.ID,
			"_uuid":                   this.uuid,
			"server_config_retrieval": "1",
		}
		isApiB = false
	} else {
		query = map[string]interface{}{
			"id":                      this.uuid,
			"server_config_retrieval": "1",
		}
		isApiB = true
	}

	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlLauncherSync,
			IsPost:  true,
			IsApiB:  isApiB,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}

func (this *Instagram) PrefillCandidates() error {
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlPrefillCandidates,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query: map[string]interface{}{
				"android_device_id": this.androidID,
				"phone_id":          this.familyID,
				"usages":            "[\"account_recovery_omnibox\"]",
				"device_id":         this.uuid,
			},
		},
	)
	return err
}

func (this *Instagram) LoggingClientEvents() error {
	type DataItem struct {
		Name         string `json:"name"`
		Time         string `json:"time"`
		SamplingRate int    `json:"sampling_rate"`
		Extra        struct {
			CurrentVersion int    `json:"current_version"`
			Pk             string `json:"pk"`
			ReleaseChannel string `json:"release_channel"`
			RadioType      string `json:"radio_type"`
		} `json:"extra"`
	}
	data := make([]DataItem, 2)
	data[0].Name = "ig_emergency_push_did_set_initial_version"
	data[0].Time = strconv.FormatInt(time.Now().Unix(), 10) + ".021"
	data[0].SamplingRate = 1
	data[0].Extra.Pk = "0"
	data[0].Extra.CurrentVersion = 48
	data[0].Extra.ReleaseChannel = "prod"
	data[0].Extra.RadioType = "wifi-none"
	data[1] = data[0]

	params := map[string]interface{}{
		"seq":              1,
		"app_id":           "567067343352427",
		"app_ver":          goInstaVersion,
		"build_num":        goInstaBuildNum,
		"device_id":        this.uuid,
		"family_device_id": this.familyID,
		"session_id":       this.sessionID,
		"channel":          "zero_latency",
		"app_uid":          "0",
		"claims":           "[\"0\"]",
		"config_version":   "v2",
		"config_checksum":  "null",
		"data":             data,
		"log_type":         "client_event",
	}
	tmp, _ := json.Marshal(params)
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath:    urlLoggingClientEvents,
			IsPost:     true,
			IsApiGraph: true,
			Query: map[string]interface{}{
				"access_token": "567067343352427|f249176f09e26ce54212b472dbab8fa8",
				"format":       "json",
				"compressed":   "0",
				"sent_time":    strconv.FormatInt(time.Now().Unix(), 10) + ".021",
				"message":      tmp,
			},
		},
	)
	return err
}

type RespLogin struct {
	BaseApiResp
	LoggedInUser struct {
		AccountBadges                  []interface{} `json:"account_badges"`
		AccountType                    int           `json:"account_type"`
		AllowContactsSync              bool          `json:"allow_contacts_sync"`
		AllowedCommenterType           string        `json:"allowed_commenter_type"`
		BizUserInboxState              int           `json:"biz_user_inbox_state"`
		CanBoostPost                   bool          `json:"can_boost_post"`
		CanSeeOrganicInsights          bool          `json:"can_see_organic_insights"`
		CanSeePrimaryCountryInSettings bool          `json:"can_see_primary_country_in_settings"`
		CountryCode                    int           `json:"country_code"`
		FbidV2                         int64         `json:"fbid_v_2"`
		FollowFrictionType             int           `json:"follow_friction_type"`
		FullName                       string        `json:"full_name"`
		HasAnonymousProfilePicture     bool          `json:"has_anonymous_profile_picture"`
		HasPlacedOrders                bool          `json:"has_placed_orders"`
		InteropMessagingUserFbid       int64         `json:"interop_messaging_user_fbid"`
		IsBusiness                     bool          `json:"is_business"`
		IsCallToActionEnabled          interface{}   `json:"is_call_to_action_enabled"`
		IsPrivate                      bool          `json:"is_private"`
		IsUsingUnifiedInboxForDirect   bool          `json:"is_using_unified_inbox_for_direct"`
		IsVerified                     bool          `json:"is_verified"`
		Nametag                        struct {
			Emoji         string `json:"emoji"`
			Gradient      int    `json:"gradient"`
			Mode          int    `json:"mode"`
			SelfieSticker int    `json:"selfie_sticker"`
		} `json:"nametag"`
		NationalNumber                             int64  `json:"national_number"`
		PhoneNumber                                string `json:"phone_number"`
		Pk                                         int64  `json:"pk"`
		ProfessionalConversionSuggestedAccountType int    `json:"professional_conversion_suggested_account_type"`
		ProfilePicUrl                              string `json:"profile_pic_url"`
		ReelAutoArchive                            string `json:"reel_auto_archive"`
		ShowInsightsTerms                          bool   `json:"show_insights_terms"`
		TotalIgtvVideos                            int    `json:"total_igtv_videos"`
		Username                                   string `json:"username"`
		WaAddressable                              bool   `json:"wa_addressable"`
		WaEligibility                              int    `json:"wa_eligibility"`
	} `json:"logged_in_user"`
	SessionFlushNonce interface{} `json:"session_flush_nonce"`
}

func (this *Instagram) Login() error {
	encodePasswd, _ := encryptPassword(this.Pass, this.ReadHeader(IGHeader_EncryptionId), this.ReadHeader(IGHeader_EncryptionKey))
	params := map[string]interface{}{
		"jazoest":             genJazoest(this.familyID),
		"country_codes":       "[{\"country_code\":\"" + strings.ReplaceAll(this.RegisterPhoneArea, "+", "") + "\",\"source\":[\"default\"]}]",
		"phone_id":            this.familyID,
		"enc_password":        encodePasswd,
		"username":            this.User,
		"adid":                this.adid,
		"guid":                this.uuid,
		"device_id":           this.androidID,
		"google_tokens":       "[]",
		"login_attempt_count": "0",
	}
	resp := &RespLogin{}
	err := this.HttpRequestJson(&reqOptions{
		ApiPath: urlLogin,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	if err != nil && this.ReadHeader(IGHeader_Authorization) != "" {
		this.IsLogin = true
		this.ID = resp.LoggedInUser.Pk
	}
	return err
}

//注册成功后触发
func (this *Instagram) contactPrefill() error {
	var query map[string]interface{}

	if this.IsLogin {
		query = map[string]interface{}{
			"_uid":      this.ID,
			"device_id": this.uuid,
			"_uuid":     this.uuid,
			"usage":     "auto_confirmation",
		}
	} else {
		query = map[string]interface{}{
			"phone_id": this.familyID,
			"usage":    "prefill",
		}
	}

	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlContactPrefill,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}
