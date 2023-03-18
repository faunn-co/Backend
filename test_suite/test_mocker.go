package test_suite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/auth_middleware"
	"github.com/aaronangxz/AffiliateManager/cmd"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"net/http/httptest"
	"reflect"
)

const (
	GET                     = "GET"
	POST                    = "POST"
	DELETE                  = "DELETE"
	GetAffiliateList        = "/api/v1/affiliate/list"
	GetAffiliateStats       = "/api/v1/affiliate/stats"
	GetAffiliateTrend       = "/api/v1/affiliate/trend"
	GetAffiliateRankingList = "/api/v1/affiliate/ranking/list"
	GetAffiliateDetailsById = "/api/v1/affiliate/:id"
	GetReferralsList        = "/api/v1/referral/list"
	GetReferralStats        = "/api/v1/referral/stats"
	GetReferralTrend        = "/api/v1/referral/trend"
	GetReferralRecentList   = "/api/v1/referral/recent/list"
	GetReferralById         = "/api/v1/referral/:id"
	GetBookingList          = "/api/v1/booking/list"
	GetUserInfo             = "/api/v1/user/info"
	GetAvailableSlot        = "/api/v1/booking/slots/available"
	UserRegistration        = "/api/v1/platform/register"
	UserAuthentication      = "/api/v1/platform/login"
	UserDeAuthentication    = "/api/v1/platform/logout"
	TrackClick              = "/api/v1/tracking/click"
	TrackCheckout           = "/api/v1/tracking/checkout"
)

type Method struct {
	Endpoint   string
	HTTPMethod string
	Model      reflect.Type
}

var methodMap = map[string]Method{
	GetAffiliateStats: {
		Endpoint:   GetAffiliateStats,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetAffiliateStatsResponse{}),
	},
	GetAffiliateList: {
		Endpoint:   GetAffiliateList,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetAffiliateListResponse{}),
	},
	GetAffiliateTrend: {
		Endpoint:   GetAffiliateTrend,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetAffiliateTrendResponse{}),
	},
	GetAffiliateRankingList: {
		Endpoint:   GetAffiliateRankingList,
		HTTPMethod: GET,
		Model:      reflect.TypeOf(pb.GetAffiliateRankingListResponse{}),
	},
	GetAffiliateDetailsById: {
		Endpoint:   GetAffiliateDetailsById,
		HTTPMethod: GET,
		Model:      reflect.TypeOf(pb.GetAffiliateDetailsByIdResponse{}),
	},
	GetReferralsList: {
		Endpoint:   GetReferralsList,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetReferralListResponse{}),
	},
	GetReferralStats: {
		Endpoint:   GetReferralStats,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetReferralStatsResponse{}),
	},
	GetReferralTrend: {
		Endpoint:   GetReferralTrend,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetReferralTrendResponse{}),
	},
	GetReferralRecentList: {
		Endpoint:   GetReferralRecentList,
		HTTPMethod: GET,
		Model:      reflect.TypeOf(pb.GetReferralRecentListResponse{}),
	},
	GetReferralById: {
		Endpoint:   GetReferralById,
		HTTPMethod: GET,
		Model:      reflect.TypeOf(pb.GetReferralDetailsByReferralIdResponse{}),
	},
	GetBookingList: {
		Endpoint:   GetBookingList,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetBookingListResponse{}),
	},
	GetUserInfo: {
		Endpoint:   GetUserInfo,
		HTTPMethod: GET,
		Model:      reflect.TypeOf(pb.GetUserInfoResponse{}),
	},
	GetAvailableSlot: {
		Endpoint:   GetAvailableSlot,
		HTTPMethod: GET,
		Model:      reflect.TypeOf(pb.GetAvailableSlotResponse{}),
	},
	UserRegistration: {
		Endpoint:   UserRegistration,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.UserRegistrationResponse{}),
	},
	UserAuthentication: {
		Endpoint:   UserAuthentication,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.UserAuthenticationResponse{}),
	},
	UserDeAuthentication: {
		Endpoint:   UserDeAuthentication,
		HTTPMethod: DELETE,
		Model:      reflect.TypeOf(pb.UserDeAuthenticationResponse{}),
	},
	TrackClick: {
		Endpoint:   TrackClick,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.TrackClickResponse{}),
	},
	TrackCheckout: {
		Endpoint:   TrackCheckout,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.TrackCheckOutResponse{}),
	},
}

type MockTest struct {
	e        *echo.Echo
	param    map[string]string
	r        *httptest.ResponseRecorder
	meta     Method
	HttpErr  int
	RespBody interface{}
}

