package goinsta

import (
	"fmt"
	"makemoney/common"
	"time"
)

type Register struct {
	Inst         *Instagram
	RegisterType string
	Account      string
	AreaCode     string
	Username     string
	RealUsername string
	Password     string
	Year         string
	Month        string
	Day          string
	signUpCode   string
	tosVersion   string
}

func (this *Register) GetSignupConfig() error {
	params := map[string]interface{}{
		"device_id":             this.Inst.Device.DeviceID,
		"main_account_selected": "0",
	}

	_, err := this.Inst.HttpRequest(&reqOptions{
		ApiPath: urlGetSignupConfig,
		IsPost:  false,
		Query:   params,
	})

	return err
}

func (this *Register) GetCommonEmailDomains() error {
	params := map[string]interface{}{}

	_, err := this.Inst.HttpRequest(&reqOptions{
		ApiPath: urlGetCommonEmailDomains,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	})

	return err
}

func (this *Register) PrecheckCloudId() error {
	params := map[string]interface{}{
		"ck_container":   "iCloud.com.burbn.instagram",
		"ck_error":       "CKErrorDomain: 9",
		"ck_environment": "production",
	}

	_, err := this.Inst.HttpRequest(&reqOptions{
		ApiPath: urlPrecheckCloudId,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	})

	return err
}

func (this *Register) IgUser() error {
	params := map[string]interface{}{}

	_, err := this.Inst.HttpRequest(&reqOptions{
		ApiPath: urlIgUser,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	})

	return err
}

type RespCheckEmail struct {
	BaseApiResp
	Valid                        bool     `json:"valid"`
	Available                    bool     `json:"available"`
	AllowSharedEmailRegistration bool     `json:"allow_shared_email_registration"`
	UsernameSuggestions          []string `json:"username_suggestions"`
	TosVersion                   string   `json:"tos_version"`
	AgeRequired                  bool     `json:"age_required"`
}

func (this *Register) CheckEmail() (*RespCheckEmail, error) {
	params := map[string]interface{}{
		"email": this.Account,
		"qe_id": this.Inst.Device.DeviceID,
	}
	resp := &RespCheckEmail{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCheckEmail,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	if err == nil {
		this.tosVersion = resp.TosVersion
	}
	return resp, err
}

type RespSendVerifyEmail struct {
	BaseApiResp
	EmailSent bool   `json:"email_sent"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

func (this *Register) SendVerifyEmail() (*RespSendVerifyEmail, error) {
	params := map[string]interface{}{
		"email":        this.Account,
		"device_id":    this.Inst.Device.DeviceID,
		"phone_id":     this.Inst.Device.DeviceID,
		"waterfall_id": this.Inst.Device.WaterID,
	}

	resp := &RespSendVerifyEmail{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		ApiPath: urlSendVerifyEmail,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespCheckConfirmationCode struct {
	BaseApiResp
	SignupCode string `json:"signup_code"`
}

func (this *Register) CheckConfirmationCode(code string) (*RespCheckConfirmationCode, error) {
	params := map[string]interface{}{
		"email":            this.Account,
		"code":             code,
		"confirm_via_link": "0",
		"device_id":        this.Inst.Device.DeviceID,
		"phone_id":         this.Inst.Device.DeviceID,
		"waterfall_id":     this.Inst.Device.WaterID,
	}

	resp := &RespCheckConfirmationCode{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCheckConfirmationCode,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	if err == nil {
		this.signUpCode = resp.SignupCode
	}
	return resp, err
}

type RespCreatUser struct {
	BaseApiResp
	AccountCreated bool `json:"account_created"`
	CreatedUser    struct {
		Pk                         int64  `json:"pk"`
		Username                   string `json:"username"`
		FullName                   string `json:"full_name"`
		IsPrivate                  bool   `json:"is_private"`
		ProfilePicUrl              string `json:"profile_pic_url"`
		IsVerified                 bool   `json:"is_verified"`
		FollowFrictionType         int    `json:"follow_friction_type"`
		HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		ReelAutoArchive            string `json:"reel_auto_archive"`
		HdProfilePicVersions       []struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Url    string `json:"url"`
		} `json:"hd_profile_pic_versions"`
		HdProfilePicUrlInfo struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"hd_profile_pic_url_info"`
		NuxPrivateEnabled            bool          `json:"nux_private_enabled"`
		NuxPrivateFirstPage          bool          `json:"nux_private_first_page"`
		HasHighlightReels            bool          `json:"has_highlight_reels"`
		IsUsingUnifiedInboxForDirect bool          `json:"is_using_unified_inbox_for_direct"`
		BizUserInboxState            int           `json:"biz_user_inbox_state"`
		InteropMessagingUserFbid     int64         `json:"interop_messaging_user_fbid"`
		AccountBadges                []interface{} `json:"account_badges"`
		AllowContactsSync            bool          `json:"allow_contacts_sync"`
	} `json:"created_user"`
	MultipleUsersOnDevice bool        `json:"multiple_users_on_device"`
	SessionFlushNonce     interface{} `json:"session_flush_nonce"`
}

