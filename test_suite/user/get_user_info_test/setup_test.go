package get_user_info_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
)

func run(reqBody *pb.GetUserInfoRequest, tokens *pb.Tokens) (int, *pb.GetUserInfoResponse) {
	err, response := test_suite.NewMockTest(test_suite.GetUserInfo).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetUserInfoResponse)
	return err, resp
}
