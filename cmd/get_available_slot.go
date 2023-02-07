package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/booking/get_available_slot"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetAvailableSlot(c echo.Context) error {
	if d, s, err := get_available_slot.New(c).GetAvailableSlotImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetAvailableSlotResponseJSON(c, d, s)
	}
}
