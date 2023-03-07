package get_available_slot_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/cmd/landing/get_available_slot"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetUserInfo_Happy(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_SUCCESS)
		date             = utils.ConvertTimeStampYearMonthDay(time.Now().Unix())
	)

	req := get_available_slot.New().Request()
	err, resp := run(req, nil, date)

	assert.Equal(t, expectedHTTPCode, err)
	if assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode()) {
		assert.Equal(t, date, resp.GetDate())
		assert.NotNil(t, resp.GetBookingSlots())
	}
}

func TestGetUserInfo_InvalidDate(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	)

	req := get_available_slot.New().Request()
	err, resp := run(req, nil, "123")

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetBookingSlots())
}

func TestGetUserInfo_PreviousDay(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
		date             = utils.ConvertTimeStampYearMonthDay(time.Now().Unix() - utils.DAY)
	)

	req := get_available_slot.New().Request()
	err, resp := run(req, nil, date)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetBookingSlots())
}

func TestGetUserInfo_SubsequentMonth(t *testing.T) {
	var (
		expectedHTTPCode = http.StatusOK
		expectedErrCode  = int64(pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
		date             = utils.ConvertTimeStampYearMonthDay(time.Now().Unix() + utils.MONTH)
	)

	req := get_available_slot.New().Request()
	err, resp := run(req, nil, date)

	assert.Equal(t, expectedHTTPCode, err)
	assert.Equal(t, expectedErrCode, resp.GetResponseMeta().GetErrorCode())
	assert.Nil(t, resp.GetBookingSlots())
}
