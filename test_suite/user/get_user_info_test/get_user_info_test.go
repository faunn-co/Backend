package get_user_info_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/admin"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/cmd/user/get_user_info"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetUserInfo_Happy(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_user_info.New().Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetUserInfo())
	assert.NotNil(t, resp.GetAffiliateMeta())
}

func TestGetUserInfo_NoLogin(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusUnauthorized
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NOT_AUTHORISED)
	)

	req := get_user_info.New().Request()
	err, resp := run(req, nil)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetUserInfo())
	assert.Nil(t, resp.GetAffiliateMeta())
}

func TestGetUserInfo_Role_Admin(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	a := admin.New().Build()
	defer a.TearDown()

	req := get_user_info.New().Request()
	err, resp := run(req, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetUserInfo())
	assert.NotNil(t, resp.GetAffiliateMeta())
}

func TestGetUserInfo_Role_Affiliate(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	req := get_user_info.New().Request()
	err, resp := run(req, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.NotNil(t, resp.GetUserInfo())
	assert.NotNil(t, resp.GetAffiliateMeta())
}
