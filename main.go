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
	e.POST("api/v1/affiliate/list", cmd.GetAffiliateList) //DONE
	//returns details of this affiliate
	e.GET("api/v1/affiliate/:id", cmd.GetAffiliateDetailsById)
	//get affiliate id from token
	e.POST("api/v1/affiliate/info", cmd.GetAffiliateInfo)
	//get affiliate id from token
	e.POST("api/v1/affiliate/stats", cmd.GetAffiliateStats)             //DONE
	e.POST("api/v1/affiliate/trend", cmd.GetAffiliateTrend)             //DONE
	e.GET("api/v1/affiliate/ranking/list", cmd.GetAffiliateRankingList) //DONE

	//Referral
	//Use by Admin/Affiliates
	e.POST("api/v1/referral/list", cmd.GetReferralsList)             //DONE
	e.POST("api/v1/referral/stats", cmd.GetReferralStats)            //DONE
	e.POST("api/v1/referral/trend", cmd.GetReferralTrend)            //DONE
	e.POST("api/v1/referral/recent/list", cmd.GetReferralRecentList) //DONE
	e.GET("api/v1/referral/:id", cmd.GetReferralById)

	//Booking
	e.POST("api/v1/booking/list", cmd.GetBookingList) //DONE
	e.POST("api/v1/booking/stats", cmd.GetAvailableSlot)
	e.POST("api/v1/booking/trend", cmd.GetAvailableSlot)
	e.GET("api/v1/booking/recent/list", cmd.GetAvailableSlot)
	e.GET("api/v1/booking/:id", cmd.GetAvailableSlot)
	e.PUT("api/v1/booking/:id", cmd.GetAvailableSlot)
	e.DELETE("api/v1/booking/:id", cmd.GetAvailableSlot)

	//Landing Page
	e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot) //DONE
	e.POST("api/v1/booking/transaction/begin", cmd.GetAvailableSlot)
	e.POST("api/v1/booking/transaction/complete", cmd.GetAvailableSlot)

	//Registration
	e.POST("api/v1/platform/register", cmd.GetAvailableSlot)
	e.POST("api/v1/platform/login", cmd.GetAvailableSlot)

	e.Logger.Fatal(e.Start(":8888"))
}
