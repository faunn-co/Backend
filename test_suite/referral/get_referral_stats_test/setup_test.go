package get_referral_stats_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetReferralStatsRequest, tokens *pb.Tokens) (int, *pb.GetReferralStatsResponse) {
	err, response := test_suite.NewMockTest(test_suite.GetReferralStats).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetReferralStatsResponse)
	return err, resp
}
