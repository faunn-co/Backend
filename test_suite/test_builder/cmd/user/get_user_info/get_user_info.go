package get_user_info

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
)

type GetUserInfo struct {
	req *pb.GetUserInfoRequest
}

func New() *GetUserInfo {
	g := new(GetUserInfo)
	g.req = new(pb.GetUserInfoRequest)
	return g
}

func (g *GetUserInfo) Request() *pb.GetUserInfoRequest {
	return g.req
}
