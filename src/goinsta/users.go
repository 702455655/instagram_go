package goinsta

import (
	"fmt"
	"makemoney/common"
)

type UserSimp struct {
	Pk               int64  `json:"pk"`
	Username         string `json:"username"`
	FullName         string `json:"full_name"`
	IsPrivate        bool   `json:"is_private"`
	ProfilePicUrl    string `json:"profile_pic_url"`
	ProfilePicId     string `json:"profile_pic_id"`
	FriendshipStatus struct {
		Following                 bool `json:"following"`
		Blocking                  bool `json:"blocking"`
		IsPrivate                 bool `json:"is_private"`
		IncomingRequest           bool `json:"incoming_request"`
		OutgoingRequest           bool `json:"outgoing_request"`
		IsBestie                  bool `json:"is_bestie"`
		IsRestricted              bool `json:"is_restricted"`
		IsMessagingPseudoBlocking bool `json:"is_messaging_pseudo_blocking"`
		ReachabilityStatus        int  `json:"reachability_status"`
		IsFeedFavorite            bool `json:"is_feed_favorite"`
	} `json:"friendship_status"`
	IsVerified                   bool  `json:"is_verified"`
	FollowFrictionType           int   `json:"follow_friction_type"`
	HasAnonymousProfilePicture   bool  `json:"has_anonymous_profile_picture"`
	HasThreadsApp                bool  `json:"has_threads_app"`
	HasHighlightReels            bool  `json:"has_highlight_reels"`
	IsUsingUnifiedInboxForDirect bool  `json:"is_using_unified_inbox_for_direct"`
	BizUserInboxState            int   `json:"biz_user_inbox_state"`
	WaAddressable                bool  `json:"wa_addressable"`
	WaEligibility                int   `json:"wa_eligibility"`
	InteropMessagingUserFbid     int64 `json:"interop_messaging_user_fbid"`
}

type UserDetail struct {
	ID                         int64   `json:"pk"`
	Username                   string  `json:"username"`
	FullName                   string  `json:"full_name"`
	Biography                  string  `json:"biography"`
	ProfilePicURL              string  `json:"profile_pic_url"`
	Email                      string  `json:"email"`
	PhoneNumber                string  `json:"phone_number"`
	IsBusiness                 bool    `json:"is_business"`
	Gender                     int     `json:"gender"`
	ProfilePicID               string  `json:"profile_pic_id"`
	HasAnonymousProfilePicture bool    `json:"has_anonymous_profile_picture"`
	IsPrivate                  bool    `json:"is_private"`
	IsUnpublished              bool    `json:"is_unpublished"`
	AllowedCommenterType       string  `json:"allowed_commenter_type"`
	IsVerified                 bool    `json:"is_verified"`
	MediaCount                 int     `json:"media_count"`
	FollowerCount              int     `json:"follower_count"`
	FollowingCount             int     `json:"following_count"`
	FollowingTagCount          int     `json:"following_tag_count"`
	MutualFollowersID          []int64 `json:"profile_context_mutual_follow_ids"`
	ProfileContext             string  `json:"profile_context"`
	GeoMediaCount              int     `json:"geo_media_count"`
	ExternalURL                string  `json:"external_url"`
	HasBiographyTranslation    bool    `json:"has_biography_translation"`
	ExternalLynxURL            string  `json:"external_lynx_url"`
	BiographyWithEntities      struct {
		RawText  string        `json:"raw_text"`
		Entities []interface{} `json:"entities"`
	} `json:"biography_with_entities"`
	UsertagsCount                int           `json:"usertags_count"`
	HasChaining                  bool          `json:"has_chaining"`
	IsFavorite                   bool          `json:"is_favorite"`
	IsFavoriteForStories         bool          `json:"is_favorite_for_stories"`
	IsFavoriteForHighlights      bool          `json:"is_favorite_for_highlights"`
	CanBeReportedAsFraud         bool          `json:"can_be_reported_as_fraud"`
	ShowShoppableFeed            bool          `json:"show_shoppable_feed"`
	ShoppablePostsCount          int           `json:"shoppable_posts_count"`
	ReelAutoArchive              string        `json:"reel_auto_archive"`
	HasHighlightReels            bool          `json:"has_highlight_reels"`
	PublicEmail                  string        `json:"public_email"`
	PublicPhoneNumber            string        `json:"public_phone_number"`
	PublicPhoneCountryCode       string        `json:"public_phone_country_code"`
	ContactPhoneNumber           string        `json:"contact_phone_number"`
	CityID                       int64         `json:"city_id"`
	CityName                     string        `json:"city_name"`
	AddressStreet                string        `json:"address_street"`
	DirectMessaging              string        `json:"direct_messaging"`
	Latitude                     float64       `json:"latitude"`
	Longitude                    float64       `json:"longitude"`
	Category                     string        `json:"category"`
	BusinessContactMethod        string        `json:"business_contact_method"`
	IncludeDirectBlacklistStatus bool          `json:"include_direct_blacklist_status"`
	HdProfilePicURLInfo          PicURLInfo    `json:"hd_profile_pic_url_info"`
	HdProfilePicVersions         []PicURLInfo  `json:"hd_profile_pic_versions"`
	School                       School        `json:"school"`
	Byline                       string        `json:"byline"`
	SocialContext                string        `json:"social_context,omitempty"`
	SearchSocialContext          string        `json:"search_social_context,omitempty"`
	MutualFollowersCount         float64       `json:"mutual_followers_count"`
	LatestReelMedia              int64         `json:"latest_reel_media,omitempty"`
	IsCallToActionEnabled        bool          `json:"is_call_to_action_enabled"`
	FbPageCallToActionID         string        `json:"fb_page_call_to_action_id"`
	Zip                          string        `json:"zip"`
	Friendship                   Friendship    `json:"friendship_status"`
	AccountBadges                []interface{} `json:"account_badges"`

	//from self?
	BizUserInboxState            int     `json:"biz_user_inbox_state"`
	FollowFrictionType           int     `json:"follow_friction_type"`
	InteropMessagingUserFbid     int64   `json:"interop_messaging_user_fbid"`
	IsUsingUnifiedInboxForDirect bool    `json:"is_using_unified_inbox_for_direct"`
	NuxPrivateEnabled            bool    `json:"nux_private_enabled"`
	NuxPrivateFirstPage          bool    `json:"nux_private_first_page"`
	AccountType                  int     `json:"account_type"`
	CanSeeOrganicInsights        bool    `json:"can_see_organic_insights"`
	ShowInsightsTerms            bool    `json:"show_insights_terms"`
	Nametag                      Nametag `json:"nametag"`
	AllowContactsSync            bool    `json:"allow_contacts_sync"`
	CanBoostPost                 bool    `json:"can_boost_post"`
}

