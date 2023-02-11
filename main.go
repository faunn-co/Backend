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

	//Affiliate
	e.POST("api/v1/affiliate/list", cmd.GetAffiliateList)
	//returns details of this affiliate
	e.GET("api/v1/affiliate/:id", cmd.GetAffiliateDetailsById)
	//get affiliate id from token
	e.POST("api/v1/affiliate/info", cmd.GetAffiliateInfo)
	//get affiliate id from token
	e.POST("api/v1/affiliate/stats", cmd.GetAffiliateStats)             //DONE
	e.POST("api/v1/affiliate/trend", cmd.GetAffiliateTrend)             //DONE
	e.GET("api/v1/affiliate/ranking/list", cmd.GetAffiliateRankingList) //DONE
	//Use by Affiliates to see their referrals
	e.GET("api/v1/affiliate/referral/list", cmd.GetAffiliateReferralsList)

	//Referral
	e.POST("api/v1/referral/list", cmd.GetReferralsList)
	e.POST("api/v1/referral/stats", cmd.GetReferralStats)
	e.POST("api/v1/referral/trend", cmd.GetReferralTrend)
	e.POST("api/v1/referral/recent/list", cmd.GetReferralRecentList)
	e.GET("api/v1/referral/:id", cmd.GetReferralById)

	//Booking
	e.GET("api/v1/booking/list", cmd.GetAvailableSlot)
	e.POST("api/v1/booking/stats", cmd.GetAvailableSlot)
	e.POST("api/v1/booking/trend", cmd.GetAvailableSlot)
	e.GET("api/v1/booking/recent/list", cmd.GetAvailableSlot)
	e.GET("api/v1/booking/:id", cmd.GetAvailableSlot)
	e.PUT("api/v1/booking/:id", cmd.GetAvailableSlot)
	e.DELETE("api/v1/booking/:id", cmd.GetAvailableSlot)

	//Landing Page
	e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot)
	e.POST("api/v1/booking/transaction/begin", cmd.GetAvailableSlot)
	e.POST("api/v1/booking/transaction/complete", cmd.GetAvailableSlot)

	//Registration
	e.POST("api/v1/platform/register", cmd.GetAvailableSlot)
	e.POST("api/v1/platform/login", cmd.GetAvailableSlot)

	e.Logger.Fatal(e.Start(":8888"))
}