func NewMockTest(method string) *MockTest {
	m := new(MockTest)
	m.meta = methodMap[method]
	m.e = echo.New()
	//Affiliate
	//Allows admin / dev only
	a := m.e.Group("api/v1/affiliate")
	a.Use(auth_middleware.AdminAuthorization)
	a.POST("/list", cmd.GetAffiliateList)               //DONE
	a.GET("/:id", cmd.GetAffiliateDetailsById)          //DONE, not tested
	a.POST("/stats", cmd.GetAffiliateStats)             //DONE
	a.POST("/trend", cmd.GetAffiliateTrend)             //DONE
	a.GET("/ranking/list", cmd.GetAffiliateRankingList) //DONE

	//Referral
	//Allows admin / affiliate / dev
	r := m.e.Group("api/v1/referral")
	r.Use(auth_middleware.AffiliateAuthorization)
	r.POST("/list", cmd.GetReferralsList)            //DONE
	r.POST("/stats", cmd.GetReferralStats)           //DONE
	r.POST("/trend", cmd.GetReferralTrend)           //DONE
	r.GET("/recent/list", cmd.GetReferralRecentList) //DONE
	r.GET("/:id", cmd.GetReferralById)               //DONE, not tested

	//Booking
	//Allows admin / dev only
	b := m.e.Group("api/v1/booking")
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
	u := m.e.Group("api/v1/user")
	u.Use(auth_middleware.AffiliateAuthorization)
	u.GET("/info", cmd.GetUserInfo) //DONE

	//Landing Page
	m.e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot) //DONE

	//Registration
	m.e.POST("api/v1/platform/register", cmd.UserRegistration)     //DONE
	m.e.POST("api/v1/platform/login", cmd.UserAuthentication)      //DONE
	m.e.DELETE("api/v1/platform/logout", cmd.UserDeAuthentication) //DONE

	//Tracking
	m.e.POST("api/v1/welcome/click", cmd.TrackClick)       //DONE
	m.e.POST("api/v1/welcome/checkout", cmd.TrackCheckout) //DONE

	//Stripe
	m.e.POST("api/v1/payment/create-payment-intent", cmd.CreatePaymentIntent) //DONE
	orm.DIR = "../../../orm/queries/"
	orm.ENV = "TEST"
	return m
}

func (m *MockTest) QueryParam(key, value string) *MockTest {
	if m.param == nil {
		m.param = make(map[string]string)
	}
	m.param[key] = value
	return m
}

func (m *MockTest) Req(body interface{}, tokens *pb.Tokens) *MockTest {
	url := m.meta.Endpoint
	if m.param != nil {
		url += "?"
		for k, v := range m.param {
			url += fmt.Sprintf("%v=%v", k, v)
		}
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return m
	}
	val, _ := json.MarshalIndent(body, "", "    ")
	request, _ := http.NewRequest(m.meta.HTTPMethod, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokens.GetAccessToken()))
	fmt.Println("******************************************")
	fmt.Println("\t\t\t\tRequest")
	fmt.Println("******************************************")
	fmt.Println(m.meta.HTTPMethod, url)
	fmt.Println("Token:", request.Header.Get("Authorization"))
	fmt.Println(string(val))
	writer := httptest.NewRecorder()
	m.e.ServeHTTP(writer, request)
	m.r = writer
	m.decode()
	return m
}

func (m *MockTest) decode() *MockTest {
	m.HttpErr = m.r.Code
	switch m.meta.Endpoint {
	case GetAffiliateList:
		var dest *pb.GetAffiliateListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetAffiliateStats:
		var dest *pb.GetAffiliateStatsResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetAffiliateTrend:
		var dest *pb.GetAffiliateTrendResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetAffiliateRankingList:
		var dest *pb.GetAffiliateRankingListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetAffiliateDetailsById:
		var dest *pb.GetAffiliateDetailsByIdResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetReferralsList:
		var dest *pb.GetReferralListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetReferralStats:
		var dest *pb.GetReferralStatsResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetReferralTrend:
		var dest *pb.GetReferralTrendResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetReferralRecentList:
		var dest *pb.GetReferralRecentListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetReferralById:
		var dest *pb.GetReferralDetailsByReferralIdRequest
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetBookingList:
		var dest *pb.GetBookingListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetUserInfo:
		var dest *pb.GetUserInfoResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetAvailableSlot:
		var dest *pb.GetAvailableSlotResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case UserRegistration:
		var dest *pb.UserRegistrationResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case UserAuthentication:
		var dest *pb.UserAuthenticationResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case UserDeAuthentication:
		var dest *pb.UserDeAuthenticationResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case TrackClick:
		var dest *pb.TrackClickResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case TrackCheckout:
		var dest *pb.TrackCheckOutResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	}
	val, _ := json.MarshalIndent(m.RespBody, "", "    ")
	fmt.Println("******************************************")
	fmt.Println("\t\t\t\tResponse")
	fmt.Println("******************************************")
	fmt.Println("HTTP", m.HttpErr)
	fmt.Println(string(val))
	return m
}

func (m *MockTest) Response() (int, interface{}) {
	return m.HttpErr, m.RespBody
}
