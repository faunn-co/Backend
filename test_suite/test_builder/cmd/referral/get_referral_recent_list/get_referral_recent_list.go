package get_referral_recent_list

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
)

type GetReferralRecentList struct {
	req *pb.GetReferralRecentListRequest
}

func New() *GetReferralRecentList {
	g := new(GetReferralRecentList)
	return g
}

func (g *GetReferralRecentList) Request() *pb.GetReferralRecentListRequest {
	return g.req
}
