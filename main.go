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
	"time"
)

func main() {
	logger.InitializeLogger()
	loadEnv()
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 20, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
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
	r.POST("/list", cmd.GetReferralsList)            //DONE
	r.POST("/stats", cmd.GetReferralStats)           //DONE
	r.POST("/trend", cmd.GetReferralTrend)           //DONE
	r.GET("/recent/list", cmd.GetReferralRecentList) //DONE
	r.GET("/:id", cmd.GetReferralById)               //DONE, not tested

	//Booking
	//Allows admin / dev only
	b := e.Group("api/v1/booking")
	b.Use(auth_middleware.AdminAuthorization)
	b.POST("/list", cmd.GetBookingList) //DONE
	//e.POST("api/v1/booking/list", cmd.GetBookingList) //DONE
	//e.POST("api/v1/booking/stats", cmd.GetAvailableSlot)
	//e.POST("api/v1/booking/trend", cmd.GetAvailableSlot)
	//e.GET("api/v1/booking/recent/list", cmd.GetAvailableSlot)
	//e.GET("api/v1/booking/:id", cmd.GetAvailableSlot)
	//e.PUT("api/v1/booking/:id", cmd.GetAvailableSlot)
	//e.DELETE("api/v1/booking/:id", cmd.GetAvailableSlot)

	//User
	u := e.Group("api/v1/user")
	u.Use(auth_middleware.AffiliateAuthorization)
	u.GET("/info", cmd.GetUserInfo) //DONE

	//Endpoints below require no Auth
	//Landing Page
	e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot) //DONE

	//Registration
	e.POST("api/v1/platform/register", cmd.UserRegistration)     //DONE
	e.POST("api/v1/platform/login", cmd.UserAuthentication)      //DONE
	e.DELETE("api/v1/platform/logout", cmd.UserDeAuthentication) //DONE

	//Tracking
	e.POST("api/v1/welcome/click", cmd.TrackClick)       //DONE
	e.POST("api/v1/welcome/payment", cmd.TrackClick)     //DONE
	e.POST("api/v1/welcome/checkout", cmd.TrackCheckout) //DONE

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
