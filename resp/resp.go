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

func GetAffiliateStatsResponseJSON(c echo.Context, curr *pb.AffiliateStats, prev *pb.AffiliateStats) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateStatsResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateStats:              curr,
		AffiliateStatsPreviousCycle: prev,
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

func GetReferralStatsResponseJSON(c echo.Context, curr *pb.ReferralStats, prev *pb.ReferralStats) error {
	return c.JSON(http.StatusOK, pb.GetReferralStatsResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		ReferralStats:              curr,
		ReferralStatsPreviousCycle: prev,
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