func (this *Register) setInstRegisterInfo(pk int64) {
	this.Inst.User = this.RealUsername
	this.Inst.Pass = this.Password
	this.Inst.RegisterEmail = this.Account
	this.Inst.RegisterTime = time.Now().Unix()
	this.Inst.RegisterIpCountry = this.Inst.Proxy.Country
	this.Inst.RegisterPhoneArea = this.AreaCode
	this.Inst.IsLogin = true
	this.Inst.ID = pk
}

func (this *Register) CreateEmail() (*RespCreatUser, error) {
	encodePasswd, err := encryptPassword(this.Password, this.Inst.GetHeader(IGHeader_EncryptionId), this.Inst.GetHeader(IGHeader_EncryptionKey))
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"do_not_auto_login_if_credentials_match": "0",
		"tos_version":                            this.tosVersion,
		"month":                                  this.Month,
		"device_id":                              this.Inst.Device.DeviceID,
		"ck_container":                           "iCloud.com.burbn.instagram",
		"has_seen_aart_on":                       "0",
		"ck_error":                               "CKErrorDomain: 9",
		"day":                                    this.Day,
		"waterfall_id":                           this.Inst.Device.WaterID,
		"year":                                   this.Year,
		"email":                                  this.Account,
		"enc_password":                           encodePasswd,
		"force_create_account":                   "0",
		"attribution_details":                    "{\n  \"Version3.1\" : {\n    \"iad-attribution\" : \"false\"\n  }\n}",
		"ck_environment":                         "production",
		"force_sign_up_code":                     this.signUpCode,
		"adid":                                   this.Inst.Device.IDFA,
		"phone_id":                               this.Inst.Device.DeviceID,
		"first_name":                             this.Username,
		"username":                               this.RealUsername,
	}

	resp := &RespCreatUser{}
	err = this.Inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCreate,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	if err == nil {
		this.setInstRegisterInfo(resp.CreatedUser.Pk)
	}
	return resp, err
}

//phone
func (this *Register) CheckPhoneNumber() error {
	params := map[string]interface{}{
		"phone_id":        this.Inst.Device.FamilyID,
		"login_nonce_map": "{}",
		"phone_number":    this.Account,
		"guid":            this.Inst.Device.DeviceID,
		"prefill_shown":   "False",
	}

	_, err := this.Inst.HttpRequest(&reqOptions{
		ApiPath: urlCheckPhoneNumber,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	})

	return err
}

type RespSendSignupSmsCode struct {
	BaseApiResp
	TosVersion  string `json:"tos_version"`
	AgeRequired bool   `json:"age_required"`
}

