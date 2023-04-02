package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/tracking/track_payment"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func TrackPayment(c echo.Context) error {
	if id, err := track_payment.New(c).TrackPaymentImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.TrackPaymentResponseJSON(c, id)
	}
}
