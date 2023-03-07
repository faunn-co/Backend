package get_available_slot

import (
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
)

type GetAvailableSlot struct {
	req *pb.GetAvailableSlotRequest
}

func New() *GetAvailableSlot {
	g := new(GetAvailableSlot)
	g.req = new(pb.GetAvailableSlotRequest)
	return g
}

func (g *GetAvailableSlot) Request() *pb.GetAvailableSlotRequest {
	return g.req
}
