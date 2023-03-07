package get_referral_recent_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/cmd/referral/get_referral_recent_list"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetReferralRecentList_Happy(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_referral_recent_list.New().Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetReferralRecent())
	assert.NotNil(t, resp.GetReferralRecent().GetRecentEarnings())
	assert.NotNil(t, resp.GetReferralRecent().GetRecentClicks())
}

func TestGetReferralRecentList_NoLogin(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusUnauthorized
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NOT_AUTHORISED)
	)

	req := get_referral_recent_list.New().Request()
	err, resp := run(req, nil)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetReferralRecent())
}
