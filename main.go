package main

import (
	"context"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
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
	//Allows admin / dev only
	a := e.Group("api/v1/affiliate")
	a.Use(auth_middleware.AdminAuthorization)
	a.POST("/list", cmd.GetAffiliateList)               //DONE
	a.GET("/:id", cmd.GetAffiliateDetailsById)          //DONE, not tested
	a.POST("/stats", cmd.GetAffiliateStats)             //DONE
	a.POST("/trend", cmd.GetAffiliateTrend)             //DONE
	a.GET("/ranking/list", cmd.GetAffiliateRankingList) //DONE

	//Referral
	//Allows admin / affiliate / dev
	r := e.Group("api/v1/referral")
	r.Use(auth_middleware.AffiliateAuthorization)
	r.POST("/list", cmd.GetReferralsList)             //DONE
	r.POST("/stats", cmd.GetReferralStats)            //DONE
	r.POST("/trend", cmd.GetReferralTrend)            //DONE
	r.POST("/recent/list", cmd.GetReferralRecentList) //DONE
	r.GET("/:id", cmd.GetReferralById)                //DONE, not tested

	//Booking
	//Allows admin / dev only
	b := e.Group("api/v1/booking")
	b.Use(auth_middleware.AdminAuthorization)
	b.POST("api/v1/booking/list", cmd.GetBookingList) //DONE
	//e.POST("api/v1/booking/list", cmd.GetBookingList) //DONE
	//e.POST("api/v1/booking/stats", cmd.GetAvailableSlot)
	//e.POST("api/v1/booking/trend", cmd.GetAvailableSlot)
	//e.GET("api/v1/booking/recent/list", cmd.GetAvailableSlot)
	//e.GET("api/v1/booking/:id", cmd.GetAvailableSlot)
	//e.PUT("api/v1/booking/:id", cmd.GetAvailableSlot)
	//e.DELETE("api/v1/booking/:id", cmd.GetAvailableSlot)

	//Endpoints below require no Auth
	e.GET("api/v1/user/info", cmd.GetUserInfo) //DONE

	//Landing Page
	e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot) //DONE

	//Registration
	e.POST("api/v1/platform/register", cmd.UserRegistration) //DONE
	e.POST("api/v1/platform/login", cmd.UserAuthentication)  //DONE

	//Tracking
	e.POST("api/v1/tracking/click", cmd.TrackClick)       //DONE
	e.POST("api/v1/tracking/checkout", cmd.TrackCheckout) //DONE

	//Stripe
	e.POST("api/v1/payment/create-payment-intent", cmd.CreatePaymentIntent) //DONE

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
