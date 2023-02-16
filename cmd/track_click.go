package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/tracking/track_click"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func TrackClick(c echo.Context) error {
	if id, err := track_click.New(c).TrackClickImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.TrackClickResponseJSON(c, id)
	}
}
