package test_suite

import (
	"bytes"
	"encoding/json"
	"fmt"
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
}

type MockTest struct {
	e        *echo.Echo
	param    map[string]string
	r        *httptest.ResponseRecorder
	meta     Method
	httpErr  int
	respBody interface{}
}

func NewMockTest(method string) *MockTest {
	m := new(MockTest)
	m.meta = methodMap[method]
	m.e = echo.New()
	//Affiliate
	m.e.POST("api/v1/affiliate/list", cmd.GetAffiliateList)
	//returns details of this affiliate
	m.e.GET("api/v1/affiliate/:id", cmd.GetAffiliateDetailsById)
	//get affiliate id from token
	m.e.POST("api/v1/affiliate/info", cmd.GetAffiliateInfo)
	//get affiliate id from token
	m.e.POST("api/v1/affiliate/stats", cmd.GetAffiliateStats)             //DONE
	m.e.POST("api/v1/affiliate/trend", cmd.GetAffiliateTrend)             //DONE
	m.e.GET("api/v1/affiliate/ranking/list", cmd.GetAffiliateRankingList) //DONE
	//Use by Affiliates to see their referrals
	m.e.GET("api/v1/affiliate/referral/list", cmd.GetAffiliateReferralsList)

	//Referral
	m.e.POST("api/v1/referral/list", cmd.GetReferralsList)
	m.e.POST("api/v1/referral/stats", cmd.GetReferralStats)
	m.e.POST("api/v1/referral/trend", cmd.GetReferralTrend)
	m.e.POST("api/v1/referral/recent/list", cmd.GetReferralRecentList)
	m.e.GET("api/v1/referral/:id", cmd.GetReferralById)

	//Booking
	m.e.POST("api/v1/booking/list", cmd.GetBookingList)
	m.e.POST("api/v1/booking/stats", cmd.GetAvailableSlot)
	m.e.POST("api/v1/booking/trend", cmd.GetAvailableSlot)
	m.e.GET("api/v1/booking/recent/list", cmd.GetAvailableSlot)
	m.e.GET("api/v1/booking/:id", cmd.GetAvailableSlot)
	m.e.PUT("api/v1/booking/:id", cmd.GetAvailableSlot)
	m.e.DELETE("api/v1/booking/:id", cmd.GetAvailableSlot)

	//Landing Page
	m.e.GET("api/v1/booking/slots/available", cmd.GetAvailableSlot)
	m.e.POST("api/v1/booking/transaction/begin", cmd.GetAvailableSlot)
	m.e.POST("api/v1/booking/transaction/complete", cmd.GetAvailableSlot)

	//Registration
	m.e.POST("api/v1/platform/register", cmd.GetAvailableSlot)
	m.e.POST("api/v1/platform/login", cmd.GetAvailableSlot)
	orm.DIR = "../orm/queries/"
	orm.ENV = "TEST"
	return m
}

func (m *MockTest) queryParam(key, value string) *MockTest {
	if m.param == nil {
		m.param = make(map[string]string)
	}
	m.param[key] = value
	return m
}

func (m *MockTest) req(body interface{}) *MockTest {
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
	fmt.Println(m.meta.HTTPMethod, url)
	fmt.Println(string(val))
	request, _ := http.NewRequest(m.meta.HTTPMethod, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()
	m.e.ServeHTTP(writer, request)
	m.r = writer
	return m
}

func (m *MockTest) decode() *MockTest {
	m.httpErr = m.r.Code
	switch m.meta.Endpoint {
	case GetAffiliateStats:
		var dest *pb.GetAffiliateStatsResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.respBody = dest
		break
	case GetAffiliateList:
		var dest *pb.GetAffiliateListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.respBody = dest
		break
	case GetAffiliateTrend:
		var dest *pb.GetAffiliateTrendResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.respBody = dest
		break
	case GetAffiliateRankingList:
		var dest *pb.GetAffiliateRankingListResponse
		err := json.Unmarshal(m.r.Body.Bytes(), &dest)
		if err != nil {
			log.Error(err)
			return m
		}
		m.respBody = dest
		break
	}
	val, _ := json.MarshalIndent(m.respBody, "", "    ")
	fmt.Println("************************************")
	fmt.Println(string(val))
	return m
}
