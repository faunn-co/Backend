package resp

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func GetAvailableSlotResponseJSON(c echo.Context, date *string, slots []*pb.BookingSlot) error {
	return c.JSON(http.StatusOK, pb.GetAvailableSlotResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		Date:         date,
		BookingSlots: slots,
	})
}

func GetAffiliateListResponseJSON(c echo.Context, list []*pb.AffiliateMeta, start, end *int64) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateList: list,
		StartTime:     start,
		EndTime:       end,
	})
}

func GetAffiliateDetailsByIdResponseJSON(c echo.Context, meta *pb.AffiliateMeta, list []*pb.ReferralDetails) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateDetailsByIdResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateMeta: meta,
		ReferralList:  list,
	})
}

func GetAffiliateStatsResponseJSON(c echo.Context, r *pb.GetAffiliateStatsResponse) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateStatsResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateStats:              r.AffiliateStats,
		AffiliateStatsPreviousCycle: r.AffiliateStatsPreviousCycle,
	})
}

func GetAffiliateRankingListResponseJSON(c echo.Context, curr *pb.AffiliateRanking) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateRankingListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateRanking: curr,
	})
}

func GetAffiliateTrendResponseJSON(c echo.Context, trend []*pb.AffiliateCoreTimedStats) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateTrendResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		TimesStats: trend,
	})
}

func GetReferralStatsResponseJSON(c echo.Context, stats *pb.GetReferralStatsResponse) error {
	return c.JSON(http.StatusOK, pb.GetReferralStatsResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ReferralStats:              stats.ReferralStats,
		ReferralStatsPreviousCycle: stats.ReferralStatsPreviousCycle,
	})
}

func GetReferralTrendResponseJSON(c echo.Context, trend []*pb.ReferralCoreTimedStats) error {
	return c.JSON(http.StatusOK, pb.GetReferralTrendResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		TimesStats: trend,
	})
}

func GetReferralRecentListResponseJSON(c echo.Context, l *pb.ReferralRecent) error {
	return c.JSON(http.StatusOK, pb.GetReferralRecentListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ReferralRecent: l,
	})
}

func GetReferralListResponseJSON(c echo.Context, list []*pb.ReferralBasic, start, end *int64) error {
	return c.JSON(http.StatusOK, pb.GetReferralListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ReferralList: list,
		StartTime:    start,
		EndTime:      end,
	})
}

func GetBookingListResponseJSON(c echo.Context, list []*pb.BookingBasic, start, end *int64) error {
	return c.JSON(http.StatusOK, pb.GetBookingListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		Bookings:  list,
		StartTime: start,
		EndTime:   end,
	})
}

func GetReferralByIdResponseJSON(c echo.Context, details *pb.ReferralDetails) error {
	return c.JSON(http.StatusOK, pb.GetReferralDetailsByReferralIdResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ReferralDetails: details,
	})
}

func TrackClickResponseJSON(c echo.Context, id *int64) error {
	return c.JSON(http.StatusOK, pb.TrackClickResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ReferralId: id,
	})
}

func GetUserInfoResponseJSON(c echo.Context, meta *pb.AffiliateProfileMeta, user *pb.User) error {
	return c.JSON(http.StatusOK, pb.GetUserInfoResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateMeta: meta,
		UserInfo:      user,
	})
}

func UserRegistrationResponseJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, pb.UserRegistrationResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
	})
}

func UserAuthenticationResponseJSON(c echo.Context, a *pb.AuthCookie) error {
	return c.JSON(http.StatusOK, pb.UserAuthenticationResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AuthCookie: a,
	})
}

func CreatePaymentIntentResponseJSON(c echo.Context, secret *string) error {
	return c.JSON(http.StatusOK, pb.CreatePaymentIntentResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ClientSecret: secret,
	})
}

func TrackCheckoutResponseJSON(c echo.Context, details *pb.BookingDetails) error {
	return c.JSON(http.StatusOK, pb.TrackCheckOutResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		BookingDetails: details,
	})
}

func UserDeAuthenticationResponseJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, pb.UserDeAuthenticationResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
	})
}

func TrackPaymentResponseJSON(c echo.Context, id *int64) error {
	return c.JSON(http.StatusOK, pb.TrackPaymentResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		BookingId: id,
	})
}

func RollbackCheckoutResponseJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, pb.TrackCheckOutResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("successfully roll-backed transaction"),
		},
	})
}

func UpdateReferralByIdResponseJSON(c echo.Context, status *int64) error {
	return c.JSON(http.StatusOK, pb.UpdateReferralByIdResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		Status: status,
	})
}

func DeleteReferralByIdResponseJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, pb.DeleteReferralByIdResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
	})
}

func GenerateMockDataResponseJSON(c echo.Context, count *pb.MockDataCount) error {
	return c.JSON(http.StatusOK, pb.GenerateMockDataResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		MockDateCount: count,
	})
}
