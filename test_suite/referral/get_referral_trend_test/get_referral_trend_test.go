package get_referral_trend_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/cmd/referral/get_referral_trend"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetTimesStats_Happy(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_NoLogin(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusUnauthorized
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NOT_AUTHORISED)
	)

	req := get_referral_trend.New().Request()
	err, resp := run(req, nil)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_None(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().SetPeriod(pb.TimeSelectorPeriod_PERIOD_NONE).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Day(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().SetPeriod(pb.TimeSelectorPeriod_PERIOD_DAY).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Week(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().SetPeriod(pb.TimeSelectorPeriod_PERIOD_WEEK).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Month(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().SetPeriod(pb.TimeSelectorPeriod_PERIOD_MONTH).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Last7Days(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().SetPeriod(pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Last28Days(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().SetPeriod(pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Range(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().
		SetPeriod(pb.TimeSelectorPeriod_PERIOD_RANGE).
		SetStartTimeStamp(time.Now().Unix() - 3*utils.DAY).
		SetEndTimeStamp(time.Now().Unix()).
		Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetTimesStats())
}

func TestGetTimesStats_Period_Range_StartBeforeEnd(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_trend.New().
		SetPeriod(pb.TimeSelectorPeriod_PERIOD_RANGE).
		SetStartTimeStamp(time.Now().Unix()).
		SetEndTimeStamp(time.Now().Unix() - 3*utils.DAY).
		Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetTimesStats())
}