type User struct {
	inst    *Instagram `bson:"-"`
	HadSync bool       `bson:"-"`

	ID                         int64         `json:"pk" bson:"pk"`
	AccountBadges              []interface{} `json:"account_badges" bson:"account_badges"`
	FollowFrictionType         int           `json:"follow_friction_type" bson:"follow_friction_type"`
	FullName                   string        `json:"full_name" bson:"full_name"`
	HasAnonymousProfilePicture bool          `json:"has_anonymous_profile_picture" bson:"has_anonymous_profile_picture"`
	IsPrivate                  bool          `json:"is_private" bson:"is_private"`
	IsVerified                 bool          `json:"is_verified" bson:"is_verified"`
	LatestReelMedia            int           `json:"latest_reel_media" bson:"latest_reel_media"`
	ProfilePicId               string        `json:"profile_pic_id" bson:"profile_pic_id"`
	ProfilePicUrl              string        `json:"profile_pic_url" bson:"profile_pic_url"`
	Username                   string        `json:"username" bson:"username"`
	ShowPrivacyScreen          bool          `json:"show_privacy_screen" bson:"show_privacy_screen"`
	HasHighlightReels          bool          `json:"has_highlight_reels" bson:"has_highlight_reels"`
	Detail                     *UserDetail   `bson:"detail"`
}

type RespUserInfo struct {
	BaseApiResp
	User UserDetail `json:"user"`
}

func (this *User) Sync() error {
	var resp RespUserInfo
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: fmt.Sprintf(urlUserInfo, this.ID),
		Query: map[string]interface{}{
			"device_id": this.inst.AccountInfo.Device.DeviceID,
		}},
		resp)

	if err == nil {
		this.Detail = &resp.User
		this.HadSync = true
	}
	return err
}

type Followers struct {
	inst      *Instagram
	User      int64
	maxId     string
	rankToken string
	HasMore   bool
}

func (this *Followers) SetAccount(inst *Instagram) {
	this.inst = inst
}

type RespNexFollowers struct {
	BaseApiResp
	BigList   bool   `json:"big_list"`
	NextMaxId string `json:"next_max_id"`
	PageSize  int    `json:"page_size"`
	Users     []User `json:"users"`
}

func (this *Followers) Next() ([]User, error) {
	if !this.HasMore {
		return nil, &common.MakeMoneyError{
			ErrType: common.NoMoreError,
		}
	}

	var params map[string]interface{}

	if this.maxId != "" {
		params = map[string]interface{}{
			"max_id": this.maxId,
		}
	} else {
		params = map[string]interface{}{
			"rank_mutual": 0,
			"target_id":   this.User,
			"rank_token":  this.rankToken,
		}
	}

	resp := &RespNexFollowers{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: fmt.Sprintf(urlFriendFollowers, this.User),
		IsPost:  false,
		Query:   params,
	}, resp)
	err = resp.CheckError(err)

	if err == nil {
		if resp.NextMaxId != "" {
			this.maxId = resp.NextMaxId
		} else {
			this.HasMore = false
		}
	}

	return resp.Users[:], err
}
