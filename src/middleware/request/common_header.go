package request

import (
	"context"
	"net/http"
	"strconv"
)

const (
	xHeaderKeyDeviceBrand = "X-Device-Brand"
	xHeaderKeyDeviceModel = "X-Device-Model"
	xHeaderKeyDeviceOS    = "X-Device-OS"
	xHeaderKeyPlatform    = "X-Platform"
	xHeaderKeyVersionName = "X-Version-Name"
	xHeaderKeyVersionCode = "X-Version-Code"
	xHeaderUserLocale     = "X-User-Locale"
	HeaderAcceptLanguage  = "Accept-Language"
	headerKeyLanguage     = "Lang"
)

type (
	ctxKeyCommonHeaders struct{}

	CommonHeaders struct {
		DeviceBrand string
		DeviceModel string
		DeviceOS    string
		Language    string
		Platform    string
		VersionName string
		VersionCode int64
	}
)

var CtxKeyCommonHeaders = ctxKeyCommonHeaders{}

func RequestAttributesContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versionCode, _ := strconv.ParseInt(r.Header.Get(xHeaderKeyVersionCode), 10, 64)
		commonHeader := CommonHeaders{
			DeviceBrand: r.Header.Get(xHeaderKeyDeviceBrand),
			DeviceModel: r.Header.Get(xHeaderKeyDeviceModel),
			DeviceOS:    r.Header.Get(xHeaderKeyDeviceOS),
			Platform:    r.Header.Get(xHeaderKeyPlatform),
			VersionName: r.Header.Get(xHeaderKeyVersionName),
			VersionCode: versionCode,
		}

		ctx := context.WithValue(r.Context(), CtxKeyCommonHeaders, commonHeader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCommonHeaders(ctx context.Context) CommonHeaders {
	if v, ok := ctx.Value(CtxKeyCommonHeaders).(CommonHeaders); ok {
		return v
	}
	return CommonHeaders{}
}

func GetLanguage(ctx context.Context) string {
	return GetCommonHeaders(ctx).Language
}

func GetPlatform(ctx context.Context) string {
	return GetCommonHeaders(ctx).Platform
}
