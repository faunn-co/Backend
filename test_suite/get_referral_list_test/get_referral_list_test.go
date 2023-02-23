package get_referral_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"net/http"
	"testing"
	"time"
)

func TestGetReferralList_Affiliate_ValidAffiliateId(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		reqBody          *pb.GetReferralListRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	reqBody = &pb.GetReferralListRequest{
		TimeSelector: &pb.TimeSelector{
			BaseTs: proto.Int64(time.Now().Unix()),
			Period: proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
		AffiliateId: r.Affiliate.UserInfo.UserId,
	}

	err, resp := run(reqBody)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.Equal(t, r.Affiliate.UserInfo.GetUserName(), resp.GetReferralList()[0].GetAffiliateName())
		assert.Equal(t, r.ReferralDb.GetReferralId(), resp.GetReferralList()[0].GetReferralId())
		assert.Equal(t, r.ReferralDb.GetReferralClickTime(), resp.GetReferralList()[0].GetReferralClickTime())
		assert.Equal(t, r.ReferralDb.GetReferralStatus(), resp.GetReferralList()[0].GetReferralStatus())
		assert.Equal(t, r.ReferralDb.GetReferralCommission(), resp.GetReferralList()[0].GetReferralCommission())
		assert.Equal(t, r.ReferralDb.GetBookingId(), resp.GetReferralList()[0].GetBookingRefId())
		assert.Equal(t, r.Booking.BookingDetails.GetCitizenTicketCount()+r.Booking.BookingDetails.GetTouristTicketCount(), resp.GetReferralList()[0].GetTotalTicketCount())
		assert.Equal(t, r.Booking.BookingDetails.GetCitizenTicketTotal()+r.Booking.BookingDetails.GetTouristTicketTotal(), resp.GetReferralList()[0].GetTotalTicketAmount())
	}
}

func TestGetReferralList_Affiliate_ValidAffiliateName(t *testing.T) {
	r := referral.New().Build()
	defer r.TearDown()

	var (
		reqBody          *pb.GetReferralListRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
	)

	reqBody = &pb.GetReferralListRequest{
		TimeSelector: &pb.TimeSelector{
			BaseTs: proto.Int64(time.Now().Unix()),
			Period: proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
		AffiliateName: r.Affiliate.UserInfo.UserName,
	}

	err, resp := run(reqBody)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.Equal(t, r.Affiliate.UserInfo.GetUserName(), resp.GetReferralList()[0].GetAffiliateName())
		assert.Equal(t, r.ReferralDb.GetReferralId(), resp.GetReferralList()[0].GetReferralId())
		assert.Equal(t, r.ReferralDb.GetReferralClickTime(), resp.GetReferralList()[0].GetReferralClickTime())
		assert.Equal(t, r.ReferralDb.GetReferralStatus(), resp.GetReferralList()[0].GetReferralStatus())
		assert.Equal(t, r.ReferralDb.GetReferralCommission(), resp.GetReferralList()[0].GetReferralCommission())
		assert.Equal(t, r.ReferralDb.GetBookingId(), resp.GetReferralList()[0].GetBookingRefId())
		assert.Equal(t, r.Booking.BookingDetails.GetCitizenTicketCount()+r.Booking.BookingDetails.GetTouristTicketCount(), resp.GetReferralList()[0].GetTotalTicketCount())
		assert.Equal(t, r.Booking.BookingDetails.GetCitizenTicketTotal()+r.Booking.BookingDetails.GetTouristTicketTotal(), resp.GetReferralList()[0].GetTotalTicketAmount())
	}
}

func TestGetReferralList_Affiliate_InvalidAffiliateId(t *testing.T) {
	var (
		reqBody          *pb.GetReferralListRequest
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)

	reqBody = &pb.GetReferralListRequest{
		TimeSelector: &pb.TimeSelector{
			BaseTs: proto.Int64(time.Now().Unix()),
			Period: proto.Int64(int64(pb.TimeSelectorPeriod_PERIOD_MONTH)),
		},
		AffiliateId: proto.Int64(69420),
	}

	err, resp := run(reqBody)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetReferralList())
}
