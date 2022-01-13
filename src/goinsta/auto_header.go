package goinsta

import (
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"strconv"
	"time"
)

func GetAutoHeaderFunc(header []string) []AutoSetHeaderFun {
	ret := make([]AutoSetHeaderFun, len(header))
	index := 0
	for _, item := range header {
		switch item {
		case "X-Entity-Length":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
			}
			index++
			break
		case "X-Idfa":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Idfa", inst.version.IDFA)
			}
			index++
			break
		case "X-Ig-App-Locale":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-App-Locale", inst.version.AppLocale)
			}
			index++
			break
		case "X-Mid":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Mid", inst.GetHeader("X-Mid"))
			}
			index++
			break
		case "X-Ig-Eu-Configure-Disabled":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Eu-Configure-Disabled", "true")
			}
			index++
			break
		case "X-Ig-Timezone-Offset":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Timezone-Offset", inst.version.TimezoneOffset)
			}
			index++
			break
		case "Ig-U-Ds-User-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				userID := inst.GetHeader("Ig-U-Ds-User-Id")
				if userID != "" {
					req.Header.Set("Ig-U-Ds-User-Id", userID)
				} else {
					log.Warn("user: %s ignore header Ig-U-Ds-User-Id", inst.User)
				}

			}
			index++
			break
		case "Media_hash":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("Media_hash")
			}
			index++
			break
		case "X-Pigeon-Session-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Pigeon-Session-Id", inst.sessionID)
			}
			index++
			break
		case "X-Fb-Http-Engine":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Fb-Http-Engine", "Liger")
			}
			index++
			break
		case "Ig-U-Rur":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				rur := inst.GetHeader("Ig-U-Rur")
				if rur != "" {
					req.Header.Set("Ig-U-Rur", rur)
				} else {
					log.Warn("user: %s ignore header Ig-U-Rur", inst.User)
				}
			}
			index++
			break
		case "X-Ig-Connection-Speed":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Connection-Speed", InstagramReqSpeed[common.GenNumber(0, len(InstagramReqSpeed))])
			}
			index++
			break
		case "X-Ig-Bandwidth-Speed-Kbps":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Bandwidth-Speed-Kbps", "0.000")
			}
			index++
			break
		case "Connection":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Connection", "close")
			}
			index++
			break
		case "X-Fb-Server-Cluster":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Fb-Server-Cluster", "True")
			}
			index++
			break
		case "X-Instagram-Rupload-Params":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X-Instagram-Rupload-Params")
			}
			index++
			break
		case "X-Ig-App-Startup-Country":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-App-Startup-Country", inst.version.StartupCountry)
			}
			index++
			break
		case "X-Pigeon-Rawclienttime":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Pigeon-Rawclienttime", strconv.FormatInt(time.Now().Unix(), 10)+".000000")
			}
			index++
			break
		case "Content-Type":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			}
			index++
			break
		case "X-Ig-Prefetch-Request":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X-Ig-Prefetch-Request")
			}
			index++
			break
		case "X-Bloks-Is-Panorama-Enabled":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Bloks-Is-Panorama-Enabled", "true")
			}
			index++
			break
		case "X-Fb-Client-Ip":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Fb-Client-Ip", "True")
			}
			index++
			break
		case "X-Ig-Device-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Device-Id", inst.DeviceID)
			}
			index++
			break
		case "X-Fb":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X-Fb")
			}
			index++
			break
		case "Ig-Intended-User-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Ig-Intended-User-Id", strconv.FormatInt(inst.ID, 10))
			}
			index++
			break
		case "X-Ig-Capabilities":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Capabilities", "36r/Fx8=")
			}
			index++
			break
		case "Accept-Language":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Accept-Language", inst.version.AcceptLanguage)
			}
			index++
			break
		case "X-Ig-Connection-Type":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Connection-Type", inst.version.NetWorkType)
			}
			index++
			break
		case "X-Bloks-Version-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Bloks-Version-Id", inst.version.BloksVersionID)
			}
			index++
			break
		case "Authorization-Others":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Authorization-Others", "")
			}
			index++
			break
		case "Authorization":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Authorization", inst.GetHeader("Authorization"))
			}
			index++
			break
		case "Content-Length":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("Content-Length")
			}
			index++
			break
		case "X-Entity-Type":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X-Entity-Type",)
			}
			index++
			break
		case "X-Ads-Opt-Out":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X-Ads-Opt-Out")
			}
			index++
			break
		case "X-Ig-Device-Locale":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Device-Locale", inst.version.AppLocale)
			}
			index++
			break
		case "X-Ig-Www-Claim":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				claim := inst.GetHeader("X-Ig-Www-Claim")
				if claim == "" {
					claim = "0"
				}
				req.Header.Set("X-Ig-Www-Claim", claim)
			}
			index++
			break
		case "X-Device-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Device-Id", inst.DeviceID)
			}
			index++
			break
		case "X-Ig-Family-Device-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Family-Device-Id", inst.familyID)
			}
			index++
			break
		case "X-Ig-App-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-App-Id", InstagramAppID)
			}
			index++
			break
		case "X-Ig-Mapped-Locale":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Mapped-Locale", inst.version.AppLocale)
			}
			index++
			break
		case "Accept-Encoding":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Accept-Encoding", "gzip, deflate")
			}
			index++
			break
		case "X-Entity-Name":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X-Entity-Name",)
			}
			index++
			break
		case "X-Ig-Abr-Connection-Speed-Kbps":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Abr-Connection-Speed-Kbps", "35")
			}
			index++
			break
		case "Offset":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set()
			}
			index++
			break
		case "Family_device_id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Family_device_id", inst.familyID)
			}
			index++
			break
		case "X_fb_photo_waterfall_id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				//req.Header.Set("X_fb_photo_waterfall_id",)
			}
			index++
			break
		case "User-Agent":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("User-Agent", inst.version.UserAgent)
			}
			index++
			break
		case "X-Tigon-Is-Retry":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Tigon-Is-Retry", "False")
			}
			index++
			break
		default:

		}
	}
	return ret
}