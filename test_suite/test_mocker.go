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
	GetAffiliateStats       = "/api/v1/affiliate/stats"
	GetAffiliateList        = "/api/v1/affiliate/list"
	GetAffiliateTrend       = "/api/v1/affiliate/trend"
	GetAffiliateRankingList = "/api/v1/affiliate/ranking/list"
	GetReferralsList        = "/api/v1/referral/list"
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
	GetReferralsList: {
		Endpoint:   GetReferralsList,
		HTTPMethod: POST,
		Model:      reflect.TypeOf(pb.GetReferralListResponse{}),
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
	r.POST("/list", cmd.GetReferralsList)             //DONE
	r.POST("/stats", cmd.GetReferralStats)            //DONE
	r.POST("/trend", cmd.GetReferralTrend)            //DONE
	r.POST("/recent/list", cmd.GetReferralRecentList) //DONE
	r.GET("/:id", cmd.GetReferralById)                //DONE, not tested

	//Booking
	//Allows admin / dev only
	b := m.e.Group("api/v1/booking")
	b.Use(auth_middleware.AdminAuthorization)
	b.POST("api/v1/booking/list", cmd.GetBookingList) //DONE

	//Endpoints below require no Auth
	m.e.GET("api/v1/user/info", cmd.GetUserInfo) //DONE

	//Landing Page
	m.e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot) //DONE

	//Registration
	m.e.POST("api/v1/platform/register", cmd.UserRegistration) //DONE
	m.e.POST("api/v1/platform/login", cmd.UserAuthentication)  //DONE

	//Tracking
	m.e.POST("api/v1/tracking/click", cmd.TrackClick)       //DONE
	m.e.POST("api/v1/tracking/checkout", cmd.TrackCheckout) //DONE

	//Stripe
	m.e.POST("api/v1/payment/create-payment-intent", cmd.CreatePaymentIntent) //DONE
	orm.DIR = "../../orm/queries/"
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
	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return m
	}
	val, _ := json.MarshalIndent(body, "", "    ")

	url := m.meta.Endpoint
	if m.param != nil {
		url += "?"
		for k, v := range m.param {
			url += fmt.Sprintf("%v=%v", k, v)
		}
	}
	request, _ := http.NewRequest(m.meta.HTTPMethod, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokens.GetAccessToken()))
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
	case GetAffiliateStats:
		var dest *pb.GetAffiliateStatsResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	case GetAffiliateList:
		var dest *pb.GetAffiliateListResponse
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
	case GetReferralsList:
		var dest *pb.GetReferralListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.RespBody = dest
		break
	}
	val, _ := json.MarshalIndent(m.RespBody, "", "    ")
	fmt.Println("************************************")
	fmt.Println(string(val))
	return m
}

func (m *MockTest) Response() (int, interface{}) {
	return m.HttpErr, m.RespBody
}
