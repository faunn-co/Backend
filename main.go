package main

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/cmd"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

func main() {
	logger.InitializeLogger()
	loadEnv()
	e := echo.New()
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
	e.GET("api/v1/affiliate/info", cmd.GetAffiliateInfo)                //DONE
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
	e.POST("api/v1/platform/register", cmd.UserRegistration) //DONE
	e.POST("api/v1/platform/login", cmd.UserAuthentication)  //DONE

	e.POST("api/v1/tracking/click", cmd.TrackClick)       //DONE
	e.POST("api/v1/tracking/checkout", cmd.TrackCheckout) //DONE

	e.POST("api/v1/payment/create-payment-intent", cmd.CreatePaymentIntent)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", getPort())))
}

func getPort() string {
	p := os.Getenv("PORT")
	if p == "" {
		return "8888"
	}
	return p
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Warn(context.Background(), "Error loading .env file")
	}

	if os.Getenv("ENV") == "LOCAL" {
		orm.ENV = "LOCAL"
		logger.Info(context.Background(), "Running on LOCAL")
	}
}
