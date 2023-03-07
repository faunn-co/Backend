package get_affiliate_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetAffiliateListRequest, tokens *pb.Tokens) (int, *pb.GetAffiliateListResponse) {
	err, response := test_suite.NewMockTest(test_suite.GetAffiliateList).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetAffiliateListResponse)
	return err, resp
}
