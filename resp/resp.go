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

func GetAffiliateListResponseJSON(c echo.Context, list []*pb.AffiliateMeta) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateList: list,
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

func GetAffiliateRankingListResponseJSON(c echo.Context, curr *pb.AffiliateRanking, prev *pb.AffiliateRanking) error {
	return c.JSON(http.StatusOK, pb.GetAffiliateRankingListResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode: proto.Int64(int64(pb.GlobalErrorCode_SUCCESS)),
			ErrorMsg:  proto.String("success"),
		},
		AffiliateRanking:              curr,
		AffiliateRankingPreviousCycle: prev,
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