func (this *Register) SendSignupSmsCode() (*RespSendSignupSmsCode, error) {
	params := map[string]interface{}{
		"device_id":    this.Inst.Device.DeviceID,
		"phone_number": this.AreaCode + this.Account,
		"phone_id":     this.Inst.Device.DeviceID,
		"source":       "regular",
	}
	resp := &RespSendSignupSmsCode{}
	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSendSignupSmsCode,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespValidateSignupSmsCode struct {
	BaseApiResp
	PnTaken  bool `json:"pn_taken"`
	Verified bool `json:"verified"`
}

func (this *Register) ValidateSignupSmsCode(code string) (*RespValidateSignupSmsCode, error) {
	this.signUpCode = code
	params := map[string]interface{}{
		"device_id":         this.Inst.Device.DeviceID,
		"phone_number":      this.AreaCode + this.Account,
		"waterfall_id":      this.Inst.Device.WaterID,
		"verification_code": code,
	}
	resp := &RespValidateSignupSmsCode{}

	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlValidateSignupSmsCode,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespUsernameSuggestions struct {
	BaseApiResp
	Suggestions []string `json:"suggestions"`
}

func (this *Register) usernameSuggestions() (*RespUsernameSuggestions, error) {
	params := map[string]interface{}{}
	if this.RegisterType == "email" {
		params = map[string]interface{}{
			"email":        this.Account,
			"device_id":    this.Inst.Device.DeviceID,
			"name":         this.Username,
			"waterfall_id": this.Inst.Device.WaterID,
		}
	} else {
		params = map[string]interface{}{
			"name":         this.Username,
			"device_id":    this.Inst.Device.DeviceID,
			"waterfall_id": this.Inst.Device.WaterID,
		}
	}

	resp := &RespUsernameSuggestions{}

	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlUsernameSuggestions,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *Register) GenUsername() string {
	usernameSuggestions, err := this.usernameSuggestions()
	if err != nil || len(usernameSuggestions.Suggestions) == 0 {
		this.RealUsername = this.Username + fmt.Sprintf("%d%d%d",
			common.GenNumber(1990, 2020),
			common.GenNumber(1, 12),
			common.GenNumber(1, 27))
	} else {
		this.RealUsername = usernameSuggestions.Suggestions[0]
	}
	return this.RealUsername
}

type RespCheckAge struct {
	BaseApiResp
	EligibleToRegister      bool `json:"eligible_to_register"`
	ParentalConsentRequired bool `json:"parental_consent_required"`
	IsSupervisedUser        bool `json:"is_supervised_user"`
}

func (this *Register) CheckAgeEligibility() (*RespCheckAge, error) {
	params := map[string]interface{}{
		"year":  this.Year,
		"month": this.Month,
		"day":   this.Day,
	}

	resp := &RespCheckAge{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCheckAgeEligibility,
		IsPost:  true,
		Signed:  false,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespCheckUsernameError struct {
	Error string `json:"error"`
}

type RespCheckUsername struct {
	BaseApiResp
	RespCheckUsernameError
	Username             string `json:"username"`
	Available            bool   `json:"available"`
	ExistingUserPassword bool   `json:"existing_user_password"`
}

func (this *Register) CheckUsername() (*RespCheckUsername, error) {
	encodePasswd, err := encryptPassword(this.Password, this.Inst.GetHeader(IGHeader_EncryptionId), this.Inst.GetHeader(IGHeader_EncryptionKey))
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"enc_password": encodePasswd,
		"username":     this.Username,
		"device_id":    this.Inst.Device.DeviceID,
	}
	resp := &RespCheckUsername{}

	err = this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlCheckUsername,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *Register) NewUserFlowBegins() (*BaseApiResp, error) {
	params := map[string]interface{}{
		"device_id": this.Inst.Device.DeviceID,
	}

	resp := &BaseApiResp{}
	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlNewUserFlowBegins,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespCreateValidated struct {
	BaseApiResp
	AccountCreated        bool        `json:"account_created"`
	CreatedUser           UserDetail  `json:"created_user"`
	ExistingUser          bool        `json:"existing_user"`
	MultipleUsersOnDevice bool        `json:"multiple_users_on_device"`
	SessionFlushNonce     interface{} `json:"session_flush_nonce"`
}

func (this *Register) CreatePhone() (*RespCreateValidated, error) {
	encodePasswd, err := encryptPassword(this.Password, this.Inst.GetHeader(IGHeader_EncryptionId), this.Inst.GetHeader(IGHeader_EncryptionKey))
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"tos_version":                            this.tosVersion,
		"verification_code":                      this.signUpCode,
		"do_not_auto_login_if_credentials_match": "0",
		"phone_id":                               this.Inst.Device.DeviceID,
		"enc_password":                           encodePasswd,
		"phone_number":                           this.AreaCode + this.Account,
		"username":                               this.RealUsername,
		"first_name":                             this.Username,
		"day":                                    this.Day,
		"year":                                   this.Year,
		"device_id":                              this.Inst.Device.DeviceID,
		"month":                                  this.Month,
		"has_seen_aart_on":                       "0",
		"force_create_account":                   "0",
		"waterfall_id":                           this.Inst.Device.WaterID,
		"ck_error":                               "CKErrorDomain: 9",
		"has_sms_consent":                        "true",
		"ck_environment":                         "production",
		"ck_container":                           "iCloud.com.burbn.instagram",
	}
	resp := &RespCreateValidated{}

	err = this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlCreateValidated,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	if err == nil {
		this.setInstRegisterInfo(resp.CreatedUser.ID)
	}
	return resp, err
}

func (this *Register) NewAccountNuxSeen() (*BaseApiResp, error) {
	params := map[string]interface{}{}

	if this.RegisterType == "email" {
		params = map[string]interface{}{
			"_uuid":           this.Inst.Device.DeviceID,
			"_uid":            this.Inst.ID,
			"is_fb_installed": "false",
		}
	} else {
		params = map[string]interface{}{
			"is_fb_installed": false,
		}
	}

	resp := &BaseApiResp{}

	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlNewAccountNuxSeen,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *Register) GetSteps() (*BaseApiResp, error) {
	params := map[string]interface{}{}
	if this.RegisterType == "email" {
		params = map[string]interface{}{
			"_uuid":                         this.Inst.Device.DeviceID,
			"_uid":                          this.Inst.ID,
			"device_id":                     this.Inst.Device.DeviceID,
			"is_secondary_account_creation": "0",
			"push_permission_requested":     "0",
			"network_type":                  this.Inst.Device.NetWorkType + "-none",
			"is_account_linking_flow":       "0",
		}
	} else {
		params = map[string]interface{}{
			"device_id":                     this.Inst.Device.DeviceID,
			"is_secondary_account_creation": "0",
			"push_permission_requested":     "0",
			"network_type":                  this.Inst.Device.NetWorkType + "-none",
			"is_account_linking_flow":       "0",
		}
	}

	resp := &BaseApiResp{}

	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlGetSteps,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}
