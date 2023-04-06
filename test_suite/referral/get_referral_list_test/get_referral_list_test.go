package get_referral_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/admin"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/cmd/referral/get_referral_list"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetReferralList_Affiliate_Happy(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func TestGetReferralList_Affiliate_InputRandomAffiliateName(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetAffiliateName("A").Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func TestGetReferralList_Affiliate_Period_None(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_NONE
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetReferralList())
}

func TestGetReferralList_Affiliate_Period_Day(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_DAY
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func TestGetReferralList_Affiliate_Period_Week(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_WEEK
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func TestGetReferralList_Affiliate_Period_Month(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_MONTH
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func Skip_TestGetReferralList_Affiliate_Period_Last7Days(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	time.Sleep(1 * time.Second)
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func Skip_TestGetReferralList_Affiliate_Period_Last28Days(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	time.Sleep(1 * time.Second)
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func TestGetReferralList_Affiliate_Period_Range(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_RANGE
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
		for _, l := range resp.GetReferralList() {
			assert.Equal(t, r.Affiliate.AffiliateInfo.GetEntityName(), l.GetAffiliateName())
		}
	}
}

func TestGetReferralList_Affiliate_Period_Range_StartBeforeEnd(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_RANGE
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)

	req := get_referral_list.New().
		SetPeriod(period).
		SetStartTimeStamp(time.Now().Unix()).
		SetEndTimeStamp(time.Now().Unix() - 3*utils.DAY).
		Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetReferralList())
}

func TestGetReferralList_Admin_Happy(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func TestGetReferralList_Admin_InputAffiliateName(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetAffiliateName(a.Referrals[0].Affiliate.AffiliateInfo.GetEntityName()).Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func TestGetReferralList_Admin_Period_None(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_NONE
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetReferralList())
}

func TestGetReferralList_Admin_Period_Day(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_DAY
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func TestGetReferralList_Admin_Period_Week(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_WEEK
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func TestGetReferralList_Admin_Period_Month(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_MONTH
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

// TODO Empty response when run via GH actions
func Skip_TestGetReferralList_Admin_Period_Last7Days(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_LAST_7_DAYS
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	time.Sleep(1 * time.Second)
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func Skip_TestGetReferralList_Admin_Period_Last28Days(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_LAST_28_DAYS
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	time.Sleep(1 * time.Second)
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func TestGetReferralList_Admin_Period_Range(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_RANGE
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	req := get_referral_list.New().SetPeriod(period).Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.NotNil(t, resp.GetReferralList())
	}
}

func TestGetReferralList_Admin_Period_Range_StartBeforeEnd(t *testing.T) {
	a := admin.New().Build()
	defer a.TearDown()

	var (
		period           = pb.TimeSelectorPeriod_PERIOD_RANGE
		expectedHTTPCode = http.StatusBadRequest
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)

	req := get_referral_list.New().
		SetPeriod(period).
		SetStartTimeStamp(time.Now().Unix()).
		SetEndTimeStamp(time.Now().Unix() - 3*utils.DAY).
		Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetReferralList())
}
