package generate_mock_data

import (
	"context"
	"errors"
	"github.com/aaronangxz/AffiliateManager/logger"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/referral"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/test_utils"
	"github.com/aaronangxz/AffiliateManager/test_suite/test_builder/user"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

type GenerateMockData struct {
	c   echo.Context
	ctx context.Context
	req *pb.GenerateMockDataRequest
}

func New(c echo.Context) *GenerateMockData {
	g := new(GenerateMockData)
	g.c = c
	g.ctx = logger.NewCtx(g.c)
	logger.Info(g.ctx, "GenerateMockData Initialized")
	return g
}

func (g *GenerateMockData) GenerateMockDataImpl() (*pb.MockDataCount, *resp.Error) {
	if err := g.verifyGenerateMockData(); err != nil {
		return nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var (
		affiliateIds       []int64
		totalReferralCount int64
		referralCount      int64
		totalBookings      int64
		bookingSuccessRate int64
		startTime          int64
		endTime            int64
	)

	for i := 0; i < int(g.req.GetMockData().GetAffiliateCount()); i++ {
		u := user.New().SetNeedLogin(false).Build()
		affiliateIds = append(affiliateIds, u.AffiliateInfo.GetUserId())
	}

	if g.req.GetMockData().ReferralsEachAffiliate != nil {
		referralCount = g.req.GetMockData().GetReferralsEachAffiliate()
	}

	if g.req.GetMockData().BookingSuccessRate != nil {
		bookingSuccessRate = 100 / g.req.GetMockData().GetBookingSuccessRate()
	}

	if g.req.GetMockData().StartDate != nil && g.req.GetMockData().EndDate != nil {
		convertedStart, _ := utils.GetDayTSWithDate(g.req.GetMockData().GetStartDate())
		convertedEnd, _ := utils.GetDayTSWithDate(g.req.GetMockData().GetEndDate())

		startTime, _ = utils.DayStartEndDate(convertedStart)
		_, endTime = utils.DayStartEndDate(convertedEnd)
	}

	for _, a := range affiliateIds {
		if g.req.GetMockData().MinReferrals != nil && g.req.GetMockData().MaxReferrals != nil {
			referralCount = test_utils.RandomRange(int(g.req.GetMockData().GetMinReferrals()), int(g.req.GetMockData().GetMaxReferrals()))
			totalReferralCount += referralCount
		}

		if g.req.GetMockData().MinBookingSuccessRate != nil && g.req.GetMockData().MaxBookingSuccessRate != nil {
			bookingSuccessRate = 100 / test_utils.RandomRange(int(g.req.GetMockData().GetMinBookingSuccessRate()), int(g.req.GetMockData().GetMaxBookingSuccessRate()))
		}

		for i := 0; i < int(referralCount); i++ {
			isSuccess := false
			if int64(i)%bookingSuccessRate == 0 {
				totalBookings++
				isSuccess = true
			}
			refClickTime := test_utils.RandomRangeInt64(startTime, endTime)
			bookingTime := test_utils.RandomRangeInt64(refClickTime, endTime)
			referral.New().SetAffiliateId(a).SetHasBooking(isSuccess).SetReferralClickTime(refClickTime).SetBookingTime(bookingTime).Build()
		}
	}
	return &pb.MockDataCount{
		TotalAffiliates: proto.Int64(int64(len(affiliateIds))),
		TotalReferrals:  proto.Int64(totalReferralCount),
		TotalBookings:   proto.Int64(totalBookings),
	}, nil
}

func (g *GenerateMockData) verifyGenerateMockData() error {
	g.req = new(pb.GenerateMockDataRequest)
	if err := g.c.Bind(g.req); err != nil {
		return err
	}

	if g.req.MockData == nil {
		return errors.New("request body is required")
	}

	if g.req.MockData.AffiliateCount == nil {
		return errors.New("specify affiliate_count")
	}

	if g.req.MockData.ReferralsEachAffiliate == nil && g.req.MockData.MinReferrals == nil && g.req.MockData.MaxReferrals == nil {
		return errors.New("specify either referrals_each_affiliate or min_referrals/max_referrals")
	}

	if g.req.MockData.ReferralsEachAffiliate != nil && (g.req.MockData.MinReferrals != nil || g.req.MockData.MaxReferrals != nil) {
		return errors.New("specify either referrals_each_affiliate or min_referrals/max_referrals")
	}

	if g.req.MockData.MinReferrals == nil && g.req.MockData.MaxReferrals == nil {
		return errors.New("specify both min_referrals and max_referrals")
	}

	if g.req.MockData.MinReferrals == nil && g.req.MockData.MaxReferrals != nil {
		return errors.New("specify both min_referrals and max_referrals")
	}

	if g.req.MockData.MinReferrals != nil && g.req.MockData.MaxReferrals == nil {
		return errors.New("specify both min_referrals and max_referrals")
	}

	if g.req.MockData.BookingSuccessRate == nil && g.req.MockData.MinBookingSuccessRate == nil && g.req.MockData.MaxBookingSuccessRate == nil {
		return errors.New("specify booking_success_rate or min_booking_success_rate/max_booking_success_rate")
	} else {
		if g.req.MockData.MinBookingSuccessRate == nil && g.req.MockData.MaxBookingSuccessRate == nil {
			return errors.New("specify both min_booking_success_rate and max_booking_success_rate")
		}

		if g.req.MockData.MinBookingSuccessRate == nil && g.req.MockData.MaxBookingSuccessRate != nil {
			return errors.New("specify both min_booking_success_rate and max_booking_success_rate")
		}

		if g.req.MockData.MinBookingSuccessRate != nil && g.req.MockData.MaxBookingSuccessRate == nil {
			return errors.New("specify both min_booking_success_rate and max_booking_success_rate")
		}
	}

	if g.req.MockData.StartDate == nil || g.req.MockData.EndDate == nil {
		return errors.New("specify start_date and end_date")
	}
	return nil
}
