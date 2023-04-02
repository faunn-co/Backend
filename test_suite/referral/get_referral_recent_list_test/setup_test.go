package get_referral_recent_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetReferralRecentListRequest, tokens *pb.Tokens) (int, *pb.GetReferralRecentListResponse) {
	err, response := test_suite.NewMockTest(test_suite.ReferralGetReferralRecentList).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetReferralRecentListResponse)
	return err, resp
}
