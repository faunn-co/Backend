package get_referral_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetReferralListRequest) (int, *pb.GetReferralListResponse) {
	err, response := test_suite.NewMockTest(test_suite.GetReferralsList).Req(reqBody).Response()
	resp := response.(*pb.GetReferralListResponse)
	return err, resp
}
