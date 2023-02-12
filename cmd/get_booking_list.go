package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/booking/get_booking_list"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GetBookingList(c echo.Context) error {
	if l, s, e, err := get_booking_list.New(c).GetBookingListImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GetBookingListResponseJSON(c, l, s, e)
	}
}
