package test_suite

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/user"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGetAffiliateStats(t *testing.T) {
	u := user.New().SetUserRole(int64(pb.UserRole_ROLE_ADMIN)).Build()
	defer u.TearDown()

	var (
		reqBody          *pb.GetAffiliateStatsRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	reqBody = &pb.GetAffiliateStatsRequest{
		RequestMeta: &pb.RequestMeta{
			UserToken: nil,
		},
		TimeSelector: &pb.TimeSelector{
			BaseTs:  proto.Int64(time.Now().Unix()),
			StartTs: nil,
			EndTs:   nil,
			Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
	}

	err, response := test_suite.NewMockTest(test_suite.GetAffiliateStats).Req(reqBody, u.Token).Response()
	resp := response.(*pb.GetAffiliateStatsResponse)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
}

func TestGetAffiliateList(t *testing.T) {
	u := user.New().SetUserRole(int64(pb.UserRole_ROLE_ADMIN)).Build()
	defer u.TearDown()

	var (
		reqBody          *pb.GetAffiliateListRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	reqBody = &pb.GetAffiliateListRequest{
		RequestMeta: &pb.RequestMeta{
			UserToken: nil,
		},
		TimeSelector: &pb.TimeSelector{
			BaseTs:  proto.Int64(time.Now().Unix()),
			StartTs: nil,
			EndTs:   nil,
			Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
	}

	err, response := test_suite.NewMockTest(test_suite.GetAffiliateList).Req(reqBody, u.Token).Response()
	resp := response.(*pb.GetAffiliateListResponse)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.GreaterOrEqual(t, len(resp.GetAffiliateList()), 0)
}

func TestGetAffiliateTrend(t *testing.T) {
	u := user.New().SetUserRole(int64(pb.UserRole_ROLE_ADMIN)).Build()
	defer u.TearDown()

	var (
		reqBody          *pb.GetAffiliateTrendRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	reqBody = &pb.GetAffiliateTrendRequest{
		RequestMeta: &pb.RequestMeta{
			UserToken: nil,
		},
		TimeSelector: &pb.TimeSelector{
			BaseTs:  proto.Int64(time.Now().Unix()),
			StartTs: nil,
			EndTs:   nil,
			Period:  proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
	}

	err, response := test_suite.NewMockTest(test_suite.GetAffiliateTrend).Req(reqBody, u.Token).Response()
	resp := response.(*pb.GetAffiliateTrendResponse)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.GreaterOrEqual(t, len(resp.GetTimesStats()), 0)
}

func TestGetAffiliateRankingList(t *testing.T) {
	u := user.New().SetUserRole(int64(pb.UserRole_ROLE_ADMIN)).Build()
	defer u.TearDown()

	var (
		reqBody          *pb.GetAffiliateRankingListRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	reqBody = &pb.GetAffiliateRankingListRequest{
		RequestMeta: &pb.RequestMeta{
			UserToken: nil,
		},
	}

	err, response := test_suite.NewMockTest(test_suite.GetAffiliateRankingList).QueryParam("period", strconv.FormatInt(int64(pb.TimeSelectorPeriod_PERIOD_MONTH), 10)).Req(reqBody, u.Token).Response()
	resp := response.(*pb.GetAffiliateRankingListResponse)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
}
