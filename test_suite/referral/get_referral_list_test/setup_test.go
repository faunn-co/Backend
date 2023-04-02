package get_referral_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetReferralListRequest, tokens *pb.Tokens) (int, *pb.GetReferralListResponse) {
	err, response := test_suite.NewMockTest(test_suite.ReferralGetReferralsList).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetReferralListResponse)
	return err, resp
}
