package get_affiliate_stats_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetAffiliateStatsRequest, tokens *pb.Tokens) (int, *pb.GetAffiliateStatsResponse) {
	err, response := test_suite.NewMockTest(test_suite.AffiliateGetAffiliateStats).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetAffiliateStatsResponse)
	return err, resp
}
