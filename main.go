package main

import (
	"github.com/aaronangxz/AffiliateManager/cmd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	e := echo.New()
	logger, _ := zap.NewProduction()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("ID", v.RequestID),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	//Affiliate
	e.POST("api/v1/affiliate/list", cmd.GetAffiliateList)               //DONE
	e.GET("api/v1/affiliate/:id", cmd.GetAffiliateDetailsById)          //DONE, not tested
	e.GET("api/v1/affiliate/info", cmd.GetAffiliateInfo)                //DONE, not tested
	e.POST("api/v1/affiliate/stats", cmd.GetAffiliateStats)             //DONE
	e.POST("api/v1/affiliate/trend", cmd.GetAffiliateTrend)             //DONE
	e.GET("api/v1/affiliate/ranking/list", cmd.GetAffiliateRankingList) //DONE

	//Referral
	//Use by Admin/Affiliates
	e.POST("api/v1/referral/list", cmd.GetReferralsList)             //DONE
	e.POST("api/v1/referral/stats", cmd.GetReferralStats)            //DONE
	e.POST("api/v1/referral/trend", cmd.GetReferralTrend)            //DONE
	e.POST("api/v1/referral/recent/list", cmd.GetReferralRecentList) //DONE
	e.GET("api/v1/referral/:id", cmd.GetReferralById)                //DONE, not tested

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

	e.POST("api/v1/tracking/click", cmd.TrackClick) //DONE, not tested
	e.POST("api/v1/tracking/checkout", cmd.TrackClick)

	e.Logger.Fatal(e.Start(":8888"))
}
