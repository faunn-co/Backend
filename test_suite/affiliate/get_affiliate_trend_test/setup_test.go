package get_affiliate_trend_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetAffiliateTrendRequest, tokens *pb.Tokens) (int, *pb.GetAffiliateTrendResponse) {
	err, response := test_suite.NewMockTest(test_suite.AffiliateGetAffiliateTrend).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetAffiliateTrendResponse)
	return err, resp
}
