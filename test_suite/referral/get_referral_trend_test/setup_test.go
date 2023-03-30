package get_referral_trend_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetReferralTrendRequest, tokens *pb.Tokens) (int, *pb.GetReferralTrendResponse) {
	err, response := test_suite.NewMockTest(test_suite.ReferralGetReferralTrend).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetReferralTrendResponse)
	return err, resp
}
