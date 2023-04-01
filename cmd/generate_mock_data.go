package cmd

import (
	"github.com/aaronangxz/AffiliateManager/impl/dev/generate_mock_data"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

func GenerateMockData(c echo.Context) error {
	if count, err := generate_mock_data.New(c).GenerateMockDataImpl(); err != nil {
		return resp.JSONResp(c, err)
	} else {
		return resp.GenerateMockDataResponseJSON(c, count)
	}
}
