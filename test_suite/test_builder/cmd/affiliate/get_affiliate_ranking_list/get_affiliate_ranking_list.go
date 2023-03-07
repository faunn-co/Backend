package get_affiliate_ranking_list

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
)

type GetAffiliateRankingList struct {
	req *pb.GetAffiliateRankingListRequest
}

func New() *GetAffiliateRankingList {
	g := new(GetAffiliateRankingList)
	return g
}

func (g *GetAffiliateRankingList) Request() *pb.GetAffiliateRankingListRequest {
	return g.req
}
