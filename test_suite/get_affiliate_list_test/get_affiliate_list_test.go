package get_affiliate_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/admin"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"net/http"
	"testing"
	"time"
)

func TestGetAffiliateList_Happy(t *testing.T) {
	var (
		reqBody          *pb.GetAffiliateListRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
		refCount         = 2
	)
	a := admin.New().GenerateReferrals(int32(refCount)).Build()
	defer a.TearDown()

	reqBody = &pb.GetAffiliateListRequest{
		RequestMeta: nil,
		TimeSelector: &pb.TimeSelector{
			BaseTs: proto.Int64(time.Now().Unix()),
			Period: proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
	}

	err, resp := run(reqBody, a.User.Token)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.Equal(t, refCount, len(resp.GetAffiliateList()))
		for i, l := range resp.GetAffiliateList() {
			assert.Equal(t, a.Referrals[i].Affiliate.AffiliateInfo.GetUserId(), l.GetAffiliateId())
			assert.Equal(t, a.Referrals[i].Affiliate.UserInfo.GetUserName(), l.GetAffiliateName())
			assert.Equal(t, a.Referrals[i].Affiliate.AffiliateInfo.GetAffiliateType(), l.GetAffiliateType())
			assert.Equal(t, a.Referrals[i].Affiliate.AffiliateInfo.GetUniqueReferralCode(), l.GetUniqueReferralCode())
			assert.Equal(t, a.Referrals[i].ReferralDb.GetReferralCommission(), l.GetReferralCommission())
			assert.Equal(t, a.Referrals[i].Booking.BookingDetails.GetCitizenTicketTotal()+a.Referrals[i].Booking.BookingDetails.GetTouristTicketTotal(), l.GetTotalRevenue())
		}
	}
}

func TestGetAffiliateList_NoLogin(t *testing.T) {
	var (
		reqBody          *pb.GetAffiliateListRequest
		expectedHTTPCode = http.StatusUnauthorized
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NOT_AUTHORISED)
	)

	reqBody = &pb.GetAffiliateListRequest{
		RequestMeta: nil,
		TimeSelector: &pb.TimeSelector{
			BaseTs: proto.Int64(time.Now().Unix()),
			Period: proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
	}

	err, resp := run(reqBody, nil)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
}

func TestGetAffiliateList_NoPermission(t *testing.T) {
	var (
		reqBody          *pb.GetAffiliateListRequest
		expectedHTTPCode = http.StatusForbidden
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_NO_ACCESS)
	)
	r := referral.New().Build()
	defer r.TearDown()

	reqBody = &pb.GetAffiliateListRequest{
		RequestMeta: nil,
		TimeSelector: &pb.TimeSelector{
			BaseTs: proto.Int64(time.Now().Unix()),
			Period: proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
	}

	err, resp := run(reqBody, r.Affiliate.Token)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
}
