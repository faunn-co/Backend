package main

import (
	"github.com/aaronangxz/AffiliateManager/cmd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	//sort and filter
	e.GET("api/v1/affiliate/list", cmd.GetAffiliateList)
	//returns details of this affiliate
	e.GET("api/v1/affiliate/:id", cmd.GetAffiliateById)
	//get affiliate id from token
	e.GET("api/v1/affiliate/details", cmd.GetAffiliateDetails)
	//get affiliate id from token
	//time range etc
	e.POST("api/v1/affiliate/stats", cmd.GetAffiliateStats)

	//sort and filter
	e.GET("api/v1/referral/list", cmd.GetReferralsList)
	//time range etc
	e.POST("api/v1/referral/stats", cmd.GetReferralStats)
	//returns details of this referral click
	e.GET("api/v1/referral/:id", cmd.GetReferralById)
	//sort and filter
	e.GET("api/v1/affiliate/referral/list", cmd.GetAffiliateReferralsList)

	e.Logger.Fatal(e.Start(":8888"))
}
