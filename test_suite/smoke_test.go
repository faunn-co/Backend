package test_suite

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGetAffiliateStats(t *testing.T) {
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

	m := NewMockTest(GetAffiliateStats)
	resp := m.req(reqBody).decode().respBody.(*pb.GetAffiliateStatsResponse)

	assert.Equal(t, expectedHTTPCode, m.httpErr)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
}

func TestGetAffiliateList(t *testing.T) {
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

	m := NewMockTest(GetAffiliateList)
	resp := m.req(reqBody).decode().respBody.(*pb.GetAffiliateListResponse)

	assert.Equal(t, expectedHTTPCode, m.httpErr)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.GreaterOrEqual(t, len(resp.GetAffiliateList()), 0)
}

func TestGetAffiliateTrend(t *testing.T) {
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

	m := NewMockTest(GetAffiliateTrend)
	resp := m.req(reqBody).decode().respBody.(*pb.GetAffiliateTrendResponse)

	assert.Equal(t, expectedHTTPCode, m.httpErr)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.GreaterOrEqual(t, len(resp.GetTimesStats()), 0)
}

func TestGetAffiliateRankingList(t *testing.T) {
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

	m := NewMockTest(GetAffiliateRankingList)
	resp := m.queryParam("period", strconv.FormatInt(int64(pb.TimeSelectorPeriod_PERIOD_MONTH), 10)).req(reqBody).decode().respBody.(*pb.GetAffiliateRankingListResponse)

	assert.Equal(t, expectedHTTPCode, m.httpErr)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
}
