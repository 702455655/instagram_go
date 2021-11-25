package goinsta

import (
	"makemoney/common"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"strings"
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
	id                  string
	httpHeader          map[string]string
	registerPhoneNumber string
	registerPhoneArea   string
	registerIpCountry   string
	IsLogin             bool

	ReqSuccessCount  int
	ReqErrorCount    int
	ReqApiErrorCount int

	Proxy *common.Proxy

	//Challenge *Challenge
	//Profiles *Profiles
	//Account *Account
	//Timeline *Timeline
	//Activity *Activity
	//Inbox *Inbox
	//Feed *Feed
	//Locations *LocationInstance

	c *http.Client
}

func (this *Instagram) SetCookieJar(jar http.CookieJar) error {
	url, err := neturl.Parse(goInstaAPIUrl)
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
	if !this.IsLogin {
		return nil
	}
	return newSearch(this, q)
}

func (this *Instagram) GetUpload() *Upload {
	if !this.IsLogin {
		return nil
	}
	return newUpload(this)
}

func (this *Instagram) GetAccount() *Account {
	if !this.IsLogin {
		return nil
	}
	return newAccount(this)
}

// SetProxy sets proxy for connection.
func (this *Instagram) SetProxy(_proxy *common.Proxy) {
	this.Proxy = _proxy
	this.c.Transport = _proxy.GetProxy()
	common.DebugHttpClient(this.c)
}

func (this *Instagram) ReadHeader(key string) string {
	return this.httpHeader[key]
}

func (this *Instagram) readMsisdnHeader() error {
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlMsisdnHeader,
			IsPost:  true,
			Query: map[string]interface{}{
				"device_id": this.uuid,
			},
		},
	)
	return err
}

//注册成功后触发
func (this *Instagram) contactPrefill() error {
	var query map[string]interface{}

	if this.IsLogin {
		query = map[string]interface{}{
			"_uid":      this.id,
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

func (this *Instagram) launcherSync() error {
	var query map[string]interface{}

	if this.IsLogin {
		query = map[string]interface{}{
			"id":                      this.id,
			"_uid":                    this.id,
			"_uuid":                   this.uuid,
			"server_config_retrieval": "1",
		}
	} else {
		query = map[string]interface{}{
			"id":                      this.uuid,
			"server_config_retrieval": "1",
		}
	}

	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlLauncherSync,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}

func (this *Instagram) zrToken() error {
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
			HeaderKey: []string{IGHeader_Authorization},
		},
	)
	return err
}

//早于注册登录?
func (this *Instagram) sendAdID() error {
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

func (this *Instagram) PrepareNewClient() {
	_ = this.readMsisdnHeader()
	_ = this.syncFeatures()
	_ = this.zrToken()
	_ = this.contactPrefill()
	_ = this.sendAdID()
	_ = this.launcherSync()
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
		"country_codes":       "[{\"country_code\":\"" + strings.ReplaceAll(this.registerPhoneArea, "+", "") + "\",\"source\":[\"default\"]}]",
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
		Login:   false,
		ApiPath: urlLogin,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	if err != nil && this.ReadHeader(IGHeader_Authorization) != "" {
		this.IsLogin = true
		this.id = strconv.FormatInt(resp.LoggedInUser.Pk, 10)
	}
	return err
}

// Logout closes current session
func (this *Instagram) Logout() error {
	_, err := this.sendSimpleRequest(urlLogout)
	this.c.Jar = nil
	this.c = nil
	return err
}

func (this *Instagram) syncFeatures() error {
	var params map[string]interface{}
	if this.IsLogin {
		params = map[string]interface{}{
			"id":          this.id,
			"_uid":        this.id,
			"_uuid":       this.uuid,
			"experiments": goInstaExperiments,
		}
	} else {
		params = map[string]interface{}{
			"id":          this.uuid,
			"experiments": goInstaExperiments,
		}
	}
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlQeSync,
			Query:   params,
			IsPost:  true,
			Login:   true,
			Signed:  true,
		},
	)
	return err
}

//func (this *Instagram) megaphoneLog() error {
//	_, err := this.HttpRequest(
//		&reqOptions{
//			ApiPath: urlMegaphoneLog,
//			Query: map[string]interface{}{
//				"id":        strconv.FormatInt(this.Account.ID, 10),
//				"type":      "feed_aysf",
//				"action":    "seen",
//				"reason":    "",
//				"device_id": this.androidID,
//				"uuid":      common.GenerateMD5Hash(string(time.Now().Unix())),
//			},
//			IsPost: true,
//			Login:  true,
//		},
//	)
//	return err
//}

//func (inst *Instagram) expose() error {
//	data, err := inst.prepareData(
//		map[string]interface{}{
//			"id":         inst.Account.ID,
//			"experiment": "ig_android_profile_contextual_feed",
//		},
//	)
//	if err != nil {
//		return err
//	}
//
//	_, err = inst.sendRequest(
//		&reqOptions{
//			ApiPath: urlExpose,
//			Query:    generateSignature(data),
//			IsPost:   true,
//		},
//	)
//
//	return err
//}

// GetMedia returns media specified by id.
//
// The argument can be int64 or string
//
// See example: examples/media/like.go
//func (inst *Instagram) GetMedia(o interface{}) (*FeedMedia, error) {
//	media := &FeedMedia{
//		inst:   inst,
//		NextID: o,
//	}
//	return media, media.Sync()
//}
