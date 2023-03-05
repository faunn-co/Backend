package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/tracking/track_checkout"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func TrackCheckout(c echo.Context) error {
	if d, err := track_checkout.New(c).TrackCheckOutImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.TrackCheckoutResponseJSON(c, d)
	}
}
