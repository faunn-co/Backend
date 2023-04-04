package get_affiliate_ranking_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/admin"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/cmd/affiliate/get_affiliate_ranking_list"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAffiliateRankingList_Happy(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_WEEK
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
		refCount         = 2
	)
	a := admin.New().GenerateReferrals(int32(refCount)).Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_NoLogin(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusUnauthorized
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NOT_AUTHORISED)
	)

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, nil, nil)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_NoPermission(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusForbidden
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NO_ACCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, r.Affiliate.Token, nil)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_Period_None(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_NONE
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_Period_Day(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_DAY
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_Period_Week(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_WEEK
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_Period_Month(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_MONTH
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_Period_Last7Days(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetAffiliateRanking())
}

func TestGetAffiliateRankingList_Period_Last28Days(t *testing.T) {
	var (
		period           = pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_affiliate_ranking_list.New().Request()
	err, resp := run(req, a.User.Token, &period)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetAffiliateRanking())
}
