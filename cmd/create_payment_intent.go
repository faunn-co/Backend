package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/payment/create_payment_intent"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func CreatePaymentIntent(c echo.Context) error {
	if s, err := create_payment_intent.New(c).CreatePaymentIntentImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.CreatePaymentIntentResponseJSON(c, s)
	}
}
