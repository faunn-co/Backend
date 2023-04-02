package get_affiliate_ranking_list_test

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/test_suite"
	"strconv"
)

func run(reqBody *pb.GetAffiliateRankingListRequest, tokens *pb.Tokens, period *pb.TimeSelectorPeriod) (int, *pb.GetAffiliateRankingListResponse) {
	var q string
	if period == nil {
		q = ""
	} else {
		q = strconv.FormatInt(int64(*period), 10)
	}
	err, response := test_suite.NewMockTest(test_suite.AffiliateGetAffiliateRankingList).QueryParam("period", q).Req(reqBody, tokens).Response()
	resp := response.(*pb.GetAffiliateRankingListResponse)
	return err, resp
}
