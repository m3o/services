package handler

import (
	"context"

	proto "github.com/kytra-app/welcome-api/proto"
)

// RequestDownloadLink send the Google Play Store & Apple App Store links to the user
func (h Handler) RequestDownloadLink(ctx context.Context, req *proto.RequestDownload, rsp *proto.RequestDownload) error {
	number := req.GetDialCode() + req.GetPhoneNumber()
	return h.sms.Send(number, "Welcome to Kytra 👋\n\n Download for iOS: https://apps.apple.com/us/app/kytra/id1484048641\n\nDownload for Android: https://play.google.com/store/apps/details?id=com.kytra.app")
}
